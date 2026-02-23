package model

import "time"

// CardType represents the card_types table.
// A card type defines a loyalty card template for a store.
type CardType struct {
	ID        int       `json:"id"`
	StoreID   int       `json:"store_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CardTypeCreate is used when creating a new card type.
type CardTypeCreate struct {
	StoreID int    `json:"store_id"`
	Name    string `json:"name"`
}

// CardTypeUpdate is used when updating an existing card type.
type CardTypeUpdate struct {
	Name string `json:"name"`
}
