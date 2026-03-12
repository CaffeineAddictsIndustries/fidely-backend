package auth

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashAndVerify(t *testing.T) {
	manager := NewDefaultPasswordManager()

	hash, err := manager.Hash("StrongPassword123!")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	if hash == "StrongPassword123!" {
		t.Fatal("hash should never equal plaintext password")
	}

	if err := manager.Verify(hash, "StrongPassword123!"); err != nil {
		t.Fatalf("Verify returned error for valid credentials: %v", err)
	}
}

func TestVerifyInvalidPassword(t *testing.T) {
	manager := NewDefaultPasswordManager()

	hash, err := manager.Hash("StrongPassword123!")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	err = manager.Verify(hash, "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestNeedsRehash(t *testing.T) {
	manager, err := NewPasswordManager(RecommendedBcryptCost)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	lowCostHash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), RecommendedBcryptCost-2)
	if err != nil {
		t.Fatalf("failed to generate low-cost hash: %v", err)
	}

	if !manager.NeedsRehash(string(lowCostHash)) {
		t.Fatal("expected low-cost hash to require rehash")
	}

	currentHash, err := manager.Hash("StrongPassword123!")
	if err != nil {
		t.Fatalf("failed to generate current-cost hash: %v", err)
	}

	if manager.NeedsRehash(currentHash) {
		t.Fatal("did not expect current-cost hash to require rehash")
	}
}
