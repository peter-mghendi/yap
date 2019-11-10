package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
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

// GetGalleries handles the "/blog/posts/galleries" route.
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

// GetPublicGalleries handles the "/blog/posts/galleries/public" route.
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

// GetGalleryByID handles the "/blog/posts/galleries/:id" route.
func GetGalleryByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	gallery := &model.Gallery{}
	gallery.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(gallery.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err := gallery.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}

// GetPublicGalleryByID handles the "/blog/posts/galleries/public/:id" route.
func GetPublicGalleryByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	gallery := &model.Gallery{
		PostBase: model.PostBase{
			Base: model.Base{
				ID: uuid.FromStringOrNil(c.Param("id")),
			},
			Release: true,
		},
	}

	if uuid.Equal(gallery.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err := gallery.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}

// CreateGallery handles the "/blog/posts/galleries/create" route.
func CreateGallery(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := PostResponse{}, 0
	gallery := &model.Gallery{}
	if err := c.Bind(gallery); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil) {
		status = http.StatusForbidden
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	gallery.Creator = claims.User
	if status, err := gallery.Create(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}

// UpdateGallery handles the "/blog/posts/galleries/:id/update" route.
func UpdateGallery(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

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

	if gallery.Creator != claims.User {
		if (!gallery.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
			(gallery.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
			status = http.StatusForbidden
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}
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

// TransferGallery handles the "/blog/posts/galleries/:id/transfer" route.
func TransferGallery(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

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

	if gallery.Creator != claims.User {
		if (!gallery.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
			(gallery.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
			status = http.StatusForbidden
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}
	}

	if uuid.Equal(gallery.Creator, g.Creator) {
		status = http.StatusNotModified
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	gallery.Creator = g.Creator
	status, err := gallery.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), gallery
	return c.JSON(status, resp)
}
