package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/l3njo/yap-api/db"
)

// Article represents prose posts
type Article struct {
	PostBase
	Content string `json:"content"`
}

// Create makes an Article
func (a *Article) Create() (int, error) {
	article := Article{
		PostBase: PostBase{
			Subject: a.Subject,
			Summary: a.Summary,
			Overlay: a.Overlay,
			Section: a.Section,
			Creator: a.Creator,
			Markers: a.Markers,
		},
		Content: a.Content,
	}

	if err := db.DB.Create(&article).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	*a = article
	return http.StatusCreated, nil
}

// Read fetches an Article
func (a *Article) Read() (int, error) {
	if err := db.DB.Set("gorm:auto_preload", true).First(a).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	a.Summons++
	db.DB.Save(a)
	return http.StatusOK, nil
}

// Update edits an Article
func (a *Article) Update() (int, error) {
	article := Article{
		PostBase: PostBase{
			Subject: a.Subject,
			Summary: a.Summary,
			Overlay: a.Overlay,
			Section: a.Section,
			Markers: a.Markers,
		},
		Content: a.Content,
	}

	if err := db.DB.Model(a).Updates(article).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(a).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes an Article
func (a *Article) Delete() (int, error) {
	res := db.DB.Delete(a)
	if num, err := res.RowsAffected, res.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusAccepted, nil
}

// Publish makes an Article public
func (a *Article) Publish() (int, error) {
	if a.Release {
		return http.StatusNotModified, nil
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(a).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	a.Release, a.Summons = true, 0
	db.DB.Save(a)
	return http.StatusAccepted, nil
}

// Retract makes an Article private
func (a *Article) Retract() (int, error) {
	if !a.Release {
		return http.StatusNotModified, nil
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(a).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	a.Release, a.Summons = false, 0
	db.DB.Save(a)
	return http.StatusAccepted, nil
}
