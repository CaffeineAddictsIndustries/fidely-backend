package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"fidely-backend/internal/auth"
)

const (
	MessageLoginSuccessful  = "login successful"
	MessageLoginFailed      = "login failed"
	MessageLogoutSuccessful = "logout successful"
	MessageLogoutFailed     = "logout failed"
	MessageAccessDenied     = "access denied"

	ReasonInvalidCredentials     = "invalid credentials"
	ReasonMissingCredentials     = "username and password are required"
	ReasonInvalidSession         = "invalid session"
	ReasonInsufficientPrivileges = "insufficient privileges"
)

// dummyPasswordHash is used to normalize login timing when username doesn't exist.
// Hash corresponds to a random throwaway value with bcrypt cost 12.
const dummyPasswordHash = "$2a$12$Qf7kBfko5i7vAlNfL9Zwme6x6A0jvEeeQ6J0Pmv6Qp0WQ3f4o0XfS"

var (
	ErrAuthRepositoryRequired  = errors.New("auth repository is required")
	ErrPasswordManagerRequired = errors.New("password manager is required")
	ErrSessionManagerRequired  = errors.New("session manager is required")
	ErrInvalidSessionTTL       = errors.New("session ttl must be positive")
)

// AdminPrincipal is a normalized admin identity used by auth flows.
type AdminPrincipal struct {
	AdminType    auth.AdminType
	AdminID      int
	Username     string
	PasswordHash string
	Role         string
	StoreID      *int
}

// LoginResult follows the UI contract: general message plus optional reason.
type LoginResult struct {
	Success      bool
	Message      string
	Reason       string
	SessionToken string
	ExpiresAt    time.Time
	AdminType    auth.AdminType
	AdminID      int
}

// LogoutResult follows the same message contract pattern.
type LogoutResult struct {
	Success bool
	Message string
	Reason  string
}

// AdminAuthRepository defines persistence needed by auth service.
type AdminAuthRepository interface {
	FindByUsername(ctx context.Context, username string) (*AdminPrincipal, error)
	FindByTypeAndID(ctx context.Context, adminType auth.AdminType, adminID int) (*AdminPrincipal, error)
	CreateSession(ctx context.Context, session auth.AdminSession) (auth.AdminSession, error)
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*auth.AdminSession, error)
	UpdateSession(ctx context.Context, session auth.AdminSession) error
}

// AdminAuthService implements admin-only login/session workflows.
type AdminAuthService struct {
	repo       AdminAuthRepository
	passwords  *auth.PasswordManager
	sessions   *auth.SessionManager
	sessionTTL time.Duration
}

func NewAdminAuthService(
	repo AdminAuthRepository,
	passwords *auth.PasswordManager,
	sessions *auth.SessionManager,
	sessionTTL time.Duration,
) (*AdminAuthService, error) {
	if repo == nil {
		return nil, ErrAuthRepositoryRequired
	}
	if passwords == nil {
		return nil, ErrPasswordManagerRequired
	}
	if sessions == nil {
		return nil, ErrSessionManagerRequired
	}
	if sessionTTL <= 0 {
		return nil, ErrInvalidSessionTTL
	}

	return &AdminAuthService{
		repo:       repo,
		passwords:  passwords,
		sessions:   sessions,
		sessionTTL: sessionTTL,
	}, nil
}

// Login authenticates an admin and creates a new session token.
func (service *AdminAuthService) Login(ctx context.Context, username string, password string) (LoginResult, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return LoginResult{
			Success: false,
			Message: MessageLoginFailed,
			Reason:  ReasonMissingCredentials,
		}, nil
	}

	principal, err := service.repo.FindByUsername(ctx, username)
	if err != nil {
		return LoginResult{}, fmt.Errorf("find admin by username: %w", err)
	}
	if principal == nil {
		// Keep timing closer to valid-username path to reduce enumeration signals.
		_ = service.passwords.Verify(dummyPasswordHash, password)
		return LoginResult{
			Success: false,
			Message: MessageLoginFailed,
			Reason:  ReasonInvalidCredentials,
		}, nil
	}

	if err := service.passwords.Verify(principal.PasswordHash, password); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return LoginResult{
				Success: false,
				Message: MessageLoginFailed,
				Reason:  ReasonInvalidCredentials,
			}, nil
		}
		return LoginResult{}, fmt.Errorf("verify password: %w", err)
	}

	session, token, err := service.sessions.NewSession(principal.AdminType, principal.AdminID, service.sessionTTL)
	if err != nil {
		return LoginResult{}, fmt.Errorf("create session: %w", err)
	}

	persistedSession, err := service.repo.CreateSession(ctx, session)
	if err != nil {
		return LoginResult{}, fmt.Errorf("persist session: %w", err)
	}

	return LoginResult{
		Success:      true,
		Message:      MessageLoginSuccessful,
		Reason:       "",
		SessionToken: token,
		ExpiresAt:    persistedSession.ExpiresAt,
		AdminType:    principal.AdminType,
		AdminID:      principal.AdminID,
	}, nil
}

// Logout revokes a session identified by its plaintext token.
func (service *AdminAuthService) Logout(ctx context.Context, sessionToken string) (LogoutResult, error) {
	tokenHash, err := service.sessions.HashToken(sessionToken)
	if err != nil {
		return LogoutResult{
			Success: false,
			Message: MessageLogoutFailed,
			Reason:  ReasonInvalidSession,
		}, nil
	}

	session, err := service.repo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return LogoutResult{}, fmt.Errorf("get session by token hash: %w", err)
	}
	if session == nil {
		return LogoutResult{
			Success: false,
			Message: MessageLogoutFailed,
			Reason:  ReasonInvalidSession,
		}, nil
	}

	revokedSession, err := service.sessions.RevokeSession(*session)
	if err != nil {
		return LogoutResult{}, fmt.Errorf("revoke session: %w", err)
	}
	if err := service.repo.UpdateSession(ctx, revokedSession); err != nil {
		return LogoutResult{}, fmt.Errorf("update session: %w", err)
	}

	return LogoutResult{
		Success: true,
		Message: MessageLogoutSuccessful,
	}, nil
}

// Authenticate resolves an admin principal from a plaintext session token.
func (service *AdminAuthService) Authenticate(ctx context.Context, sessionToken string) (*AdminPrincipal, error) {
	tokenHash, err := service.sessions.HashToken(sessionToken)
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	session, err := service.repo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("get session by token hash: %w", err)
	}
	if session == nil {
		return nil, auth.ErrInvalidSession
	}
	if !service.sessions.IsSessionActive(*session) {
		return nil, auth.ErrInvalidSession
	}

	principal, err := service.repo.FindByTypeAndID(ctx, session.AdminType, session.AdminID)
	if err != nil {
		return nil, fmt.Errorf("find principal by type and id: %w", err)
	}
	if principal == nil {
		return nil, auth.ErrInvalidSession
	}

	now := time.Now().UTC()
	session.LastSeenAt = &now
	if err := service.repo.UpdateSession(ctx, *session); err != nil {
		return nil, fmt.Errorf("touch session: %w", err)
	}

	return principal, nil
}
