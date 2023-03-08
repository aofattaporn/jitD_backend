package models

import (
	"time"
)

type Comment2 struct {
	CommentID string         `json:"commentID"`
	Content   string         `json:"content"`
	UserId    string         `json:"userId"`
	Like      []*LikeComment `json:"like,omitempty"`
	Date      time.Time      `json:"date,omitempty"`
	IsPin     bool           `json:"isPin"`
}
