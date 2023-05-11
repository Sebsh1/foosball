package controllers

import (
	"foosball/internal/match"
	"foosball/internal/organization"
	"foosball/internal/rating"
	"foosball/internal/rest/handlers"
	"foosball/internal/rest/helpers"
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

	org, err := h.organizationService.GetOrganization(ctx, req.OrganiziationID)
	if err != nil {
		if err == organization.ErrNotFound {
			return echo.ErrNotFound
		}

		h.logger.WithError(err).Error("failed to get organization")
		return echo.ErrInternalServerError
	}

	if err = h.matchService.CreateMatch(ctx, req.TeamA, req.TeamB, req.Sets); err != nil {
		h.logger.WithError(err).Error("failed to create match")
		return echo.ErrInternalServerError
	}

	draw, winner, loser := h.matchService.DetermineResult(ctx, req.TeamA, req.TeamB, req.Sets)

	if err := h.ratingService.UpdateRatings(ctx, rating.Method(org.RatingMethod), draw, winner, loser); err != nil {
		h.logger.WithError(err).Error("failed to update ratings")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}
