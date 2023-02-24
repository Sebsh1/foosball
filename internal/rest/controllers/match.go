package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) PostMatch(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, nil)
}

func (h *Handlers) GetMatch(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, nil)
}

func (h *Handlers) GetMatchesWithPlayer(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, nil)
}
