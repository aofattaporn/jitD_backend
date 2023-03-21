package models

import "cloud.google.com/go/firestore"

type User struct {
	PetName  string                   `json:"petName"`
	PetHP    int                      `json:"petHP"`
	Point    int                      `json:"point"`
	BookMark []*firestore.DocumentRef `json:"bookMark"`
}

type UserResponse struct {
	UserID   string   `json:"userID"`
	PetName  string   `json:"petName"`
	PetHP    int      `json:"petHP"`
	Point    int      `json:"point"`
	BookMark []string `json:"bookMark"`
}
