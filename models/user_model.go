package models

import (
	"cloud.google.com/go/firestore"
)

type User struct {
	PetName       string   `json:"petName"`
	Point         int      `json:"point"`
	HistorySearch []string `json:"historySearch"`

	// Reference collection
	DailyQuests map[string]*DailyQuestProgress `json:"dailyQuests,omitempty"`
	Posts       []*firestore.DocumentRef       `json:"posts,omitempty"`
	Comments    []*firestore.DocumentRef       `json:"comments,omitempty"`
	Likes       []*firestore.DocumentRef       `json:"likes,omitempty"`
}

type UserResponse struct {
	UserId        string `json:"userId"`
	PetName       string `json:"petName"`
	Point         int    `json:"point"`
	CountPosts    int    `json:"countPosts"`
	CountComments int    `json:"countComments"`
	CountLikes    int    `json:"countLikes"`
}
