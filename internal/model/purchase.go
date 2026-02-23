package model

import "time"

// Purchase represents the purchases table.
// Records a purchase transaction that earned points.
type Purchase struct {
	ID             int       `json:"id"`
	PurchaseTypeID int       `json:"purchase_type_id"`
	CardID         int       `json:"card_id"`
	Points         int       `json:"points"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PurchaseCreate is used when creating a new purchase.
type PurchaseCreate struct {
	PurchaseTypeID int `json:"purchase_type_id"`
	CardID         int `json:"card_id"`
	Points         int `json:"points"`
}
