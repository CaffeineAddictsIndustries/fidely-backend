package model

import "time"

// PushNotificationCard represents the push_notification_cards table.
// Links push notifications to specific cards (many-to-many relationship).
type PushNotificationCard struct {
	PushNotificationID int       `json:"push_notification_id"`
	CardID             int       `json:"card_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// PushNotificationCardCreate is used when linking a notification to a card.
type PushNotificationCardCreate struct {
	PushNotificationID int `json:"push_notification_id"`
	CardID             int `json:"card_id"`
}
