package models

import (
	"cloud.google.com/go/firestore"
)

type User struct {
	UserID  string `json:"UserID"`
	PetName string `json:"petName"`
	Point   int    `json:"point"`
	Config  struct {
		Noti bool `json:"noti"`
	} `json:"config"`
	Posts    []*firestore.DocumentRef `json:posts,omitempty"`
	Comments []string                 `json:"comments"`
	Likes    []string                 `json:"likes"`
}
