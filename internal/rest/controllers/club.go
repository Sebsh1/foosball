package controllers

import (
	"matchlog/internal/club"
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) CreateClub(c handlers.AuthenticatedContext) error {
	type request struct {
		Name string `json:"name" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	_, err = h.clubService.CreateClub(ctx, req.Name, c.Claims.UserID)
	if err != nil {
		h.logger.Error("failed to create Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handlers) DeleteClub(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubID uint `json:"ClubId" validate:"required,gt=0"`
	}

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	ctx := c.Request().Context()

	if err := h.clubService.DeleteClub(ctx, req.ClubID); err != nil {
		h.logger.Error("failed to delete Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateClub(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubID uint   `json:"ClubId" validate:"required,gt=0"`
		Name   string `json:"name" validate:"required"`
	}

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	ctx := c.Request().Context()

	if err := h.clubService.UpdateClub(ctx, req.ClubID, req.Name); err != nil {
		h.logger.Error("failed to update Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateUserRole(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubID uint      `param:"ClubId" validate:"required,gt=0"`
		UserID uint      `param:"userId" validate:"required,gt=0"`
		Role   club.Role `json:"role" validate:"required"`
	}

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	ctx := c.Request().Context()

	if err := h.clubService.UpdateUserRole(ctx, req.UserID, req.ClubID, req.Role); err != nil {
		h.logger.Error("failed to update user role",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) GetUsersInClub(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubID uint `query:"ClubId" validate:"required,gt=0"`
	}

	type responseUser struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	type response struct {
		Users []responseUser `json:"users"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	userIDs, err := h.clubService.GetUserIDsInClub(ctx, req.ClubID)
	if err != nil {
		h.logger.Error("failed to get userIDs in Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	users, err := h.userService.GetUsers(ctx, userIDs)
	if err != nil {
		h.logger.Error("failed to get users",
			"error", err)
		return echo.ErrInternalServerError
	}

	respUsers := make([]responseUser, len(users))
	for i, u := range users {
		respUsers[i] = responseUser{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	resp := response{
		Users: respUsers,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetUserInvites(c handlers.AuthenticatedContext) error {
	type responseClub struct {
		ID     uint   `json:"id"`
		ClubID uint   `json:"org_id"`
		Name   string `json:"name"`
	}

	type response struct {
		Invites []responseClub `json:"invites"`
	}

	ctx := c.Request().Context()

	orgUsers, err := h.clubService.GetInvitesByUserID(ctx, c.Claims.UserID)
	if err != nil {
		h.logger.Error("failed to get user invites",
			"error", err)
		return echo.ErrInternalServerError
	}

	orgIDs := make([]uint, len(orgUsers))
	for i, orgUser := range orgUsers {
		orgIDs[i] = orgUser.ClubId
	}

	orgs, err := h.clubService.GetClubs(ctx, orgIDs)
	if err != nil {
		h.logger.Error("failed to get Clubs",
			"error", err)
		return echo.ErrInternalServerError
	}

	invites := make([]responseClub, len(orgs))
	for i, org := range orgs {
		invites[i] = responseClub{
			ID:     orgUsers[i].Id,
			ClubID: org.Id,
			Name:   org.Name,
		}
	}

	resp := response{
		Invites: invites,
	}

	return c.JSON(http.StatusOK, resp)
}
