package models

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

// Post represents all post types
type Post interface {
	// Create makes a Post
	Create() (int, error)
	// Read fetches a Post
	Read() (int, error)
	// Update edits a Post
	Update() (int, error)
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
	Pattern   string         `json:"pattern"`
	Release   bool           `json:"release"`
	Creator   uuid.UUID      `json:"creator" gorm:"type:uuid"`
	Markers   pq.StringArray `json:"markers" gorm:"type:varchar(255)[]"`
	Reactions []Reaction     `json:"reactions,omitempty" sql:"-" gorm:"foreignkey:Post"`
}

// Article represents prose posts
type Article struct {
	PostBase
	Content string
}

// Gallery represents image posts
type Gallery struct {
	PostBase
	Content pq.StringArray `json:"content" gorm:"type:varchar(255)[]"`
	Caption pq.StringArray `json:"caption" gorm:"type:varchar(255)[]"`
}

// Flicker represents video posts
type Flicker struct {
	PostBase
	Content string `json:"content"`
	Caption string `json:"caption"`
}

// Create makes an Article
func (a *Article) Create() (int, error) {
	a.Release = false
	if err := DB.Create(a).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// Read fetches an Article
func (a *Article) Read() (int, error) {
	if err := DB.Preload(strReactions).First(a).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	a.Summons++
	DB.Save(a)
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

	if err := DB.Model(a).Updates(article).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes an Article
func (a *Article) Delete() (int, error) {
	db := DB.Delete(a)
	if num, err := db.RowsAffected, db.Error; num == 0 {
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

	if err := DB.First(a).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	a.Release, a.Summons = true, 0
	DB.Save(a)
	return http.StatusAccepted, nil
}

// Retract makes an Article private
func (a *Article) Retract() (int, error) {
	if !a.Release {
		return http.StatusNotModified, nil
	}

	if err := DB.First(a).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	a.Release, a.Summons = false, 0
	DB.Save(a)
	return http.StatusAccepted, nil
}

// Create makes a Gallery
func (g *Gallery) Create() (int, error) {
	g.Release = false
	if err := DB.Create(g).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// Read fetches a Gallery
func (g *Gallery) Read() (int, error) {
	if err := DB.Preload(strReactions).First(g).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	g.Summons++
	DB.Save(g)
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

	if err := DB.Model(g).Updates(gallery).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes a Gallery
func (g *Gallery) Delete() (int, error) {
	db := DB.Delete(g)
	if num, err := db.RowsAffected, db.Error; num == 0 {
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

	if err := DB.First(g).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	g.Release, g.Summons = true, 0
	DB.Save(g)
	return http.StatusAccepted, nil
}

// Retract makes makes a Gallery private
func (g *Gallery) Retract() (int, error) {
	if !g.Release {
		return http.StatusNotModified, nil
	}

	if err := DB.First(g).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	g.Release, g.Summons = false, 0
	DB.Save(g)
	return http.StatusAccepted, nil
}

// Create makes a Flicker
func (f *Flicker) Create() (int, error) {
	f.Release = false
	if err := DB.Create(f).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

// Read fetches a Flicker
func (f *Flicker) Read() (int, error) {
	if err := DB.Preload(strReactions).First(f).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	}

	f.Summons++
	DB.Save(f)
	return http.StatusOK, nil
}

// Update edits a Flicker
func (f *Flicker) Update() (int, error) {
	flicker := Flicker{
		PostBase: PostBase{
			Subject: f.Subject,
			Summary: f.Summary,
			Overlay: f.Overlay,
			Section: f.Section,
			Markers: f.Markers,
		},
		Content: f.Content,
		Caption: f.Caption,
	}

	if err := DB.Model(f).Updates(flicker).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

// Delete removes a Flicker
func (f *Flicker) Delete() (int, error) {
	db := DB.Delete(f)
	if num, err := db.RowsAffected, db.Error; num == 0 {
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

	if err := DB.First(f).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	f.Release, f.Summons = true, 0
	DB.Save(f)
	return http.StatusAccepted, nil
}

// Retract makes makes a Flicker private
func (f *Flicker) Retract() (int, error) {
	if !f.Release {
		return http.StatusNotModified, nil
	}

	if err := DB.First(f).Error; gorm.IsRecordNotFoundError(err) {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	f.Release, f.Summons = true, 0
	DB.Save(f)
	return http.StatusAccepted, nil
}
