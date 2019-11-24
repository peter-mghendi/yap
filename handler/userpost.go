package handler

import (
	"net/http"

	"github.com/l3njo/yap/model"
	"github.com/l3njo/yap/util"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// GetUserPublicBlogArticles handles the "/users/:id/blog/posts/articles" route.
func GetUserPublicBlogArticles(c echo.Context) error {
	resp, status := ArticlesResponse{}, 0
	articles, status, err := model.ReadAllArticles()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	userID := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(userID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	articles = util.FilterA(articles, func(a model.Article) bool {
		return (uuid.Equal(a.Creator, userID)) && a.Release
	})

	resp.Status, resp.Message, resp.Articles = true, http.StatusText(status), articles
	return c.JSON(status, resp)
}

// GetUserPublicBlogGalleries handles the "/users/:id/blog/posts/galleries" route.
func GetUserPublicBlogGalleries(c echo.Context) error {
	resp, status := GalleriesResponse{}, 0
	galleries, status, err := model.ReadAllGalleries()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	userID := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(userID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	galleries = util.FilterG(galleries, func(g model.Gallery) bool {
		return (uuid.Equal(g.Creator, userID)) && g.Release
	})

	resp.Status, resp.Message, resp.Galleries = true, http.StatusText(status), galleries
	return c.JSON(status, resp)
}

// GetUserPublicBlogFlickers handles the "/users/:id/blog/posts/flickers" route.
func GetUserPublicBlogFlickers(c echo.Context) error {
	resp, status := FlickersResponse{}, 0
	flickers, status, err := model.ReadAllFlickers()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	userID := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(userID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	flickers = util.FilterF(flickers, func(f model.Flicker) bool {
		return (uuid.Equal(f.Creator, userID)) && f.Release
	})

	resp.Status, resp.Message, resp.Flickers = true, http.StatusText(status), flickers
	return c.JSON(status, resp)
}
