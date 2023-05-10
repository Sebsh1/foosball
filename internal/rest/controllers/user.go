package controllers

import (
	"foosball/internal/rest/handlers"
	"foosball/internal/rest/helpers"
	"foosball/internal/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) DeleteUser(c handlers.AuthenticatedContext) error {
	type deleteUserRequest struct {
		ID uint `param:"userId" validate:"required, gte=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[deleteUserRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.ID != c.Claims.UserID {
		return echo.ErrUnauthorized
	}

	if err := h.userService.DeleteUser(ctx, req.ID); err != nil {
		h.logger.WithError(err).Error("failed to delete user")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) GetUsersInOrganization(c handlers.AuthenticatedContext) error {
	type getUsersInOrgRequest struct {
		OrgId uint `json:"orgId" validate:"required, gte=1"`
	}

	type getUsersInOrgResponse struct {
		Users []user.User `json:"users"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[getUsersInOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if req.OrgId != c.Claims.OrganizationID {
		return echo.ErrUnauthorized
	}

	users, err := h.userService.GetUsersInOrganization(ctx, req.OrgId)
	if err != nil {
		h.logger.WithError(err).Error("failed to get users in organization")
		return echo.ErrInternalServerError
	}

	reponse := getUsersInOrgResponse{
		Users: users,
	}

	return c.JSON(http.StatusOK, reponse)
}

func (h *Handlers) RemoveUserFromOrganization(c handlers.AuthenticatedContext) error {
	type removeUserFromOrgRequest struct {
		OrgId  uint `param:"orgId" validate:"required, gte=1"`
		UserId uint `param:"userId" validate:"required, gte=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[removeUserFromOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if !(c.Claims.Admin && req.OrgId == c.Claims.OrganizationID) || req.UserId == c.Claims.UserID {
		return echo.ErrUnauthorized
	}

	users, err := h.userService.GetUsers(ctx, []uint{req.UserId})
	if err != nil {
		h.logger.WithError(err).Error("failed to get users")
		return echo.ErrInternalServerError
	}

	u := users[0]
	if *u.OrganizationID != req.OrgId {
		return echo.ErrBadRequest
	}

	if err := h.userService.UpdateUser(ctx, u.ID, u.Email, u.Name, u.Hash, 0, false); err != nil {
		h.logger.WithError(err).Error("failed to remove user from organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) SetUserAsAdmin(c handlers.AuthenticatedContext) error {
	type setUserAsAdminRequest struct {
		OrgId  uint `param:"orgId" validate:"required, gte=1"`
		UserId uint `param:"userId" validate:"required, gte=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[setUserAsAdminRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if !(c.Claims.Admin && req.OrgId == c.Claims.OrganizationID) {
		return echo.ErrUnauthorized
	}

	users, err := h.userService.GetUsers(ctx, []uint{req.UserId})
	if err != nil {
		h.logger.WithError(err).Error("failed to get users")
		return echo.ErrInternalServerError
	}

	u := users[0]
	if *u.OrganizationID != req.OrgId {
		return echo.ErrBadRequest
	}

	if err := h.userService.UpdateUser(ctx, u.ID, u.Email, u.Name, u.Hash, *u.OrganizationID, true); err != nil {
		h.logger.WithError(err).Error("failed to set user as admin")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) RemoveAdminFromUser(c handlers.AuthenticatedContext) error {
	type removeAdminFromUserRequest struct {
		OrgId  uint `param:"orgId" validate:"required, gte=1"`
		UserId uint `param:"userId" validate:"required, gte=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[removeAdminFromUserRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if !(c.Claims.Admin && req.OrgId == c.Claims.OrganizationID) {
		return echo.ErrUnauthorized
	}

	users, err := h.userService.GetUsers(ctx, []uint{req.UserId})
	if err != nil {
		h.logger.WithError(err).Error("failed to get users")
		return echo.ErrInternalServerError
	}

	u := users[0]
	if *u.OrganizationID != req.OrgId {
		return echo.ErrBadRequest
	}

	if err := h.userService.UpdateUser(ctx, u.ID, u.Email, u.Name, u.Hash, *u.OrganizationID, false); err != nil {
		h.logger.WithError(err).Error("failed to remove admin from user")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
