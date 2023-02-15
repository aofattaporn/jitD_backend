package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Like struct {
	User        User                   `json:"user"`
	Type        string                 `json:"type"`
	DocumentRef *firestore.DocumentRef `json:"documentRef,omitempty"`
	Date        time.Time              `json:"date,omitempty"`
}
