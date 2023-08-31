package middleware

import (
	"matchlog/internal/authentication"
	"matchlog/internal/user"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func AdminGuard(logger *zap.SugaredLogger) func(next echo.HandlerFunc) echo.HandlerFunc {
	l := logger.With("middleware", "admin_guard")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("jwt_claims").(*authentication.Claims)
			if !ok {
				l.Debug("missing jwt_claims")

				return echo.ErrForbidden
			}

			if claims.Role != string(user.AdminRole) {
				l.Debug("user attempted to perform admin action without admin role",
					"name", claims.Name,
					"user_id", claims.UserID,
					"org_id", claims.OrganizationID)

				return echo.ErrForbidden
			}

			return next(c)
		}
	}
}
