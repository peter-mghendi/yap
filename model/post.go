package model

import (
	"errors"
	"net/http"

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

type postPattern string

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

	status := http.StatusNotFound
	return nil, status, errors.New(http.StatusText(status))
}
