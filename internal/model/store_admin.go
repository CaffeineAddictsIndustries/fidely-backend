package model

import "time"

// StoreAdmin represents the store_admin table.
// A store admin manages a store on the platform.
type StoreAdmin struct {
	ID           int       `json:"id"`
	StoreID      *int      `json:"store_id,omitempty"`
	FranchiseID  *int      `json:"franchise_id,omitempty"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose password hashes in JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// StoreAdminCreate is used when creating a new store admin.
type StoreAdminCreate struct {
	StoreID     *int   `json:"store_id,omitempty"`
	FranchiseID *int   `json:"franchise_id,omitempty"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

// StoreAdminUpdate is used when updating an existing store admin.
type StoreAdminUpdate struct {
	StoreID     *int   `json:"store_id,omitempty"`
	FranchiseID *int   `json:"franchise_id,omitempty"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}
