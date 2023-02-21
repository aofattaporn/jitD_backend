package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Like struct {
	UserRef *firestore.DocumentRef `json:"userRef,omitempty"`
	PostRef *firestore.DocumentRef `json:"postRef,omitempty"`
	Date    time.Time              `json:"date,omitempty"`
}

type LikeComment struct {
	UserID    string    `json:"userID,omitempty"`
	CommentID string    `json:"commentID,omitempty"`
	Date      time.Time `json:"date,omitempty"`
}
