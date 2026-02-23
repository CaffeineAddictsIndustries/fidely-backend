package model

import "time"

// User represents the users table.
// A user is a customer who collects loyalty points.
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose password in JSON
	Token     string    `json:"-"` // Never expose token in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCreate is used when creating a new user.
type UserCreate struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// UserUpdate is used when updating an existing user.
type UserUpdate struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
