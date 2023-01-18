package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Post struct {
	Content  string                   `json:"content"`
	Date     time.Time                `json:"Date,omitempty"`
	IsPublic bool                     `json:"IsPublic"`
	Category []string                 `json:"Category,omitempty"`
	Comment  []*firestore.DocumentRef `json:"comment,omitempty"`
	Like     []*firestore.DocumentRef `json:"like,omitempty"`
}
