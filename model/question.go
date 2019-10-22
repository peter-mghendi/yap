package model

import (
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

// Question represents a forum Question
type Question struct {
	Base
	Subject string         `json:"subject"`
	Content string         `json:"content"`
	Markers pq.StringArray `json:"markers" gorm:"type:varchar(255)[]"`
	Section string         `json:"section"`
	Creator uuid.UUID      `json:"creator" gorm:"type:uuid"`
	Summons int            `json:"summons"`
	Answers pq.Int64Array  `json:"answers" gorm:"type:varchar(255)[]"`
}
