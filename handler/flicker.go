package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/l3njo/yap/model"
	"github.com/l3njo/yap/util"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// FlickersResponse is a response containing a slice of Flickers
type FlickersResponse struct {
	Response
	Flickers []model.Flicker `json:"data"`
}

// GetFlickers handles the "/blog/posts/flickers" route.
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

// GetPublicFlickers handles the "/blog/posts/flickers/public" route.
func GetPublicFlickers(c echo.Context) error {
	resp, status := FlickersResponse{}, 0
	flickers, status, err := model.ReadAllFlickers()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	flickers = util.FilterF(flickers, func(f model.Flicker) bool {
		return f.Release
	})

	resp.Status, resp.Message, resp.Flickers = true, http.StatusText(status), flickers
	return c.JSON(status, resp)
}

// GetFlickerByID handles the "/blog/posts/flickers/:id" route.
func GetFlickerByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	flicker := &model.Flicker{}
	flicker.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(flicker.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err := flicker.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), flicker
	return c.JSON(status, resp)
}

// GetPublicFlickerByID handles the "/blog/posts/flickers/public/:id" route.
func GetPublicFlickerByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	flicker := &model.Flicker{
		PostBase: model.PostBase{
			Base: model.Base{
				ID: uuid.FromStringOrNil(c.Param("id")),
			},
			Release: true,
		},
	}

	if uuid.Equal(flicker.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err := flicker.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), flicker
	return c.JSON(status, resp)
}

// CreateFlicker handles the "/blog/posts/flickers/create" route.
func CreateFlicker(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := PostResponse{}, 0
	flicker := &model.Flicker{}
	if err := c.Bind(flicker); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil) {
		status = http.StatusForbidden
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	flicker.Creator = claims.User
	if status, err := flicker.Create(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), flicker
	return c.JSON(status, resp)
}

// UpdateFlicker handles the "/blog/posts/flickers/:id/update" route.
func UpdateFlicker(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

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

	if flicker.Creator != claims.User {
		if (!flicker.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
			(flicker.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
			status = http.StatusForbidden
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}
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

// TransferFlicker handles the "/blog/posts/flickers/:id/transfer" route.
func TransferFlicker(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

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

	if flicker.Creator != claims.User {
		if (!flicker.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
			(flicker.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
			status = http.StatusForbidden
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}
	}

	if uuid.Equal(flicker.Creator, f.Creator) {
		status = http.StatusNotModified
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	flicker.Creator = f.Creator
	status, err := flicker.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), flicker
	return c.JSON(status, resp)
}
