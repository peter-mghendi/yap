package handler

import (
	"net/http"

	"github.com/l3njo/yap-api/model"
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

// GetReactions handles the "/reactions" route.
func GetReactions(c echo.Context) error {
	resp, status := ReactionsResponse{}, 0
	reactions, status, err := model.ReadAllReactions()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Reactions = true, http.StatusText(status), reactions
	return c.JSON(status, resp)
}

// GetReactionByID handles the "/reactions/:id" route.
func GetReactionByID(c echo.Context) error {
	resp, status := ReactionResponse{}, 0
	id := uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(id, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction := model.Reaction{
		Base: model.Base{ID: id},
	}

	status, err := reaction.Read()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Reaction = true, http.StatusText(status), reaction
	return c.JSON(status, resp)
}

// CreateReaction handles the "/reactions/create" route.
func CreateReaction(c echo.Context) error {
	resp, status := ReactionResponse{}, 0
	reaction := model.Reaction{}
	if err := c.Bind(&reaction); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := reaction.Create(); err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message, resp.Reaction = true, http.StatusText(status), reaction
	return c.JSON(status, resp)
}

// UpdateReaction handles the "/reactions/:id/update" route.
func UpdateReaction(c echo.Context) error {
	resp, status := ReactionResponse{}, 0
	reaction, r := model.Reaction{}, model.Reaction{}
	if err := c.Bind(&r); err != nil {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	reaction.ID = uuid.FromStringOrNil(c.Param("id"))
	if uuid.Equal(reaction.ID, uuid.Nil) {
		status = http.StatusBadRequest
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	if status, err := reaction.Read(); err != nil {
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

// DeleteReaction handles the "/reactions/:id/delete" route.
func DeleteReaction(c echo.Context) error {
	resp, status := ReactionResponse{}, 0
	reaction := model.Reaction{
		Base: model.Base{ID: uuid.FromStringOrNil(c.Param("id"))},
	}

	status, err := reaction.Delete()
	if err != nil {
		resp.Message = http.StatusText(status)
		return c.JSON(status, resp)
	}

	resp.Status, resp.Message = true, http.StatusText(status)
	return c.JSON(status, resp)
}
