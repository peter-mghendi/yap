package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/l3njo/yap/db"
	uuid "github.com/satori/go.uuid"
)

// Response represents a forum Response
type Response struct {
	Base
	Content string    `json:"content"`
	Creator uuid.UUID `json:"creator" gorm:"type:uuid"`
	Summons int       `json:"summons"`
	ReplyTo uuid.UUID `json:"replyTo" gorm:"type:uuid"`
}

// Create makes a new Response
func (r *Response) Create() (int, error) {
	response := Response{
		Content: r.Content,
		Creator: r.Creator,
		ReplyTo: r.ReplyTo,
	}

	if err := db.DB.Create(&response).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	*r = response
	return http.StatusCreated, nil
}

// Read returns an existing Response
func (r *Response) Read() (int, error) {
	if err := db.DB.First(r).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	r.Summons++
	db.DB.Save(r)
	return http.StatusOK, nil
}

// Update edits an existing Response.
func (r *Response) Update() (int, error) {
	err := db.DB.Model(r).Updates(Response{Content: r.Content}).Error
	if gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(r).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes an existing Response.
func (r *Response) Delete() (int, error) {
	db := db.DB.Delete(r)
	if num, err := db.RowsAffected, db.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}
	return http.StatusAccepted, nil
}

// ReadAllResponses fetches all Responses.
func ReadAllResponses() ([]Response, int, error) {
	responses := []Response{}
	if err := db.DB.Set("gorm:auto_preload", true).Find(&responses).Error; gorm.IsRecordNotFoundError(err) {
		return responses, http.StatusNotFound, err
	}

	return responses, http.StatusOK, nil
}
