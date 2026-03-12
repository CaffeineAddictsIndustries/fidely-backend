package auth

import (
	"testing"
	"time"
)

func TestGenerateTokenReturnsHash(t *testing.T) {
	manager := NewSessionManager("pepper")
	manager.nowFunc = func() time.Time { return time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC) }

	token, tokenHash, err := manager.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if tokenHash == "" {
		t.Fatal("expected non-empty token hash")
	}
	if token == tokenHash {
		t.Fatal("token hash must differ from plaintext token")
	}
}

func TestNewSession(t *testing.T) {
	manager := NewSessionManager("pepper")
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	manager.nowFunc = func() time.Time { return now }

	session, token, err := manager.NewSession(AdminTypeStoreAdmin, 7, 2*time.Hour)
	if err != nil {
		t.Fatalf("NewSession returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected plaintext token")
	}
	if session.AdminType != AdminTypeStoreAdmin {
		t.Fatalf("unexpected admin type: %s", session.AdminType)
	}
	if session.AdminID != 7 {
		t.Fatalf("unexpected admin id: %d", session.AdminID)
	}
	if session.SessionTokenHash == "" {
		t.Fatal("expected token hash")
	}
	if !session.CreatedAt.Equal(now) {
		t.Fatalf("unexpected created_at: %v", session.CreatedAt)
	}
	if !session.ExpiresAt.Equal(now.Add(2 * time.Hour)) {
		t.Fatalf("unexpected expires_at: %v", session.ExpiresAt)
	}
}

func TestIsSessionActive(t *testing.T) {
	manager := NewSessionManager("pepper")
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	manager.nowFunc = func() time.Time { return now }

	session := AdminSession{
		ID:        1,
		AdminType: AdminTypeFidelyAdmin,
		AdminID:   2,
		ExpiresAt: now.Add(30 * time.Minute),
	}
	if !manager.IsSessionActive(session) {
		t.Fatal("expected session to be active")
	}

	expired := session
	expired.ExpiresAt = now.Add(-1 * time.Minute)
	if manager.IsSessionActive(expired) {
		t.Fatal("expected expired session to be inactive")
	}

	revoked := session
	revokedAt := now.Add(-5 * time.Minute)
	revoked.RevokedAt = &revokedAt
	if manager.IsSessionActive(revoked) {
		t.Fatal("expected revoked session to be inactive")
	}
}

func TestRotateSessionToken(t *testing.T) {
	manager := NewSessionManager("pepper")
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	manager.nowFunc = func() time.Time { return now }

	session := AdminSession{
		ID:               3,
		AdminType:        AdminTypeStoreAdmin,
		AdminID:          22,
		SessionTokenHash: "old-hash",
		CreatedAt:        now.Add(-1 * time.Hour),
		ExpiresAt:        now.Add(1 * time.Hour),
	}

	updated, token, err := manager.RotateSessionToken(session, 3*time.Hour)
	if err != nil {
		t.Fatalf("RotateSessionToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected rotated plaintext token")
	}
	if updated.SessionTokenHash == "" || updated.SessionTokenHash == "old-hash" {
		t.Fatal("expected updated session token hash")
	}
	if !updated.ExpiresAt.Equal(now.Add(3 * time.Hour)) {
		t.Fatalf("unexpected updated expires_at: %v", updated.ExpiresAt)
	}
	if updated.LastSeenAt == nil || !updated.LastSeenAt.Equal(now) {
		t.Fatalf("unexpected last_seen_at: %v", updated.LastSeenAt)
	}
}

func TestRevokeSession(t *testing.T) {
	manager := NewSessionManager("pepper")
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	manager.nowFunc = func() time.Time { return now }

	session := AdminSession{ID: 9, AdminType: AdminTypeFidelyAdmin, AdminID: 5}
	revoked, err := manager.RevokeSession(session)
	if err != nil {
		t.Fatalf("RevokeSession returned error: %v", err)
	}
	if revoked.RevokedAt == nil || !revoked.RevokedAt.Equal(now) {
		t.Fatalf("unexpected revoked_at: %v", revoked.RevokedAt)
	}
	if revoked.LastSeenAt == nil || !revoked.LastSeenAt.Equal(now) {
		t.Fatalf("unexpected last_seen_at: %v", revoked.LastSeenAt)
	}
}
