package models

import "cloud.google.com/go/firestore"

type User struct {
	PetName       string                   `json:"petName"`
	PetHP         int                      `json:"petHP"`
	Point         int                      `json:"point"`
	HistorySearch []string                 `json:"historySearch"`
	BookMark      []*firestore.DocumentRef `json:"bookMark"`
}

type UserResponse struct {
	UserID string `json:"userID"`
	// data same a request
	PetName       string   `json:"petName"`
	PetHP         int      `json:"petHP"`
	Point         int      `json:"point"`
	HistorySearch []string `json:"historySearch"`
	BookMark      []string `json:"bookMark"`
}
