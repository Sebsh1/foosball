package handlers

import (
	"matchlog/internal/authentication"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthenticatedContext struct {
	echo.Context
	Log    *logrus.Entry
	Claims authentication.Claims
	JWT    string
}

type AuthenticatedHandlerFunc func(ctx AuthenticatedContext) error

func AuthenticatedHandlerFactory(logger *logrus.Entry) func(handler AuthenticatedHandlerFunc) func(ctx echo.Context) error {
	return func(handler AuthenticatedHandlerFunc) func(ctx echo.Context) error {
		return func(ctx echo.Context) error {
			claims, ok := ctx.Get("jwt_claims").(*authentication.Claims)
			if !ok {
				logger.Debug("missing jwt_claims")

				return echo.ErrUnauthorized
			}

			jwt, ok := ctx.Get("jwt").(string)
			if !ok {
				logger.Debug("missing jwt")

				return echo.ErrUnauthorized
			}

			logger = logger.WithFields(logrus.Fields{
				"user_id": claims.UserID,
				"org_id":  claims.OrganizationID,
				"name":    claims.Name,
				"role":    claims.Role,
			})

			return handler(AuthenticatedContext{
				Context: ctx,
				Claims:  *claims,
				Log:     logger,
				JWT:     jwt,
			})
		}
	}
}
