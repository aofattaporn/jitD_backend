package models

import (
	"time"
)

type DailyQuestProgress struct {
	QuestDate time.Time `json:"questTime"`
	Quests    []*Quest  `json:"quest"`
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
