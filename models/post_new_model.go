package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Post struct {
	UserID   *firestore.DocumentRef `json:"userId"`
	Content  string                 `json:"content"`
	Date     time.Time              `json:"ddate,omitempty"`
	IsPublic bool                   `json:"isPublic"`
	Category []string               `json:"category,omitempty"`
	Comment  []*Comment2            `json:"comment,omitempty"`
	LikesRef []*Like                `json:"likesref,omitempty"`
}

type Comment2 struct {
	CommentID string         `json:"commentID"`
	Content   string         `json:"content"`
	UserId    string         `json:"userId"`
	Like      []*LikeComment `json:"like,omitempty"`
	Date      time.Time      `json:"date,omitempty"`
}
