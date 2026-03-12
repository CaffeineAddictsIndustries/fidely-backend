package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

const (
	defaultSessionTokenBytes = 32
)

var (
	ErrInvalidAdminType = errors.New("invalid admin type")
	ErrInvalidAdminID   = errors.New("invalid admin id")
	ErrInvalidSession   = errors.New("invalid session")
	ErrInvalidTTL       = errors.New("invalid session ttl")
	ErrInvalidToken     = errors.New("invalid session token")
)

// AdminType indicates which admin table a principal belongs to.
type AdminType string

const (
	AdminTypeStoreAdmin  AdminType = "store_admin"
	AdminTypeFidelyAdmin AdminType = "fidely_admin"
)

// AdminSession mirrors persisted admin session records.
type AdminSession struct {
	ID               int
	AdminType        AdminType
	AdminID          int
	SessionTokenHash string
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	CreatedAt        time.Time
	LastSeenAt       *time.Time
}

// SessionManager centralizes session token creation and lifecycle checks.
type SessionManager struct {
	tokenBytes int
	pepper     string
	nowFunc    func() time.Time
}

// NewSessionManager builds a manager with secure token defaults.
func NewSessionManager(pepper string) *SessionManager {
	return &SessionManager{
		tokenBytes: defaultSessionTokenBytes,
		pepper:     pepper,
		nowFunc:    time.Now,
	}
}

// GenerateToken returns a random opaque token and its deterministic hash.
func (manager *SessionManager) GenerateToken() (string, string, error) {
	tokenBytes := make([]byte, manager.tokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random token: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	tokenHash, err := manager.HashToken(token)
	if err != nil {
		return "", "", err
	}

	return token, tokenHash, nil
}

// HashToken hashes a session token with optional pepper.
func (manager *SessionManager) HashToken(token string) (string, error) {
	if token == "" {
		return "", ErrInvalidToken
	}

	sum := sha256.Sum256([]byte(manager.pepper + ":" + token))
	return hex.EncodeToString(sum[:]), nil
}

// NewSession creates a new session record and returns the plaintext token.
func (manager *SessionManager) NewSession(adminType AdminType, adminID int, ttl time.Duration) (AdminSession, string, error) {
	if !isValidAdminType(adminType) {
		return AdminSession{}, "", ErrInvalidAdminType
	}
	if adminID <= 0 {
		return AdminSession{}, "", ErrInvalidAdminID
	}
	if ttl <= 0 {
		return AdminSession{}, "", ErrInvalidTTL
	}

	token, tokenHash, err := manager.GenerateToken()
	if err != nil {
		return AdminSession{}, "", err
	}

	now := manager.nowFunc().UTC()
	session := AdminSession{
		AdminType:        adminType,
		AdminID:          adminID,
		SessionTokenHash: tokenHash,
		CreatedAt:        now,
		ExpiresAt:        now.Add(ttl),
	}

	return session, token, nil
}

// IsSessionActive verifies expiration and revocation state.
func (manager *SessionManager) IsSessionActive(session AdminSession) bool {
	now := manager.nowFunc().UTC()
	if session.RevokedAt != nil {
		return false
	}
	if !now.Before(session.ExpiresAt) {
		return false
	}
	return true
}

// RotateSessionToken replaces the token hash and extends expiration.
func (manager *SessionManager) RotateSessionToken(session AdminSession, ttl time.Duration) (AdminSession, string, error) {
	if ttl <= 0 {
		return AdminSession{}, "", ErrInvalidTTL
	}
	if session.ID <= 0 {
		return AdminSession{}, "", ErrInvalidSession
	}
	if !isValidAdminType(session.AdminType) || session.AdminID <= 0 {
		return AdminSession{}, "", ErrInvalidSession
	}

	token, tokenHash, err := manager.GenerateToken()
	if err != nil {
		return AdminSession{}, "", err
	}

	now := manager.nowFunc().UTC()
	session.SessionTokenHash = tokenHash
	session.ExpiresAt = now.Add(ttl)
	session.LastSeenAt = &now

	return session, token, nil
}

// RevokeSession marks a session as revoked at current time.
func (manager *SessionManager) RevokeSession(session AdminSession) (AdminSession, error) {
	if session.ID <= 0 {
		return AdminSession{}, ErrInvalidSession
	}

	now := manager.nowFunc().UTC()
	session.RevokedAt = &now
	session.LastSeenAt = &now
	return session, nil
}

func isValidAdminType(adminType AdminType) bool {
	return adminType == AdminTypeStoreAdmin || adminType == AdminTypeFidelyAdmin
}
