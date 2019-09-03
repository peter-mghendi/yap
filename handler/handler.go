package handler

import (
	"net/http"

	"github.com/l3njo/yap-api/model"
	"github.com/labstack/echo/v4"
)

// Response is the base response type
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

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

// AppController handles the "/" route.
func AppController(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Yap API!!\n")
}
