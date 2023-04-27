package controllers

import (
	"foosball/internal/invite"
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) GetUserInvites(c echo.Context) error {
	type getUserInvitesRequest struct {
		ID uint `param:"userId" validate:"required"`
	}

	type getUserInvitesResponse struct {
		Invites []invite.Invite `json:"invites"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[getUserInvitesRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	invites, err := h.inviteService.GetInvitesByUserID(ctx, req.ID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user invites")
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, invites)
}

func (h *Handlers) DeclineInvite(c echo.Context) error {
	type declineInviterequest struct {
		ID uint `param:"inviteId" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[declineInviterequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.inviteService.DeclineInvite(ctx, req.ID); err != nil {
		h.logger.WithError(err).Error("failed to decline invite")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) AcceptInvite(c echo.Context) error {
	type acceptInviteRequest struct {
		ID uint `param:"inviteId" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[acceptInviteRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.inviteService.AcceptInvite(ctx, req.ID); err != nil {
		h.logger.WithError(err).Error("failed to accept invite")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) InviteUserToOrganization(c echo.Context) error {
	type inviteUserToOrganizationRequest struct {
		OrganizationID uint   `json:"orgId" validate:"required"`
		Email          string `json:"email" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[inviteUserToOrganizationRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	exists, user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user by email")
		return echo.ErrInternalServerError
	}

	if !exists {
		return echo.ErrNotFound
	}

	if err := h.inviteService.CreateInvite(ctx, user.ID, req.OrganizationID); err != nil {
		h.logger.WithError(err).Error("failed to create invite")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
