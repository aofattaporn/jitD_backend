package models

import "cloud.google.com/go/firestore"

type User struct {
	UserID   string                   `json:"userID"`
	PetName  string                   `json:"petName"`
	PetHP    int                      `json:"petHP"`
	Point    int                      `json:"point"`
	IsAdmin  bool                     `json:"isAdmin"`
	BookMark []*firestore.DocumentRef `json:"bookMark"`
}

type UserResponse struct {
	UserID   string   `json:"userID"`
	PetName  string   `json:"petName"`
	PetHP    int      `json:"petHP"`
	Point    int      `json:"point"`
	IsAdmin  bool     `json:"isAdmin"`
	BookMark []string `json:"bookMark"`
}
