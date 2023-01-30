package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Post struct {
	UserID   *firestore.DocumentRef   `json:"userId"`
	Content  string                   `json:"content"`
	Date     time.Time                `json:"date,omitempty"`
	IsPublic bool                     `json:"isPublic"`
	Category []string                 `json:"category,omitempty"`
	Comment  []*firestore.DocumentRef `json:"comment,omitempty"`
	Like     []*firestore.DocumentRef `json:"like,omitempty"`
}

type PostResponse struct {
	UserId       string    `json:"userId"`
	PostId       string    `json:"postId"`
	Content      string    `json:"content"`
	Date         time.Time `json:"date,omitempty"`
	IsPublic     bool      `json:"isPublic"`
	Category     []string  `json:"category"`
	CountComment int       `json:"countComment"`
	CountLike    int       `json:"countLike"`
}
