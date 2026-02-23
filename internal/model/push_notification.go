package model

import "time"

// PushNotification represents the push_notifications table.
// Defines a push notification to be sent to users.
type PushNotification struct {
	ID          int       `json:"id"`
	CardTypeID  int       `json:"card_type_id"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
	Message     string    `json:"message"`
	PushDate    time.Time `json:"push_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PushNotificationCreate is used when creating a new push notification.
type PushNotificationCreate struct {
	CardTypeID  int       `json:"card_type_id"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
	Message     string    `json:"message"`
	PushDate    time.Time `json:"push_date"`
}

// PushNotificationUpdate is used when updating an existing push notification.
type PushNotificationUpdate struct {
	Status      int       `json:"status"`
	Description string    `json:"description"`
	Message     string    `json:"message"`
	PushDate    time.Time `json:"push_date"`
}
