package controllers

import (
	"foosball/internal/invite"
	"foosball/internal/rest/handlers"
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) GetUserInvites(c handlers.AuthenticatedContext) error {
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

	if req.ID != c.Claims.UserID {
		return echo.ErrUnauthorized
	}

	invites, err := h.inviteService.GetInvitesByUserID(ctx, req.ID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user invites")
		return echo.ErrInternalServerError
	}

	resp := getUserInvitesResponse{
		Invites: invites,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) DeclineInvite(c handlers.AuthenticatedContext) error {
	type declineInviterequest struct {
		UserID   uint `param:"userId" validate:"required,gte=0"`
		InviteID uint `param:"inviteId" validate:"required,gte=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[declineInviterequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.UserID != c.Claims.UserID {
		return echo.ErrUnauthorized
	}

	if err := h.inviteService.DeclineInvite(ctx, req.InviteID); err != nil {
		if err == invite.ErrNotFound {
			return echo.ErrNotFound
		}

		h.logger.WithError(err).Error("failed to decline invite")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) AcceptInvite(c handlers.AuthenticatedContext) error {
	type acceptInviteRequest struct {
		UserID   uint `param:"userId" validate:"required,gte=0"`
		InviteID uint `param:"inviteId" validate:"required,gte=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[acceptInviteRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.UserID != c.Claims.UserID {
		return echo.ErrUnauthorized
	}

	if err := h.inviteService.AcceptInvite(ctx, req.InviteID); err != nil {
		if err == invite.ErrNotFound {
			return echo.ErrNotFound
		}

		h.logger.WithError(err).Error("failed to accept invite")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) InviteUserToOrganization(c handlers.AuthenticatedContext) error {
	type inviteUserToOrganizationRequest struct {
		OrganizationID uint   `json:"orgId" validate:"required,gt=0"`
		Email          string `json:"email" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[inviteUserToOrganizationRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.OrganizationID != c.Claims.OrganizationID {
		return echo.ErrUnauthorized
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
