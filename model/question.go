package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/l3njo/yap/db"
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
	Answers pq.StringArray `json:"answers" gorm:"type:varchar(255)[]"`
}

// Create makes a new question
func (q *Question) Create() (int, error) {
	question := Question{
		Subject: q.Subject,
		Content: q.Content,
		Markers: q.Markers,
		Section: q.Section,
		Creator: q.Creator,
	}

	if err := db.DB.Create(&question).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	*q = question
	return http.StatusCreated, nil
}

// Read returns an existing question
func (q *Question) Read() (int, error) {
	if err := db.DB.First(q).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	q.Summons++
	db.DB.Save(q)
	return http.StatusOK, nil
}

// Update edits an existing question.
func (q *Question) Update() (int, error) {
	question := Question{
		Subject: q.Subject,
		Content: q.Content,
		Markers: q.Markers,
		Section: q.Section,
		Answers: q.Answers,
	}

	err := db.DB.Model(q).Updates(question).Error
	if gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(q).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes an existing question.
func (q *Question) Delete() (int, error) {
	db := db.DB.Delete(q)
	if num, err := db.RowsAffected, db.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}
	return http.StatusAccepted, nil
}

// ReadAllQuestions fetches all questions.
func ReadAllQuestions() ([]Question, int, error) {
	questions := []Question{}
	if err := db.DB.Set("gorm:auto_preload", true).Find(&questions).Error; gorm.IsRecordNotFoundError(err) {
		return questions, http.StatusNotFound, err
	}

	return questions, http.StatusOK, nil
}
