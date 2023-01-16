package models

import "time"

type Post struct {
	Content  string    `json:"content"`
	Date     time.Time `firestore:"updated,omitempty"`
	Comment  []string  `json:"comment,omitempty"`
	Like     []string  `json:"like,omitempty"`
	IsPublic bool      `json:"IsPublic"`
}
