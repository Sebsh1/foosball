package controllers

import (
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/helpers"
	"matchlog/internal/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) DeleteUser(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	if err := h.userService.DeleteUser(ctx, c.Claims.UserID); err != nil {
		h.logger.WithError(err).Error("failed to delete user")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) GetUsersInOrganization(c handlers.AuthenticatedContext) error {
	type responseUser struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	type getUsersInOrgResponse struct {
		Users []responseUser `json:"users"`
	}

	ctx := c.Request().Context()

	users, err := h.userService.GetUsersInOrganization(ctx, c.Claims.OrganizationID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get users in organization")
		return echo.ErrInternalServerError
	}

	respUsers := make([]responseUser, len(users))
	for i, u := range users {
		respUsers[i] = responseUser{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  string(u.Role),
		}
	}

	reponse := getUsersInOrgResponse{
		Users: respUsers,
	}

	return c.JSON(http.StatusOK, reponse)
}

func (h *Handlers) RemoveUserFromOrganization(c handlers.AuthenticatedContext) error {
	type removeUserFromOrgRequest struct {
		UserId uint `param:"userId" validate:"required,gt=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[removeUserFromOrgRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	u, err := h.userService.GetUser(ctx, req.UserId)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user")
		return echo.ErrInternalServerError
	}

	if *u.OrganizationID != c.Claims.OrganizationID {
		return echo.ErrBadRequest
	}

	if u.Role == user.VirtualRole {
		if err := h.userService.DeleteUser(ctx, u.ID); err != nil {
			h.logger.WithError(err).Error("failed to delete virtual user")
			return echo.ErrInternalServerError
		}

		return c.NoContent(http.StatusOK)
	}

	if err := h.userService.UpdateUser(ctx, u.ID, u.Email, u.Name, u.Hash, nil, user.NoneRole); err != nil {
		h.logger.WithError(err).Error("failed to remove user from organization")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateUserRole(c handlers.AuthenticatedContext) error {
	type setUserAsAdminRequest struct {
		UserId uint   `param:"userId" validate:"required,gt=0"`
		Role   string `json:"role" validate:"required,oneof=admin manager member"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[setUserAsAdminRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	u, err := h.userService.GetUser(ctx, req.UserId)
	if err != nil {
		h.logger.WithError(err).Error("failed to get users")
		return echo.ErrInternalServerError
	}

	if *u.OrganizationID != c.Claims.OrganizationID {
		return echo.ErrBadRequest
	}

	if err := h.userService.UpdateUser(ctx, u.ID, u.Email, u.Name, u.Hash, u.OrganizationID, user.Role(req.Role)); err != nil {
		h.logger.WithError(err).Error("failed to set user as admin")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) AddVirtualUserToOrganization(c handlers.AuthenticatedContext) error {
	type addVirtualUserToOrganizationRequest struct {
		Name string `json:"name" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[addVirtualUserToOrganizationRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}
	if err := h.userService.CreateVirtualUser(ctx, req.Name, c.Claims.OrganizationID); err != nil {
		h.logger.WithError(err).Error("failed to create virtual user")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) TransferVirtualUserToUser(c handlers.AuthenticatedContext) error {
	type transferVirtualUserToUserRequest struct {
		UserId        uint `param:"userId" validate:"required,gt=0"`
		VirtualUserId uint `param:"virtualUserId" validate:"required,gt=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[transferVirtualUserToUserRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	virtualUser, err := h.userService.GetUser(ctx, req.VirtualUserId)
	if err != nil {
		h.logger.WithError(err).Error("failed to get virtual user")
		return echo.ErrInternalServerError
	}

	if *virtualUser.OrganizationID != c.Claims.OrganizationID || virtualUser.Role != user.VirtualRole {
		return echo.ErrBadRequest
	}

	realUser, err := h.userService.GetUser(ctx, req.UserId)
	if err != nil {
		h.logger.WithError(err).Error("failed to get user")
		return echo.ErrInternalServerError
	}

	if *realUser.OrganizationID != c.Claims.OrganizationID || realUser.Role == user.VirtualRole {
		return echo.ErrBadRequest
	}

	if err := h.ratingService.TransferRatings(ctx, req.VirtualUserId, req.UserId); err != nil {
		h.logger.WithError(err).Error("failed to transfer ratings")
		return echo.ErrInternalServerError
	}

	if err := h.statisticService.TransferStatistics(ctx, req.VirtualUserId, req.UserId); err != nil {
		h.logger.WithError(err).Error("failed to transfer statistics")
		return echo.ErrInternalServerError
	}

	if err := h.userService.DeleteUser(ctx, req.VirtualUserId); err != nil {
		h.logger.WithError(err).Error("failed to delete virtual user")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
