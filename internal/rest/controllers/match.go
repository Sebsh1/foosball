package controllers

import (
	"matchlog/internal/match"
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"matchlog/internal/statistic"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) PostMatch(c handlers.AuthenticatedContext) error {
	type request struct {
		TeamA   []uint `json:"teamA" validate:"required"`
		TeamB   []uint `json:"teamB" validate:"required"`
		ScoresA []int  `json:"scoresA" validate:"required"`
		ScoresB []int  `json:"scoresB" validate:"required"`
		Rated   bool   `json:"rated" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if len(req.ScoresA) != len(req.ScoresB) {
		return echo.ErrBadRequest
	}

	result, winners, losers := h.matchService.DetermineResult(ctx, req.TeamA, req.TeamB, req.ScoresA, req.ScoresB)

	if err = h.matchService.CreateMatch(ctx, req.TeamA, req.TeamB, req.ScoresA, req.ScoresB, result); err != nil {
		h.logger.Error("failed to create match",
			"error", err)
		return echo.ErrInternalServerError
	}

	if result == match.Draw {
		allPlayers := append(req.TeamA, req.TeamB...)
		if err := h.statisticService.UpdateStatisticsByUserIDs(ctx, allPlayers, statistic.ResultDraw); err != nil {
			h.logger.Error("failed to update statistics for draw",
				"error", err)
			return echo.ErrInternalServerError
		}
	} else {
		if err := h.statisticService.UpdateStatisticsByUserIDs(ctx, winners, statistic.ResultWin); err != nil {
			h.logger.Error("failed to update statistics for winners",
				"error", err)
			return echo.ErrInternalServerError
		}

		if err := h.statisticService.UpdateStatisticsByUserIDs(ctx, losers, statistic.ResultLoss); err != nil {
			h.logger.Error("failed to update statistics for losers",
				"error", err)
			return echo.ErrInternalServerError
		}
	}

	if !req.Rated {
		return c.NoContent(http.StatusCreated)
	}

	isDraw := result == match.Draw
	if err := h.ratingService.UpdateRatings(ctx, isDraw, winners, losers); err != nil {
		h.logger.Error("failed to update ratings",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}
