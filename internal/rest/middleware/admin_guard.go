package middleware

import (
	"foosball/internal/authentication"
	"foosball/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func AdminGuard(logger *logrus.Entry) func(next echo.HandlerFunc) echo.HandlerFunc {
	l := logger.WithField("middleware", "admin_guard")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("jwt_claims").(*authentication.Claims)
			if !ok {
				l.Debug("missing jwt_claims")

				return echo.ErrForbidden
			}

			if claims.Role != string(user.AdminRole) {
				l.WithFields(logrus.Fields{
					"Name":           claims.Name,
					"UserID":         claims.UserID,
					"OrganizationID": claims.OrganizationID,
					"Role":           claims.Role,
				}).Debug("user attempted to perform admin action without admin role")

				return echo.ErrForbidden
			}

			return next(c)
		}
	}
}
