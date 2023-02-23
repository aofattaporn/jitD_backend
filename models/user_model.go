package models

type User struct {
	PetName       string              `json:"petName"`
	PetHP         int                 `json:"petHP"`
	Point         int                 `json:"point"`
	HistorySearch []string            `json:"historySearch"`
	DailyQuests   *DailyQuestProgress `json:"dailyQuests,omitempty"`
}

type UserResponse struct {
	UserID string `json:"userID"`
	// data same a request
	PetName       string              `json:"petName"`
	Point         int                 `json:"point"`
	PetHP         int                 `json:"petHP"`
	HistorySearch []string            `json:"historySearch"`
	DailyQuests   *DailyQuestProgress `json:"dailyQuests,omitempty"`
}
