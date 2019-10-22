package model

import (
	"net/http"

	"github.com/l3njo/yap/db"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// ReactionType represents a user action on a post.
type ReactionType string

// ReactionTypes represent available user actions on posts.
const (
	ReactionApprove ReactionType = "approve"
	ReactionSticker ReactionType = "sticker"
	ReactionComment ReactionType = "comment"
)

// Reaction represents a User action on a Post
type Reaction struct {
	Base
	Type ReactionType `json:"type"`
	User uuid.UUID    `gorm:"type:uuid" json:"user"`
	Post uuid.UUID    `gorm:"type:uuid" json:"post"`
	Site string       `json:"site"`
	Text string       `json:"text"`
}

// Create makes new reactions
func (r *Reaction) Create() (int, error) {
	if r.Type != ReactionComment {
		r.Text = ""
	}

	if err := db.DB.Create(r).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (r *Reaction) Read() (int, error) {
	if err := db.DB.First(r).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

// Update edits reaction text
func (r *Reaction) Update() (int, error) {
	if r.Type != ReactionComment {
		return http.StatusMethodNotAllowed, nil
	}

	err := db.DB.Model(r).Updates(Reaction{Text: r.Text}).Error
	if gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes existing reactions
func (r *Reaction) Delete() (int, error) {
	db := db.DB.Delete(r)
	if num, err := db.RowsAffected, db.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusAccepted, nil
}

// ReadAllReactions fetches all Reactions
func ReadAllReactions() ([]Reaction, int, error) {
	reactions := []Reaction{}
	if err := db.DB.Find(&reactions).Error; gorm.IsRecordNotFoundError(err) {
		return reactions, http.StatusNotFound, err
	}

	return reactions, http.StatusOK, nil
}
