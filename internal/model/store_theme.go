package model

import "time"

// StoreTheme represents the store_theme table.
// Contains visual customization settings for a store's loyalty cards.
type StoreTheme struct {
	ID                     int       `json:"id"`
	StoreID                int       `json:"store_id"`
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

// StoreThemeCreate is used when creating a new store theme.
type StoreThemeCreate struct {
	StoreID                int    `json:"store_id"`
	PrimaryColor           string `json:"primary_color"`
	SecondaryColor         string `json:"secondary_color"`
	BackgroundColor        string `json:"background_color"`
	TextColor              string `json:"text_color"`
	StampImageURL          string `json:"stamp_image_url"`
	SlotImageURL           string `json:"slot_image_url"`
	CardBackgroundImageURL string `json:"card_background_image_url"`
	LogoURL                string `json:"logo_url"`
}

// StoreThemeUpdate is used when updating an existing store theme.
type StoreThemeUpdate struct {
	PrimaryColor           string `json:"primary_color"`
	SecondaryColor         string `json:"secondary_color"`
	BackgroundColor        string `json:"background_color"`
	TextColor              string `json:"text_color"`
	StampImageURL          string `json:"stamp_image_url"`
	SlotImageURL           string `json:"slot_image_url"`
	CardBackgroundImageURL string `json:"card_background_image_url"`
	LogoURL                string `json:"logo_url"`
}
