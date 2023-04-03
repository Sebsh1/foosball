package controllers

import (
	"foosball/internal/rest/helpers"
	"foosball/internal/team"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getMatchRequest struct {
	Id uint `param:"id" validate:"required,numeric,gte=0"`
}

type getMatchResponse struct {
	SeasonID uint `json:"seasonID,omitempty"`
	TeamAID  uint `json:"teamAID" validate:"required"`
	TeamBID  uint `json:"teamBID" validate:"required"`
	GoalsA   int  `json:"goalsA" validate:"required,numeric,gte=0,lte=10"`
	GoalsB   int  `json:"goalsB" validate:"required,numeric,gte=0,lte=10"`
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

	resp := getMatchResponse{
		SeasonID: match.SeasonID,
		TeamAID:  match.TeamAID,
		TeamBID:  match.TeamBID,
		GoalsA:   match.GoalsA,
		GoalsB:   match.GoalsB,
	}

	return c.JSON(http.StatusOK, resp)
}

type postMatchRequest struct {
	SeasonID uint `json:"seasonId"`
	TeamAID  uint `json:"teamAId" validate:"required"`
	TeamBID  uint `json:"teamBId" validate:"required"`
	GoalsA   int  `json:"goalsA" validate:"required,numeric,min=0,max=10"`
	GoalsB   int  `json:"goalsB" validate:"required,numeric,min=0,max=10"`
}

type postMatchResponse struct {
	SeasonID uint `json:"seasonId"`
	TeamAID  uint `json:"teamAId" validate:"required"`
	TeamBID  uint `json:"teamBId" validate:"required"`
	GoalsA   int  `json:"goalsA" validate:"required,numeric,min=0,max=10"`
	GoalsB   int  `json:"goalsB" validate:"required,numeric,min=0,max=10"`
}

func (h *Handlers) PostMatch(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[postMatchRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.TeamAID == req.TeamBID || req.GoalsA == req.GoalsB {
		return echo.ErrBadRequest
	}

	teamA, err := h.teamService.GetTeam(ctx, req.TeamAID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get team A")
		return echo.ErrInternalServerError
	}
	teamB, err := h.teamService.GetTeam(ctx, req.TeamBID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get team B")
		return echo.ErrInternalServerError
	}

	err = h.matchService.CreateMatch(ctx, teamA.ID, teamB.ID, req.GoalsA, req.GoalsB)
	if err != nil {
		h.logger.WithError(err).Error("failed to create match")
		return echo.ErrInternalServerError
	}

	var winner, loser *team.Team
	if req.GoalsA > req.GoalsB {
		winner = teamA
		loser = teamB
	} else {
		winner = teamB
		loser = teamA
	}

	if err := h.ratingService.UpdateRatings(ctx, winner, loser); err != nil {
		h.logger.WithError(err).Error("failed to update ratings")
		return echo.ErrInternalServerError
	}

	resp := postMatchResponse{
		TeamAID: req.TeamAID,
		TeamBID: req.TeamBID,
		GoalsA:  req.GoalsA,
		GoalsB:  req.GoalsB,
	}
	return c.JSON(http.StatusCreated, resp)
}
