package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AppController handles the "/" route.
func AppController(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Yap API!!\n")
}
