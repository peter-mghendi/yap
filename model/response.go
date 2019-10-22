package model

import (
	uuid "github.com/satori/go.uuid"
)

// Response represents a forum Response
type Response struct {
	Base
	Content string    `json:"content"`
	Creator uuid.UUID `json:"creator" gorm:"type:uuid"`
	ReplyTo uuid.UUID `json:"replyTo" gorm:"type:uuid"`
}
