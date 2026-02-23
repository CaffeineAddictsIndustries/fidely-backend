package model

import "time"

// Store represents the stores table.
// A store is a business that uses the Fidely loyalty platform.
type Store struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StoreCreate is used when creating a new store.
type StoreCreate struct {
	Name string `json:"name"`
}

// StoreUpdate is used when updating an existing store.
type StoreUpdate struct {
	Name string `json:"name"`
}
