package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/l3njo/yap/model"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo/v4"
)

// PostResponse is a response containing one Post
type PostResponse struct {
	Response
	model.Post `json:"data"`
}

// PublishPost handles the "/blog/posts/:id/publish" route.
func PublishPost(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil) {
		status = http.StatusForbidden
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

// RetractPost handles the "/blog/posts/:id/retract" route.
func RetractPost(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil) {
		status = http.StatusForbidden
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

// DeletePost handles the "/blog/posts/:id/delete" route.
func DeletePost(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

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

	switch v := post.(type) {
	case *model.Article:
		if v.Creator != claims.User {
			if (!v.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
				(v.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
				status = http.StatusForbidden
				resp.Message = http.StatusText(status)
				return c.JSON(status, resp)
			}
		}
	case *model.Gallery:
		if v.Creator != claims.User {
			if (!v.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
				(v.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
				status = http.StatusForbidden
				resp.Message = http.StatusText(status)
				return c.JSON(status, resp)
			}
		}
	case *model.Flicker:
		if v.Creator != claims.User {
			if (!v.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
				(v.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
				status = http.StatusForbidden
				resp.Message = http.StatusText(status)
				return c.JSON(status, resp)
			}
		}
	default:
		status = http.StatusInternalServerError
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
