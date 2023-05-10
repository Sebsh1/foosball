package controllers

import (
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) Login(c echo.Context) error {
	type loginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	type loginResponse struct {
		JWT string `json:"jwt"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[loginRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	correct, token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.WithError(err).Error("failed to login")
		return echo.ErrInternalServerError
	}

	if !correct {
		return echo.ErrUnauthorized
	}

	resp := loginResponse{
		JWT: token,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) Signup(c echo.Context) error {
	type signupRequest struct {
		Email    string `json:"email" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[signupRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	success, err := h.authService.Signup(ctx, req.Email, req.Name, req.Password)
	if err != nil {
		h.logger.WithError(err).Error("failed to signup")
		return echo.ErrInternalServerError
	}

	if !success {
		return echo.ErrBadRequest
	}

	return c.NoContent(http.StatusCreated)
}
