package models

type Post struct {
	Post_id int    `json:"post_id"`
	Content string `json:"content"`
	Comment string `json:"comment,omitempty"`
	Like    string `json:"like,omitempty"`
}
