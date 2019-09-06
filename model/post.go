package model

import (
	"net/http"

	"github.com/l3njo/yap-api/db"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

// Post represents all post types
type Post interface {
	// Read fetches a Post
	Read() (int, error)
	// Delete removes a Post
	Delete() (int, error)
	// Publish makes a Post public
	Publish() (int, error)
	// Retract makes a Post private
	Retract() (int, error)
}

// PostBase is the underlying object for all post types
type PostBase struct {
	Base
	Subject   string         `json:"subject"`
	Summary   string         `json:"summary"`
	Overlay   string         `json:"overlay"`
	Section   string         `json:"section"`
	Summons   int            `json:"summons"`
	Release   bool           `json:"release"`
	Creator   uuid.UUID      `json:"creator" gorm:"type:uuid"`
	Markers   pq.StringArray `json:"markers" gorm:"type:varchar(255)[]"`
	Reactions []Reaction     `json:"reactions,omitempty" sql:"-" gorm:"foreignkey:Post"`
}

// ReadAllArticles fetches all Articles
func ReadAllArticles() ([]Article, int, error) {
	articles := []Article{}
	if err := db.DB.Set("gorm:auto_preload", true).Find(&articles).Error; gorm.IsRecordNotFoundError(err) {
		return articles, http.StatusNotFound, err
	}

	return articles, http.StatusOK, nil
}

// ReadAllGalleries fetches all Galleries
func ReadAllGalleries() ([]Gallery, int, error) {
	galleries := []Gallery{}
	if err := db.DB.Set("gorm:auto_preload", true).Find(&galleries).Error; gorm.IsRecordNotFoundError(err) {
		return galleries, http.StatusNotFound, err
	}

	return galleries, http.StatusOK, nil
}

// ReadAllFlickers fetches all Flickers
func ReadAllFlickers() ([]Flicker, int, error) {
	flickers := []Flicker{}
	if err := db.DB.Set("gorm:auto_preload", true).Find(&flickers).Error; gorm.IsRecordNotFoundError(err) {
		return flickers, http.StatusNotFound, err
	}

	return flickers, http.StatusOK, nil
}

// GetPost finds a Post across Post models.
func GetPost(id uuid.UUID) (Post, int, error) {
	article := &Article{PostBase: PostBase{Base: Base{ID: id}}}
	if status, err := article.Read(); err == nil {
		return article, status, nil
	} else if status != http.StatusNotFound {
		return article, status, err
	}

	gallery := &Gallery{PostBase: PostBase{Base: Base{ID: id}}}
	if status, err := gallery.Read(); err == nil {
		return gallery, status, nil
	} else if status != http.StatusNotFound {
		return gallery, status, err
	}

	flicker := &Flicker{PostBase: PostBase{Base: Base{ID: id}}}
	if status, err := flicker.Read(); err == nil {
		return flicker, status, nil
	} else if status != http.StatusNotFound {
		return flicker, status, nil
	}

	return article, http.StatusNotFound, nil
}
