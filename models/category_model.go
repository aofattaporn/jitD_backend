package models

import (
	"cloud.google.com/go/firestore"
)

type Category struct {
	CategoryName string                 `json:"categoryName"`
	PostRef      *firestore.DocumentRef `json:"date,omitempty"`
}

type CategoryResponce struct {
	CategoryName string `json:"categoryName"`
	CountPost    int    `json:"countPost"`
}
