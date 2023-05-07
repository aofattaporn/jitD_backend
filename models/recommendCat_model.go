package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type CategoryReccommend struct {
	UserID              *firestore.DocumentRef `json:"userId"`
	Date                time.Time              `json:"ddate,omitempty"`
	CountStudy          int                    `json:"countStudy"`
	CountWork           int                    `json:"countWork"`
	CountMentalHealth   int                    `json:"countMentalHealth"`
	CountLifeProblem    int                    `json:"countLifeProblem"`
	CountRelationship   int                    `json:"countRelationship"`
	CountFamily         int                    `json:"countFamily"`
	CountPhysicalHealth int                    `json:"countPhysicalHealth"`
}
