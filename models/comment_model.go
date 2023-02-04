package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Comment struct {
	Content string                   `json:"content"`
	Like    []*firestore.DocumentRef `json:"like,omitempty"`
	Date    time.Time                `json:"date,omitempty"`
}

type CommentResponse struct {
	UserId       string    `json:"userId"`
	PostId       string    `json:"postId"`
	Comment_id   string    `json:"comment_id"`
	Content      string    `json:"Content"`
	CountComment int       `json:"countComment"`
	Like         string    `json:"like,omitempty"`
	Date         time.Time `json:"date,omitempty"`
}
