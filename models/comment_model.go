package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Comment struct {
	Content string                   `json:"content"`
	UserId  string                   `json:"userId"`
	PostId  string                   `json:"postId"`
	Like    []*firestore.DocumentRef `json:"like,omitempty"`
	Date    time.Time                `json:"date,omitempty"`
}

type CommentResponse struct {
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	CommentId string    `json:"commentId"`
	Content   string    `json:"content"`
	CountLike int       `json:"countLike"`
	Date      time.Time `json:"date,omitempty"`
	IsPublic  bool      `json:"isPublic"`
}
