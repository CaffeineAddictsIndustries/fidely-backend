package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"fidely-backend/internal/auth"
)

type inMemoryAuthRepo struct {
	byUsername map[string]*AdminPrincipal
	byKey      map[string]*AdminPrincipal
	sessions   map[string]auth.AdminSession
	nextID     int
}

func newInMemoryAuthRepo() *inMemoryAuthRepo {
	return &inMemoryAuthRepo{
		byUsername: make(map[string]*AdminPrincipal),
		byKey:      make(map[string]*AdminPrincipal),
		sessions:   make(map[string]auth.AdminSession),
		nextID:     1,
	}
}

func principalKey(adminType auth.AdminType, adminID int) string {
	return fmt.Sprintf("%s:%d", adminType, adminID)
}

func (repo *inMemoryAuthRepo) FindByUsername(_ context.Context, username string) (*AdminPrincipal, error) {
	principal := repo.byUsername[username]
	return principal, nil
}

func (repo *inMemoryAuthRepo) FindByTypeAndID(_ context.Context, adminType auth.AdminType, adminID int) (*AdminPrincipal, error) {
	principal := repo.byKey[principalKey(adminType, adminID)]
	return principal, nil
}

func (repo *inMemoryAuthRepo) CreateSession(_ context.Context, session auth.AdminSession) (auth.AdminSession, error) {
	session.ID = repo.nextID
	repo.nextID++
	repo.sessions[session.SessionTokenHash] = session
	return session, nil
}

func (repo *inMemoryAuthRepo) GetSessionByTokenHash(_ context.Context, tokenHash string) (*auth.AdminSession, error) {
	session, ok := repo.sessions[tokenHash]
	if !ok {
		return nil, nil
	}
	copySession := session
	return &copySession, nil
}

func (repo *inMemoryAuthRepo) UpdateSession(_ context.Context, session auth.AdminSession) error {
	repo.sessions[session.SessionTokenHash] = session
	return nil
}

func TestAdminAuthServiceLoginSuccess(t *testing.T) {
	repo := newInMemoryAuthRepo()
	passwords := auth.NewDefaultPasswordManager()
	sessions := auth.NewSessionManager("pepper")

	hash, err := passwords.Hash("S3cur3Password!")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	principal := &AdminPrincipal{
		AdminType:    auth.AdminTypeStoreAdmin,
		AdminID:      10,
		Username:     "alice",
		PasswordHash: hash,
		Role:         "manager",
	}
	repo.byUsername[principal.Username] = principal
	repo.byKey[principalKey(principal.AdminType, principal.AdminID)] = principal

	service, err := NewAdminAuthService(repo, passwords, sessions, 2*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	result, err := service.Login(context.Background(), "alice", "S3cur3Password!")
	if err != nil {
		t.Fatalf("login returned error: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success result, got: %+v", result)
	}
	if result.Message != MessageLoginSuccessful {
		t.Fatalf("unexpected success message: %s", result.Message)
	}
	if result.SessionToken == "" {
		t.Fatal("expected non-empty session token")
	}
}

func TestAdminAuthServiceLoginUsernameNotFound(t *testing.T) {
	repo := newInMemoryAuthRepo()
	passwords := auth.NewDefaultPasswordManager()
	sessions := auth.NewSessionManager("pepper")

	service, err := NewAdminAuthService(repo, passwords, sessions, 2*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	result, err := service.Login(context.Background(), "missing", "any")
	if err != nil {
		t.Fatalf("login returned error: %v", err)
	}
	if result.Success {
		t.Fatalf("expected failed result, got: %+v", result)
	}
	if result.Message != MessageLoginFailed {
		t.Fatalf("unexpected failed message: %s", result.Message)
	}
	if result.Reason != ReasonUsernameDoesNotExist {
		t.Fatalf("unexpected reason: %s", result.Reason)
	}
}

func TestAdminAuthServiceLoginIncorrectPassword(t *testing.T) {
	repo := newInMemoryAuthRepo()
	passwords := auth.NewDefaultPasswordManager()
	sessions := auth.NewSessionManager("pepper")

	hash, err := passwords.Hash("S3cur3Password!")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	principal := &AdminPrincipal{
		AdminType:    auth.AdminTypeFidelyAdmin,
		AdminID:      55,
		Username:     "platform-admin",
		PasswordHash: hash,
		Role:         "owner",
	}
	repo.byUsername[principal.Username] = principal
	repo.byKey[principalKey(principal.AdminType, principal.AdminID)] = principal

	service, err := NewAdminAuthService(repo, passwords, sessions, 2*time.Hour)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	result, err := service.Login(context.Background(), "platform-admin", "wrong")
	if err != nil {
		t.Fatalf("login returned error: %v", err)
	}
	if result.Success {
		t.Fatalf("expected failed result, got: %+v", result)
	}
	if result.Reason != ReasonIncorrectPassword {
		t.Fatalf("unexpected reason: %s", result.Reason)
	}
}
