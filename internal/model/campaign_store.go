package model

import "time"

// CampaignStore represents the campaign_stores join table.
// A campaign can be associated with zero or multiple stores.
type CampaignStore struct {
	CampaignID int       `json:"campaign_id"`
	StoreID    int       `json:"store_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CampaignStoreCreate is used when associating a campaign with a store.
type CampaignStoreCreate struct {
	CampaignID int `json:"campaign_id"`
	StoreID    int `json:"store_id"`
}
