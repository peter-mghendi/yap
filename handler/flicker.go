package handler

import (
	"net/http"

	"github.com/l3njo/yap-api/model"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// FlickersResponse is a response containing a slice of Flickers
type FlickersResponse struct {
	Response
	Flickers []model.Flicker `json:"data"`
}

// GetFlickers handles the "/posts/flickers" route.
func GetFlickers(c echo.Context) error {
	resp, status := FlickersResponse{}, 0
	flickers, status, err := model.ReadAllFlickers()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Flickers = true, http.StatusText(status), flickers
	return c.JSON(status, resp)
}

// GetFlickerByID handles the "/posts/flickers/:id" route.
func GetFlickerByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	flicker := &model.Flicker{
		PostBase: model.PostBase{
			Base: model.Base{ID: id},
		},
	}

	status, err := flicker.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), flicker
	return c.JSON(status, resp)
}

// CreateFlicker handles the "/posts/flickers/create" route.
func CreateFlicker(c echo.Context) error {
	resp, status := PostResponse{}, 0
	flicker := &model.Flicker{}
	if err := c.Bind(flicker); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := flicker.Create(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), flicker
	return c.JSON(status, resp)
}

// UpdateFlicker handles the "/posts/flickers/:id/update" route.
func UpdateFlicker(c echo.Context) error {
	resp, status := PostResponse{}, 0
	flicker, f := &model.Flicker{}, &model.Flicker{}
	if err := c.Bind(f); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	flicker.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(flicker.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := flicker.Read(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	flicker.PostBase, flicker.Content, flicker.Caption = f.PostBase, f.Content, f.Caption
	status, err := flicker.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), flicker
	return c.JSON(status, resp)
}