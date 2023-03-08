package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Comment struct {
	Content string                 `json:"content"`
	UserId  *firestore.DocumentRef `json:"userId"`
	PostId  *firestore.DocumentRef `json:"postId"`
	Like    []*LikeComment         `json:"like,omitempty"`
	Date    time.Time              `json:"date,omitempty"`
	IsPin   bool                   `json:"isPin"`
}

type CommentResponse struct {
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	CommentId string    `json:"commentId"`
	Content   string    `json:"content"`
	CountLike int       `json:"countLike"`
	Date      time.Time `json:"date,omitempty"`
	IsLike    bool      `json:"isLike"`
	IsPin     bool      `json:"isPin"`
}
