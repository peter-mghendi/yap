package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/l3njo/yap-api/db"
	"github.com/lib/pq"
)

// Gallery represents image posts
type Gallery struct {
	PostBase
	Content pq.StringArray `json:"content" gorm:"type:varchar(255)[]"`
	Caption pq.StringArray `json:"caption" gorm:"type:varchar(255)[]"`
}

// Create makes a Gallery
func (g *Gallery) Create() (int, error) {
	gallery := Gallery{
		PostBase: PostBase{
			Subject: g.Subject,
			Summary: g.Summary,
			Overlay: g.Overlay,
			Section: g.Section,
			Creator: g.Creator,
			Markers: g.Markers,
		},
		Content: g.Content,
		Caption: g.Caption,
	}

	if err := db.DB.Create(&gallery).Error; err != nil {
		return http.StatusInternalServerError, err
	}

    *g = gallery
	return http.StatusCreated, nil
}

// Read fetches a Gallery
func (g *Gallery) Read() (int, error) {
	if err := db.DB.Set("gorm:auto_preload", true).First(g).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	g.Summons++
	db.DB.Save(g)
	return http.StatusOK, nil
}

// Update edits a Gallery
func (g *Gallery) Update() (int, error) {
	gallery := Gallery{
		PostBase: PostBase{
			Subject: g.Subject,
			Summary: g.Summary,
			Overlay: g.Overlay,
			Section: g.Section,
			Markers: g.Markers,
		},
		Content: g.Content,
		Caption: g.Caption,
	}

	if err := db.DB.Model(g).Updates(gallery).Error; gorm.IsRecordNotFoundError(err) {
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

// Delete removes a Gallery
func (g *Gallery) Delete() (int, error) {
	res := db.DB.Delete(g)
	if num, err := res.RowsAffected, res.Error; num == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusAccepted, nil
}

// Publish makes a Gallery public
func (g *Gallery) Publish() (int, error) {
	if g.Release {
		return http.StatusNotModified, nil
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(g).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	g.Release, g.Summons = true, 0
	db.DB.Save(g)
	return http.StatusAccepted, nil
}

// Retract makes makes a Gallery private
func (g *Gallery) Retract() (int, error) {
	if !g.Release {
		return http.StatusNotModified, nil
	}

	if err := db.DB.Set("gorm:auto_preload", true).First(g).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	g.Release, g.Summons = false, 0
	db.DB.Save(g)
	return http.StatusAccepted, nil
}
