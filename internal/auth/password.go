package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// RecommendedBcryptCost is a practical default for admin logins.
	RecommendedBcryptCost = 12
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmptyPassword      = errors.New("password cannot be empty")
)

// PasswordManager encapsulates password hashing policy.
type PasswordManager struct {
	cost int
}

// NewPasswordManager builds a manager with a fixed bcrypt cost.
func NewPasswordManager(cost int) (*PasswordManager, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("bcrypt cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	return &PasswordManager{cost: cost}, nil
}

// NewDefaultPasswordManager builds a manager using secure defaults.
func NewDefaultPasswordManager() *PasswordManager {
	return &PasswordManager{cost: RecommendedBcryptCost}
}

// Hash returns a bcrypt hash for the provided password.
func (manager *PasswordManager) Hash(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), manager.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// Verify compares password against a stored hash.
func (manager *PasswordManager) Verify(passwordHash string, password string) error {
	if passwordHash == "" || password == "" {
		return ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

// NeedsRehash indicates whether a hash should be regenerated with the current policy.
func (manager *PasswordManager) NeedsRehash(passwordHash string) bool {
	if passwordHash == "" {
		return true
	}

	cost, err := bcrypt.Cost([]byte(passwordHash))
	if err != nil {
		return true
	}

	return cost < manager.cost
}
