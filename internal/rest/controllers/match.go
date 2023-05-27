package controllers

import (
	"matchlog/internal/match"
	"matchlog/internal/organization"
	"matchlog/internal/rating"
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"matchlog/internal/statistic"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) PostMatch(c handlers.AuthenticatedContext) error {
	type postMatchRequest struct {
		OrganiziationID uint        `json:"organizationId" validate:"required"`
		TeamA           []uint      `json:"teamA" validate:"required"`
		TeamB           []uint      `json:"teamB" validate:"required"`
		Sets            []match.Set `json:"sets" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[postMatchRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err = h.matchService.CreateMatch(ctx, req.TeamA, req.TeamB, req.Sets); err != nil {
		h.logger.WithError(err).Error("failed to create match")
		return echo.ErrInternalServerError
	}

	draw, winners, losers := h.matchService.DetermineResult(ctx, req.TeamA, req.TeamB, req.Sets)

	org, err := h.organizationService.GetOrganization(ctx, req.OrganiziationID)
	if err != nil {
		if err == organization.ErrNotFound {
			return echo.ErrNotFound
		}

		h.logger.WithError(err).Error("failed to get organization")
		return echo.ErrInternalServerError
	}

	if draw {
		allPlayers := append(req.TeamA, req.TeamB...)
		if err := h.statisticService.UpdateStatisticsByUserIDs(ctx, allPlayers, statistic.ResultDraw); err != nil {
			h.logger.WithError(err).Error("failed to update statistics for draw")
			return echo.ErrInternalServerError
		}
	} else {
		if err := h.statisticService.UpdateStatisticsByUserIDs(ctx, winners, statistic.ResultWin); err != nil {
			h.logger.WithError(err).Error("failed to update statistics for winners")
			return echo.ErrInternalServerError
		}

		if err := h.statisticService.UpdateStatisticsByUserIDs(ctx, losers, statistic.ResultLoss); err != nil {
			h.logger.WithError(err).Error("failed to update statistics for losers")
			return echo.ErrInternalServerError
		}
	}

	if err := h.ratingService.UpdateRatings(ctx, rating.Method(org.RatingMethod), draw, winners, losers); err != nil {
		h.logger.WithError(err).Error("failed to update ratings")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}
