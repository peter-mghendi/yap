package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/l3njo/yap/model"
	"github.com/l3njo/yap/util"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo/v4"
)

// ReactionResponse is a response containing one Reaction
type ReactionResponse struct {
	Response
	model.Reaction `json:"data"`
}

// ReactionsResponse is a response containing a slice of Reactions
type ReactionsResponse struct {
	Response
	Reactions []model.Reaction `json:"data"`
}

// GetBlogPostReactions handles the "/blog/posts/:id/reactions" route.
func GetBlogPostReactions(c echo.Context) error {
	resp, status := ReactionsResponse{}, 0
	reactions, status, err := model.ReadAllReactions()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	post := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(uuid.Nil, post) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reactions = util.FilterR(reactions, func(r model.Reaction) bool {
		return (r.Site == "blog") && (r.Item == post)
	})

	resp.Status, resp.Message, resp.Reactions = true, http.StatusText(status), reactions
	return c.JSON(status, resp)
}

// GetUserBlogReactions handles the "/users/:id/blog/reactions" route.
func GetUserBlogReactions(c echo.Context) error {
	resp, status := ReactionsResponse{}, 0
	reactions, status, err := model.ReadAllReactions()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	user := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(uuid.Nil, user) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reactions = util.FilterR(reactions, func(r model.Reaction) bool {
		return (r.Site == "blog") && (r.User == user)
	})

	resp.Status, resp.Message, resp.Reactions = true, http.StatusText(status), reactions
	return c.JSON(status, resp)
}

// GetBlogPostReactionByID handles the "/blog/reactions/:id" route.
func GetBlogPostReactionByID(c echo.Context) error {
	resp, status := ReactionResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("reaction"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	post := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction := model.Reaction{
		Base: model.Base{ID: id},
		Site: "blog",
		Item: post,
	}

	status, err := reaction.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Reaction = true, http.StatusText(status), reaction
	return c.JSON(status, resp)
}

// CreateBlogReaction handles the "/blog/reactions/create" route.
func CreateBlogReaction(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := ReactionResponse{}, 0
	reaction := model.Reaction{}
	if err := c.Bind(&reaction); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction.Site = "blog"
	post := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(uuid.Nil, post) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction.User, reaction.Item = claims.User, post
	if status, err := reaction.Create(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Reaction = true, http.StatusText(status), reaction
	return c.JSON(status, resp)
}

// UpdateBlogReaction handles the "/blog/reactions/:id/update" route.
func UpdateBlogReaction(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)

	resp, status := ReactionResponse{}, 0
	reaction, r := model.Reaction{Site: "blog"}, model.Reaction{}
	if err := c.Bind(&r); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction.ID = uuid.FromStringOrNil(c.Param("reaction"))
	if uuid.Equal(reaction.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction.Item = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(reaction.Item, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := reaction.Read(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if reaction.User != claims.User {
		status = http.StatusForbidden
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction.Text = r.Text
	status, err := reaction.Update()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Reaction = true, http.StatusText(status), reaction
	return c.JSON(status, resp)
}

// DeleteBlogReaction handles the "/blog/reactions/:id/delete" route.
func DeleteBlogReaction(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*JwtCustomClaims)
	reaction, resp, status := model.Reaction{Site: "blog"}, ReactionResponse{}, 0
	reaction.ID = uuid.FromStringOrNil(c.Param("reaction"))
	if uuid.Equal(reaction.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction.Item = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(reaction.Item, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err := reaction.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if !RBAC.IsGranted(string(claims.Role), permissionReactionOps, nil) && !uuid.Equal(claims.User, reaction.User) {
		status = http.StatusForbidden
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	status, err = reaction.Delete()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message = true, http.StatusText(status)
	return c.JSON(status, resp)
}
