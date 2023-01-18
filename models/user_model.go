package models

import (
	"cloud.google.com/go/firestore"
)

type User struct {
	UserID  string `json:"UserID"`
	PetName string `json:"petName"`
	Point   int    `json:"point"`

	// Reference collection
	Posts    []*firestore.DocumentRef `json:posts,omitempty"`
	Comments []*firestore.DocumentRef `json:"comments,omitempty"`
	Likes    []*firestore.DocumentRef `json:"likes,omitempty"`
}
