package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type User struct {
	UserID       string                   `json:"userID"`
	PetName      string                   `json:"petName"`
	PetHP        int                      `json:"petHP"`
	Point        int                      `json:"point"`
	IsAdmin      bool                     `json:"isAdmin"`
	FCMToken     string                   `json:"fcmToken"`
	BookMark     []*firestore.DocumentRef `json:"bookMark"`
	Notification []*Notification          `json:"notification"`
	RegisterDate time.Time                `json:"registerDate,omitempty"`
}

type UserResponse struct {
	UserID   string   `json:"userID"`
	PetName  string   `json:"petName"`
	PetHP    int      `json:"petHP"`
	Point    int      `json:"point"`
	IsAdmin  bool     `json:"isAdmin"`
	BookMark []string `json:"bookMark"`
	// FCMToken     string    `json:"fcmToken"`
	Notification []string  `json:"notification"`
	RegisterDate time.Time `json:"registerDate,omitempty"`
}
