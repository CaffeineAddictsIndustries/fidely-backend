package model

import "time"

// FranchiseTheme represents the franchise_theme table.
// Contains visual customization settings for a franchise.
type FranchiseTheme struct {
	ID                     int       `json:"id"`
	FranchiseID            int       `json:"franchise_id"`
	PrimaryColor           string    `json:"primary_color"`
	SecondaryColor         string    `json:"secondary_color"`
	BackgroundColor        string    `json:"background_color"`
	TextColor              string    `json:"text_color"`
	StampImageURL          string    `json:"stamp_image_url"`
	SlotImageURL           string    `json:"slot_image_url"`
	CardBackgroundImageURL string    `json:"card_background_image_url"`
	LogoURL                string    `json:"logo_url"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// FranchiseThemeCreate is used when creating a new franchise theme.
type FranchiseThemeCreate struct {
	FranchiseID            int    `json:"franchise_id"`
	PrimaryColor           string `json:"primary_color"`
	SecondaryColor         string `json:"secondary_color"`
	BackgroundColor        string `json:"background_color"`
	TextColor              string `json:"text_color"`
	StampImageURL          string `json:"stamp_image_url"`
	SlotImageURL           string `json:"slot_image_url"`
	CardBackgroundImageURL string `json:"card_background_image_url"`
	LogoURL                string `json:"logo_url"`
}

// FranchiseThemeUpdate is used when updating an existing franchise theme.
type FranchiseThemeUpdate struct {
	PrimaryColor           string `json:"primary_color"`
	SecondaryColor         string `json:"secondary_color"`
	BackgroundColor        string `json:"background_color"`
	TextColor              string `json:"text_color"`
	StampImageURL          string `json:"stamp_image_url"`
	SlotImageURL           string `json:"slot_image_url"`
	CardBackgroundImageURL string `json:"card_background_image_url"`
	LogoURL                string `json:"logo_url"`
}
