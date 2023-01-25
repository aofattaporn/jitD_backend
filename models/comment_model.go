package models

import "google.golang.org/genproto/googleapis/type/date"

type Comment struct {
	Comment_id int       `json:"comment_id"`
	Content    string    `json:"content"`
	Like       string    `json:"like,omitempty"`
	Date       date.Date `json:"date"`
}
