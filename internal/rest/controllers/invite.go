package controllers

import (
	"matchlog/internal/invite"
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) GetUserInvites(c handlers.AuthenticatedContext) error {
	type responseInvite struct {
		InviteID       uint `json:"invite_id"`
		OrganizationID uint `json:"organization_id"`
	}

	type getUserInvitesResponse struct {
		Invites []responseInvite `json:"invites"`
	}

	ctx := c.Request().Context()

	invites, err := h.inviteService.GetInvitesByUserID(ctx, c.Claims.UserID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user invites")
		return echo.ErrInternalServerError
	}

	respInvites := make([]responseInvite, len(invites))
	for i, invite := range invites {
		respInvites[i] = responseInvite{
			InviteID:       invite.ID,
			OrganizationID: invite.OrganizationID,
		}
	}

	resp := getUserInvitesResponse{
		Invites: respInvites,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) RespondToInvite(c handlers.AuthenticatedContext) error {
	type acceptInviteRequest struct {
		InviteID uint `param:"inviteId" validate:"required,gte=0"`
		Accept   bool `json:"accept" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[acceptInviteRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.Accept {
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
		Emails []string `json:"emails" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[inviteUsersToOrganizationRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	userIDs := make([]uint, len(req.Emails))
	for _, email := range req.Emails {
		exists, u, err := h.userService.GetUserByEmail(ctx, email)
		if err != nil {
			h.logger.WithError(err).Error("failed to get user by email")
			return echo.ErrInternalServerError
		}

		if !exists {
			continue
		}

		if u.OrganizationID != nil && *u.OrganizationID == c.Claims.OrganizationID {
			continue
		}

		userIDs = append(userIDs, u.ID)
	}

	if err := h.inviteService.CreateInvites(ctx, userIDs, c.Claims.OrganizationID); err != nil {
		h.logger.WithError(err).Error("failed to create invites")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
