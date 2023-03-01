package controllers

import (
	"errors"
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

var ErrNotFound = errors.New("not found")

type getPlayerRequest struct {
	Id uint `param:"id" validate:"required,numeric,gte=0"`
}

type getPlayerResponse struct {
	Id     uint   `json:"id" validate:"required,numeric,gte=0"`
	Name   string `json:"name" validate:"required,ascii"`
	Rating int    `json:"rating" validate:"required,numeric,gte=0"`
}

func (h *Handlers) GetPlayer(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[getPlayerRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	player, err := h.playerService.GetPlayer(ctx, req.Id)
	if err != nil {
		// TODO Return 404 instead if not found
		h.logger.WithError(err).Error("failed to create player")
		return echo.ErrInternalServerError
	}

	resp := getPlayerResponse{
		Id:     player.ID,
		Name:   player.Name,
		Rating: player.Rating,
	}

	return c.JSON(http.StatusOK, resp)
}

type postPlayerRequest struct {
	Name string `param:"name" validate:"required,ascii"`
}

type postPlayerResponse struct {
	Name   string `json:"name" validate:"required,ascii"`
	Rating int    `json:"rating" validate:"required,numeric,gte=0"`
}

func (h *Handlers) PostPlayer(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[postPlayerRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	err = h.playerService.CreatePlayer(ctx, req.Name)
	if err != nil {
		h.logger.WithError(err).Error("failed to create player")
		return echo.ErrInternalServerError
	}

	resp := postPlayerResponse{
		Name:   req.Name,
		Rating: 1000,
	}
	return c.JSON(http.StatusCreated, resp)
}

type deletePlayerRequest struct {
	ID uint `param:"id" validate:"required,numeric"`
}

func (h *Handlers) DeletePlayer(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[deletePlayerRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	err = h.playerService.DeletePlayer(ctx, req.ID)
	if err != nil {
		h.logger.WithError(err).Error("failed to delete player")
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *Handlers) GetPlayerStatistics(c echo.Context) error {
	// TODO: name, rating, total matches played, w/l, w/l by players,
	return c.JSON(http.StatusNotImplemented, nil)
}
