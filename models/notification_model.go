package models

import (
	"time"
)

type Notification struct {
	Message         string    `firestore:"message"`
	Timestamp       time.Time `firestore:"timestamp"`
	RecipientUserID string    `firestore:"recipient_user_id"`
	SenderUserID    string    `firestore:"sender_user_id"`
	Read            bool      `firestore:"read"`
}
