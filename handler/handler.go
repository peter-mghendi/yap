package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response is the base response type
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// AppController handles the "/" route.
func AppController(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Yap API!!\n")
}
