package models

type Test struct {
	Number       int       `json:"number"`
	QuestionText string    `json:"questionText"`
	Text         [4]string `json:"text"`
}
