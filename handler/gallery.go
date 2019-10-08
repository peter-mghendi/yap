package handler

import (
	"net/http"

	"github.com/l3njo/yap/model"
	"github.com/l3njo/yap/util"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// GalleriesResponse is a response containing a slice of Galleries
type GalleriesResponse struct {
	Response
	Galleries []model.Gallery `json:"data"`
}

// GetGalleries handles the "/posts/galleries" route.
func GetGalleries(c echo.Context) error {
	resp, status := GalleriesResponse{}, 0
	galleries, status, err := model.ReadAllGalleries()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Galleries = true, http.StatusText(status), galleries
	return c.JSON(status, resp)
}

// GetPublicGalleries handles the "/posts/galleries/public" route.
func GetPublicGalleries(c echo.Context) error {
	resp, status := GalleriesResponse{}, 0
	galleries, status, err := model.ReadAllGalleries()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	galleries = util.FilterG(galleries, func(g model.Gallery) bool {
		return g.Release
	})

	resp.Status, resp.Message, resp.Galleries = true, http.StatusText(status), galleries
	return c.JSON(status, resp)
}

// GetGalleryByID handles the "/posts/galleries/:id" route.
func GetGalleryByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	gallery := &model.Gallery{
		PostBase: model.PostBase{
			Base: model.Base{ID: id},
		},
	}

	status, err := gallery.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}

// GetPublicGalleryByID handles the "/posts/galleries/public/:id" route.
func GetPublicGalleryByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	gallery := &model.Gallery{
		PostBase: model.PostBase{
			Base: model.Base{ID: id},
		},
	}

	status, err := gallery.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !gallery.Release {
		status = http.StatusNotFound
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}

// CreateGallery handles the "/posts/galleries/create" route.
func CreateGallery(c echo.Context) error {
	resp, status := PostResponse{}, 0
	gallery := &model.Gallery{}
	if err := c.Bind(gallery); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := gallery.Create(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}

// UpdateGallery handles the "/posts/galleries/:id/update" route.
func UpdateGallery(c echo.Context) error {
	resp, status := PostResponse{}, 0
	gallery, g := &model.Gallery{}, &model.Gallery{}
	if err := c.Bind(g); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	gallery.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(gallery.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := gallery.Read(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	gallery.PostBase, gallery.Content, gallery.Caption = g.PostBase, g.Content, g.Caption
	status, err := gallery.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}
