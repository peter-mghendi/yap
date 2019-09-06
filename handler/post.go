package handler

import (
	"net/http"

	"github.com/l3njo/yap-api/model"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo/v4"
)

// PostResponse is a response containing one Post
type PostResponse struct {
	Response
	model.Post `json:"data"`
}

// PublishPost handles the "/posts/:id/publish" route.
func PublishPost(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	post, status, err := model.GetPost(id)
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err = post.Publish()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), post
	return c.JSON(status, resp)
}

// RetractPost handles the "/posts/:id/retract" route.
func RetractPost(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	post, status, err := model.GetPost(id)
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err = post.Retract()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), post
	return c.JSON(status, resp)
}

// DeletePost handles the "/posts/:id/delete" route.
func DeletePost(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	post, status, err := model.GetPost(id)
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err = post.Delete()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message = true, http.StatusText(status)
	return c.JSON(status, resp)
}
