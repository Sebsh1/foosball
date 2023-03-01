package controllers

import (
	"foosball/internal/models"
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getLeaderboardRatingRequest struct {
	Top int `query:"top" validate:"numeric,max=100" default:"10"`
}

type getLeaderboardRatingResponse struct {
	Leaderboard []*models.Player `json:"leaderboard" validate:"required"`
}

func (h *Handlers) GetLeaderboardRating(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[getLeaderboardRatingRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	players, err := h.playerService.GetTopPlayersByRating(ctx, req.Top)

	resp := getLeaderboardRatingResponse{
		Leaderboard: players,
	}

	return c.JSON(http.StatusNotImplemented, resp)
}

func (h *Handlers) GetLeaderboardMatchesPlayed(c echo.Context) error {
	// TODO implement
	return c.JSON(http.StatusNotImplemented, nil)
}
