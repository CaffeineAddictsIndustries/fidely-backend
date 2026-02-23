package model

import "time"

// Card represents the cards table.
// A card is an instance of a card type belonging to a user, with accumulated points.
type Card struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CardTypeID int       `json:"card_type_id"`
	Points     int       `json:"points"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CardCreate is used when creating a new card.
type CardCreate struct {
	UserID     int `json:"user_id"`
	CardTypeID int `json:"card_type_id"`
	Points     int `json:"points"`
}

// CardUpdate is used when updating an existing card.
type CardUpdate struct {
	Points int `json:"points"`
}
