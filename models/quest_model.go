package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type DailyQuestProgress struct {
	UserID    *firestore.DocumentRef `json:"userId"`
	QuestDate time.Time              `json:"questDate"`
	Quests    []*Quest               `json:"quest"`
}

type Quest struct {
	QuestName      string    `json:"questName"`
	Progress       int       `json:"progress"`
	MaxProgress    int       `json:"maxProgress"`
	Reward         int       `json:"reward"`
	IsGetPoint     bool      `json:"isGetPoint"`
	Completed      bool      `json:"completed"`
	LastCompletion time.Time `json:"lastCompletion"`
}
