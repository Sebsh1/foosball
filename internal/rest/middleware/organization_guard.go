package middleware

import (
	"encoding/json"
	"matchlog/internal/authentication"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func OrganizationGuard(logger *logrus.Entry) func(next echo.HandlerFunc) echo.HandlerFunc {
	l := logger.WithField("middleware", "organization_guard")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("jwt_claims").(*authentication.Claims)
			if !ok {
				l.Debug("missing jwt_claims")

				return echo.ErrForbidden
			}

			requestMap := make(map[string]interface{})
			if err := json.NewDecoder(c.Request().Body).Decode(&requestMap); err != nil {
				l.Debug("failed to decode request body")

				return echo.ErrBadRequest
			}

			orgID, ok := requestMap["organization_id"].(uint)
			if !ok {
				l.Debug("missing organization_id")

				return echo.ErrBadRequest
			}

			if claims.OrganizationID != orgID {
				l.WithFields(logrus.Fields{
					"Name":           claims.Name,
					"UserID":         claims.UserID,
					"OrganizationID": claims.OrganizationID,
					"Role":           claims.Role,
				}).Debugf("user %d attempted to perform action in organization %d without being a member", claims.UserID, claims.OrganizationID)

				return echo.ErrForbidden
			}

			return next(c)
		}
	}
}
