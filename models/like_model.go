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
	UserRef    *firestore.DocumentRef `json:"userRef,omitempty"`
	CommentRef *firestore.DocumentRef `json:"commentRef,omitempty"`
	Date       time.Time              `json:"date,omitempty"`
}
