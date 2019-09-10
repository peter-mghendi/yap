package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/l3njo/yap-api/db"
)

// Flicker represents video posts
type Flicker struct {
	PostBase
	Content string `json:"content"`
	Caption string `json:"caption"`
}

// Create makes a Flicker
func (f *Flicker) Create() (int, error) {
	flicker := Flicker{
		PostBase: PostBase{
			Subject: f.Subject,
			Summary: f.Summary,
			Overlay: f.Overlay,
			Section: f.Section,
			Creator: f.Creator,
			Markers: f.Markers,
		},
		Content: f.Content,
		Caption: f.Caption,
	}

	if err := db.DB.Create(&flicker).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	*f = flicker
	return http.StatusCreated, nil
}

// Read fetches a Flicker
func (f *Flicker) Read() (int, error) {
	if err := db.DB.Set("gorm:auto_preload", true).First(f).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	f.Summons++
	db.DB.Save(f)
	return http.StatusOK, nil
}

// Update edits a Flicker
func (f *Flicker) Update() (int, error) {
	flicker := Flicker{
		PostBase: PostBase{
			Subject: f.Subject,
			Summary: f.Summary,
			Overlay: f.Overlay,
			Pattern: flickerPost,
			Section: f.Section,
			Markers: f.Markers,
		},
		Content: f.Content,
		Caption: f.Caption,
	}

	if err := db.DB.Model(f).Updates(flicker).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := db.DB.Set("gorm:auto_preload", true).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes a Flicker
func (f *Flicker) Delete() (int, error) {
	res := db.DB.Delete(f)
	if num, err := res.RowsAffected, res.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusAccepted, nil
}

// Publish makes a Flicker public
func (f *Flicker) Publish() (int, error) {
	if f.Release {
		return http.StatusNotModified, nil
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(f).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	f.Release, f.Summons = true, 0
	db.DB.Save(f)
	return http.StatusAccepted, nil
}

// Retract makes makes a Flicker private
func (f *Flicker) Retract() (int, error) {
	if !f.Release {
		return http.StatusNotModified, nil
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(f).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	f.Release, f.Summons = true, 0
	db.DB.Save(f)
	return http.StatusAccepted, nil
}

// ReadAllFlickers fetches all Flickers
func ReadAllFlickers() ([]Flicker, int, error) {
	flickers := []Flicker{}
	if err := db.DB.Set("gorm:auto_preload", true).Find(&flickers).Error; gorm.IsRecordNotFoundError(err) {
		return flickers, http.StatusNotFound, err
	}

	return flickers, http.StatusOK, nil
}
