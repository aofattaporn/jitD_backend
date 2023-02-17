package models

import "time"

type DailyQuestProgress struct {
	QuestName      string    `json:"questName"`
	Progress       int       `json:"progress"`
	MaxProgress    int       `json:"maxProgress"`
	Reward         int       `json:"reward"`
	Completed      bool      `json:"completed"`
	LastCompletion time.Time `json:"lastCompletion"`
}
