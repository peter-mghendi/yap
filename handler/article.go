package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/l3njo/yap/model"
	"github.com/l3njo/yap/util"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// ArticlesResponse is a response containing a slice of Articles
type ArticlesResponse struct {
	Response
	Articles []model.Article `json:"data"`
}

// GetArticles handles the "/blog/posts/articles" route.
func GetArticles(c echo.Context) error {
	resp, status := ArticlesResponse{}, 0
	articles, status, err := model.ReadAllArticles()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Articles = true, http.StatusText(status), articles
	return c.JSON(status, resp)
}

// GetPublicArticles handles the "/blog/posts/articles/public" route.
func GetPublicArticles(c echo.Context) error {
	resp, status := ArticlesResponse{}, 0
	articles, status, err := model.ReadAllArticles()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	articles = util.FilterA(articles, func(a model.Article) bool {
		return a.Release
	})

	resp.Status, resp.Message, resp.Articles = true, http.StatusText(status), articles
	return c.JSON(status, resp)
}

// GetArticleByID handles the "/blog/posts/articles/:id" route.
func GetArticleByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	article := &model.Article{
		PostBase: model.PostBase{
			Base: model.Base{ID: id},
		},
	}

	status, err := article.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), article
	return c.JSON(status, resp)
}

// GetPublicArticleByID handles the "/blog/posts/articles/public/:id" route.
func GetPublicArticleByID(c echo.Context) error {
	resp, status := PostResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	article := &model.Article{
		PostBase: model.PostBase{
			Base: model.Base{ID: id},
		},
	}

	status, err := article.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !article.Release {
		status = http.StatusNotFound
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), article
	return c.JSON(status, resp)
}

// CreateArticle handles the "/blog/posts/articles/create" route.
func CreateArticle(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := PostResponse{}, 0
	article := &model.Article{}
	if err := c.Bind(article); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil) {
		status = http.StatusForbidden
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	article.Creator = claims.User
	if status, err := article.Create(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), article
	return c.JSON(status, resp)
}

// UpdateArticle handles the "/blog/posts/articles/:id/update" route.
func UpdateArticle(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := PostResponse{}, 0
	article, a := &model.Article{}, &model.Article{}
	if err := c.Bind(a); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	article.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(article.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := article.Read(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if article.Creator != claims.User {
		if (!article.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
			(article.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
			status = http.StatusForbidden
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}
	}

	article.PostBase, article.Content = a.PostBase, a.Content
	status, err := article.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), article
	return c.JSON(status, resp)
}

// TransferArticle handles the "/blog/posts/articles/:id/transfer" route.
func TransferArticle(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := PostResponse{}, 0
	article, a := &model.Article{}, &model.Article{}
	if err := c.Bind(a); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	article.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(article.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := article.Read(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if article.Creator != claims.User {
		if (!article.Release && !RBAC.IsGranted(string(claims.Role), permissionDraftOps, nil)) ||
			(article.Release && !RBAC.IsGranted(string(claims.Role), permissionPostOps, nil)) {
			status = http.StatusForbidden
			resp.Message = http.StatusText(status)
			return c.JSON(status, resp)
		}
	}

	if uuid.Equal(article.Creator, a.Creator) {
		status = http.StatusNotModified
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	article.Creator = a.Creator
	status, err := article.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Post = true, http.StatusText(status), article
	return c.JSON(status, resp)
}
