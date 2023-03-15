package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Test struct {
	Number       int       `json:"number"`
	QuestionText string    `json:"questionText"`
	Choices      [4]Choice `json:"choices"`
}

type Choice struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
	Value  int    `json:"value"`
}

type TestResualt struct {
	UserID   *firestore.DocumentRef `json:"userId"`
	TestDate time.Time              `json:"testDate"`
	TestName string                 `json:"testName"`
	Point    int                    `json:"point"`
	Result   string                 `json:"result"`
	Desc     string                 `json:"desc"`
}

type TestResualtResponse struct {
	TestDate time.Time `json:"testDate"`
	TestName string    `json:"testName"`
	Point    int       `json:"point"`
	Result   string    `json:"result"`
	Desc     string    `json:"desc"`
}
