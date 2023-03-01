package controllers

import (
	"encoding/json"
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getMatchRequest struct {
	Id uint `param:"id" validate:"required,numeric,gte=0"`
}

type getMatchResponse struct {
	TeamA  []uint `json:"teamA" validate:"required"`
	TeamB  []uint `json:"teamB" validate:"required"`
	GoalsA int    `json:"goalsA" validate:"required,numeric,gte=0,lte=10"`
	GoalsB int    `json:"goalsB" validate:"required,numeric,gte=0,lte=10"`
}

func (h *Handlers) GetMatch(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[getMatchRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	match, err := h.matchService.GetMatch(ctx, req.Id)
	if err != nil {
		// TODO Return 404 instead if not found
		h.logger.WithError(err).Error("failed to get match")
		return echo.ErrInternalServerError
	}

	var teamA []uint
	err = json.Unmarshal(match.TeamA, &teamA)
	if err != nil {
		h.logger.WithError(err).Error("failed to unmarshal teamA ids")
		return echo.ErrInternalServerError
	}

	var teamB []uint
	err = json.Unmarshal(match.TeamB, &teamB)
	if err != nil {
		h.logger.WithError(err).Error("failed to unmarshal teamB ids")
		return echo.ErrInternalServerError
	}

	resp := getMatchResponse{
		TeamA:  teamA,
		TeamB:  teamB,
		GoalsA: match.GoalsA,
		GoalsB: match.GoalsB,
	}
	return c.JSON(http.StatusOK, resp)
}

type postMatchRequest struct {
	TeamA  []uint `json:"teamA" validate:"required"`
	TeamB  []uint `json:"teamB" validate:"required"`
	GoalsA int    `json:"goalsA" validate:"required,numeric,min=0,max=10"`
	GoalsB int    `json:"goalsB" validate:"required,numeric,min=0,max=10"`
	Winner string `json:"winner" validate:"required,ascii,oneof=A B"`
}

type postMatchResponse struct {
	TeamA  []uint `json:"teamA" validate:"required"`
	TeamB  []uint `json:"teamB" validate:"required"`
	GoalsA int    `json:"goalsA" validate:"required,numeric,min=0,max=10"`
	GoalsB int    `json:"goalsB" validate:"required,numeric,min=0,max=10"`
	Winner string `json:"winner" validate:"required,ascii,oneof=A B"`
}

func (h *Handlers) PostMatch(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[postMatchRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	teamA, err := h.playerService.GetPlayers(ctx, req.TeamA)
	if err != nil {
		h.logger.WithError(err).Error("failed to get team A player ids")
		return echo.ErrInternalServerError
	}
	teamB, err := h.playerService.GetPlayers(ctx, req.TeamB)
	if err != nil {
		h.logger.WithError(err).Error("failed to get team B player ids")
		return echo.ErrInternalServerError
	}

	err = h.matchService.CreateMatch(ctx, *teamA, *teamB, req.GoalsA, req.GoalsB)
	if err != nil {
		h.logger.WithError(err).Error("failed to create match")
		return echo.ErrInternalServerError
	}

	resp := postMatchResponse{
		TeamA:  req.TeamA,
		TeamB:  req.TeamB,
		GoalsA: req.GoalsA,
		GoalsB: req.GoalsB,
		Winner: req.Winner,
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handlers) GetMatchesWithPlayer(c echo.Context) error {
	// TODO implement
	return c.JSON(http.StatusNotImplemented, nil)
}
