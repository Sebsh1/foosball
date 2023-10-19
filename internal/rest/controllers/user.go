package controllers

import (
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) DeleteUser(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	if err := h.userService.DeleteUser(ctx, c.Claims.UserId); err != nil {
		h.logger.Error("failed to delete user",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) RemoveUserFromClub(c handlers.AuthenticatedContext) error {
	type request struct {
		UserId uint `param:"userId" validate:"required,gt=0"`
		ClubId uint `param:"ClubId" validate:"required,gt=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.ClubService.RemoveUserFromClub(ctx, req.UserId, req.ClubId); err != nil {
		h.logger.Error("failed to remove user from Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) AddVirtualUserToClub(c handlers.AuthenticatedContext) error {
	type request struct {
		Name string `json:"name" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}
	if err := h.userService.CreateVirtualUser(ctx, req.Name); err != nil {
		h.logger.Error("failed to create virtual user",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) RespondToInvite(c handlers.AuthenticatedContext) error {
	type request struct {
		InviteId uint `param:"inviteId" validate:"required,gt=0"`
		Accept   bool `json:"accept" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.RespondToInvite(ctx, req.InviteId, req.Accept); err != nil {
		h.logger.Error("failed to respond to invite",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
