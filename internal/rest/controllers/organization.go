package controllers

import (
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"matchlog/internal/user"
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

	orgID, err := h.organizationService.CreateOrganization(ctx, req.Name)
	if err != nil {
		h.logger.Error("failed to create organization",
			"error", err)
		return echo.ErrInternalServerError
	}

	userInfo, err := h.userService.GetUser(ctx, c.Claims.UserID)
	if err != nil {
		h.logger.Error("failed to get user",
			"error", err)
		return echo.ErrInternalServerError
	}

	if err := h.userService.UpdateUser(ctx, userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Hash, &orgID, user.AdminRole); err != nil {
		h.logger.Error("failed to add user to organization",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handlers) DeleteOrganization(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	if err := h.organizationService.DeleteOrganization(ctx, c.Claims.OrganizationID); err != nil {
		h.logger.Error("failed to delete organization",
			"error", err)
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
		h.logger.Error("failed to update organization",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
