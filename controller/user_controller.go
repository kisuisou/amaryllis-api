package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}
