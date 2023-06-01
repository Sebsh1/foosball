package controllers

import (
	"matchlog/internal/invite"
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) GetUserInvites(c handlers.AuthenticatedContext) error {
	type getUserInvitesRequest struct {
		ID uint `param:"userId" validate:"required"`
	}

	type responseInvite struct {
		OrganizationID uint `json:"organizationId"`
		UserID         uint `json:"userId"`
	}

	type getUserInvitesResponse struct {
		Invites []responseInvite `json:"invites"`
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

	respInvites := make([]responseInvite, len(invites))
	for i, invite := range invites {
		respInvites[i] = responseInvite{
			OrganizationID: invite.OrganizationID,
			UserID:         invite.UserID,
		}
	}

	resp := getUserInvitesResponse{
		Invites: respInvites,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) RespondToInvite(c handlers.AuthenticatedContext) error {
	type acceptInviteRequest struct {
		UserID   uint `param:"userId" validate:"required,gte=0"`
		InviteID uint `param:"inviteId" validate:"required,gte=0"`
		Accepted bool `json:"accepted" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[acceptInviteRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.UserID != c.Claims.UserID {
		return echo.ErrUnauthorized
	}

	if req.Accepted {
		if err := h.inviteService.AcceptInvite(ctx, req.InviteID); err != nil {
			if err == invite.ErrNotFound {
				return echo.ErrNotFound
			}

			h.logger.WithError(err).Error("failed to accept invite")
			return echo.ErrInternalServerError
		}
	} else {
		if err := h.inviteService.DeclineInvite(ctx, req.InviteID); err != nil {
			if err == invite.ErrNotFound {
				return echo.ErrNotFound
			}

			h.logger.WithError(err).Error("failed to decline invite")
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) InviteUsersToOrganization(c handlers.AuthenticatedContext) error {
	type inviteUsersToOrganizationRequest struct {
		OrganizationID uint     `param:"orgId" validate:"required,gt=0"`
		Emails         []string `json:"emails" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[inviteUsersToOrganizationRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.OrganizationID != c.Claims.OrganizationID {
		return echo.ErrUnauthorized
	}

	userIDs := make([]uint, len(req.Emails))
	for _, email := range req.Emails {
		exists, user, err := h.userService.GetUserByEmail(ctx, email)
		if err != nil {
			h.logger.WithError(err).Error("failed to get user by email")
			return echo.ErrInternalServerError
		}

		if !exists {
			return echo.ErrNotFound
		}

		userIDs = append(userIDs, user.ID)
	}

	if err := h.inviteService.CreateInvites(ctx, userIDs, req.OrganizationID); err != nil {
		h.logger.WithError(err).Error("failed to create invites")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
