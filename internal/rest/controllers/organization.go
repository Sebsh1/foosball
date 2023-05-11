package controllers

import (
	"foosball/internal/rest/handlers"
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) CreateOrganization(c handlers.AuthenticatedContext) error {
	type createOrgRequest struct {
		Name         string `json:"name" validate:"required"`
		RatingMethod string `json:"ratingMethod" validate:"required,oneof=elo rms glicko2"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[createOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if c.Claims.OrganizationID != 0 {
		return echo.ErrForbidden
	}

	if err := h.organizationService.CreateOrganization(ctx, req.Name, req.RatingMethod); err != nil {
		h.logger.WithError(err).Error("failed to create organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) DeleteOrganization(c handlers.AuthenticatedContext) error {
	type deleteOrgRequest struct {
		ID uint `param:"orgId" validate:"required,gt=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[deleteOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if !(c.Claims.Admin && c.Claims.OrganizationID == req.ID) {
		return echo.ErrUnauthorized
	}

	if err := h.organizationService.DeleteOrganization(ctx, req.ID); err != nil {
		h.logger.WithError(err).Error("failed to delete organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateOrganization(c handlers.AuthenticatedContext) error {
	type updateOrgRequest struct {
		ID           uint   `param:"orgId" validate:"required,gt=0"`
		Name         string `json:"name" validate:"required"`
		RatingMethod string `json:"ratingMethod" validate:"required,oneof=elo rms glicko2"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[updateOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if !(c.Claims.Admin && c.Claims.OrganizationID == req.ID) {
		return echo.ErrUnauthorized
	}

	if err := h.organizationService.UpdateOrganization(ctx, req.ID, req.Name, req.RatingMethod); err != nil {
		h.logger.WithError(err).Error("failed to update organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
