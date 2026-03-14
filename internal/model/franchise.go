package model

import "time"

// Franchise represents the franchises table.
// A franchise can own multiple stores, card types, and store admins.
type Franchise struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FranchiseCreate is used when creating a new franchise.
type FranchiseCreate struct {
	Name string `json:"name"`
}

// FranchiseUpdate is used when updating an existing franchise.
type FranchiseUpdate struct {
	Name string `json:"name"`
}
