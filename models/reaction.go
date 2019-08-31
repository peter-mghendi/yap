package models

import (
	"net/http"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type reactionType int

// reactionTypes represent available user actions on posts.
const (
	ReactionApprove reactionType = iota + 1
	ReactionSticker
	ReactionComment
)

// Reaction represents a User action on a Post
type Reaction struct {
	Base
	Type reactionType `json:"type"`
	User uuid.UUID    `gorm:"type:uuid" json:"user"`
	Post uuid.UUID    `gorm:"type:uuid" json:"post"`
	Text string       `json:"text"`
}

// Create makes new reactions
func (r *Reaction) Create() (int, error) {
	if r.Type != ReactionComment {
		r.Text = ""
	}

	if err := DB.Create(r).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (r *Reaction) Read() (int, error) {
	if err := DB.First(r).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

// Update edits reaction text
func (r *Reaction) Update() (int, error) {
	if r.Type != ReactionComment {
		return http.StatusMethodNotAllowed, nil
	}

	if err := DB.Model(r).Updates(Reaction{Text: r.Text}).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes existing reactions
func (r *Reaction) Delete() (int, error) {
	db := DB.Delete(r)
	if num, err := db.RowsAffected, db.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusAccepted, nil
}

// ReadReactions fetches an array of reactions per specified criteria
// HACK Some bit hacking going on here
func ReadReactions(r *Reaction) ([]Reaction, int, error) {
	temp, reactions := Reaction{}, []Reaction{}

	if (r.Type & ReactionApprove) != 0 {
		temp.Type = ReactionApprove
		if err := DB.Where(&temp).Find(reactions).Error; err != nil {
			return nil, http.StatusInternalServerError, err
		}

		reactions = append(reactions, temp)
	}

	if (r.Type & ReactionSticker) != 0 {
		temp.Type = ReactionSticker
		if err := DB.Where(&temp).Find(reactions).Error; err != nil {
			return nil, http.StatusInternalServerError, err
		}

		reactions = append(reactions, temp)
	}

	if (r.Type & ReactionComment) != 0 {
		temp.Type = ReactionComment
		if err := DB.Where(&temp).Find(reactions).Error; err != nil {
			return nil, http.StatusInternalServerError, err
		}

		reactions = append(reactions, temp)
	}

	return reactions, http.StatusOK, nil
}
