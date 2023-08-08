package controllers

import (
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) CreateOrganization(c handlers.AuthenticatedContext) error {
	type createOrgRequest struct {
		Name string `json:"name" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[createOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if c.Claims.OrganizationID != 0 {
		return echo.ErrForbidden
	}

	if err := h.organizationService.CreateOrganization(ctx, req.Name); err != nil {
		h.logger.WithError(err).Error("failed to create organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handlers) DeleteOrganization(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	if err := h.organizationService.DeleteOrganization(ctx, c.Claims.OrganizationID); err != nil {
		h.logger.WithError(err).Error("failed to delete organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateOrganization(c handlers.AuthenticatedContext) error {
	type updateOrgRequest struct {
		Name string `json:"name" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[updateOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.organizationService.UpdateOrganization(ctx, c.Claims.OrganizationID, req.Name); err != nil {
		h.logger.WithError(err).Error("failed to update organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
