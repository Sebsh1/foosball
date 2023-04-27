package controllers

import (
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) CreateOrganization(c echo.Context) error {
	type createOrgRequest struct {
		Name         string `json:"name" validate:"required"`
		RatingMethod string `json:"ratingMethod" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[createOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.organizationService.CreateOrganization(ctx, req.Name, req.RatingMethod); err != nil {
		h.logger.WithError(err).Error("failed to create organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) DeleteOrganization(c echo.Context) error {
	type deleteOrgRequest struct {
		ID uint `param:"orgId" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[deleteOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.organizationService.DeleteOrganization(ctx, req.ID); err != nil {
		h.logger.WithError(err).Error("failed to delete organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateOrganization(c echo.Context) error {
	type updateOrgRequest struct {
		ID           uint   `param:"orgId" validate:"required"`
		Name         string `json:"name" validate:"required"`
		RatingMethod string `json:"ratingMethod" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[updateOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.organizationService.UpdateOrganization(ctx, req.ID, req.Name, req.RatingMethod); err != nil {
		h.logger.WithError(err).Error("failed to update organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
