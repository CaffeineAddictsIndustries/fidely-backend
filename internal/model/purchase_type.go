package model

import "time"

// PurchaseType represents the purchase_types table.
// Defines types of purchases that can earn points for a card type.
type PurchaseType struct {
	ID          int       `json:"id"`
	CardTypeID  int       `json:"card_type_id"`
	Points      int       `json:"points"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PurchaseTypeCreate is used when creating a new purchase type.
type PurchaseTypeCreate struct {
	CardTypeID  int    `json:"card_type_id"`
	Points      int    `json:"points"`
	Description string `json:"description"`
}

// PurchaseTypeUpdate is used when updating an existing purchase type.
type PurchaseTypeUpdate struct {
	Points      int    `json:"points"`
	Description string `json:"description"`
}
