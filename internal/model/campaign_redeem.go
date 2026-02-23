package model

import "time"

// CampaignRedeem represents the campaign_redeems table.
// Records when a user redeems a campaign reward.
type CampaignRedeem struct {
	ID         int       `json:"id"`
	CampaignID int       `json:"campaign_id"`
	CardID     int       `json:"card_id"`
	Amount     int       `json:"amount"`
	Points     int       `json:"points"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CampaignRedeemCreate is used when creating a new campaign redeem.
type CampaignRedeemCreate struct {
	CampaignID int `json:"campaign_id"`
	CardID     int `json:"card_id"`
	Amount     int `json:"amount"`
	Points     int `json:"points"`
}
