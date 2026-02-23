package model

import "time"

// Campaign represents the campaigns table.
// A campaign defines a reward that can be redeemed for points.
type Campaign struct {
	ID         int       `json:"id"`
	CardTypeID int       `json:"card_type_id"`
	Points     int       `json:"points"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CampaignCreate is used when creating a new campaign.
type CampaignCreate struct {
	CardTypeID int `json:"card_type_id"`
	Points     int `json:"points"`
}

// CampaignUpdate is used when updating an existing campaign.
type CampaignUpdate struct {
	Points int `json:"points"`
}
