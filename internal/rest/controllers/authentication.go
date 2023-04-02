package controllers

import (
	"foosball/internal/rest/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Username string `json:"username" validate:"required,ascii"`
	Password string `json:"password" validate:"required,ascii"`
}

type loginResponse struct {
	JWT string `json:"jwt" validate:"required,ascii"`
}

func (h *Handlers) Login(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[loginRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	correct, token, err := h.authService.Login(ctx, req.Username, req.Password)
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
	return c.JSON(http.StatusCreated, resp)
}
