package controllers

import (
	"foosball/internal/authentication"
	"foosball/internal/invite"
	"foosball/internal/match"
	"foosball/internal/organization"
	"foosball/internal/rating"
	"foosball/internal/rest/handlers"
	"foosball/internal/rest/middleware"
	"foosball/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	logger              *logrus.Entry
	authService         authentication.Service
	userService         user.Service
	organizationService organization.Service
	inviteService       invite.Service
	matchService        match.Service
	ratingService       rating.Service
}

func Register(
	e *echo.Group,
	logger *logrus.Entry,
	authService authentication.Service,
	userService user.Service,
	organizationService organization.Service,
	inviteService invite.Service,
	matchService match.Service,
	ratingService rating.Service,
) {
	h := &Handlers{
		logger:              logger,
		authService:         authService,
		userService:         userService,
		organizationService: organizationService,
		inviteService:       inviteService,
		matchService:        matchService,
		ratingService:       ratingService,
	}

	authHandler := handlers.AuthenticatedHandlerFactory(logger)
	authGuard := middleware.AuthGuard(authService)

	// Authentication
	e.POST("/login", h.Login)
	e.POST("/signup", h.Signup)

	// Users
	userGroup := e.Group("/user", authGuard())
	userGroup.DELETE("", authHandler(h.DeleteUser))
	userGroup.GET("/:userId/invites", authHandler(h.GetUserInvites))
	userGroup.POST("/:userId/invite/:inviteId/accept", authHandler(h.AcceptInvite))
	userGroup.POST("/:userId/invite/:inviteId/decline", authHandler(h.DeclineInvite))

	// Organizations
	orgGroup := e.Group("/organization", authGuard())
	orgGroup.GET("/:orgId/users", authHandler(h.GetUsersInOrganization))
	orgGroup.DELETE("/:orgId", authHandler(h.DeleteOrganization))
	orgGroup.POST("", authHandler(h.CreateOrganization))
	orgGroup.POST("/:orgId", authHandler(h.UpdateOrganization))
	orgGroup.POST("/:orgId/invite/", authHandler(h.InviteUserToOrganization))
	orgGroup.POST("/:orgId/user/:userId/remove", authHandler(h.RemoveUserFromOrganization))
	orgGroup.POST("/:orgId/user/:userId/admin", authHandler(h.SetUserAsAdmin))
	orgGroup.POST("/:orgId/user/:userId/admin/remove", authHandler(h.RemoveAdminFromUser))

	// Matches
	matchGroup := e.Group("/match", authGuard())
	matchGroup.POST("", authHandler(h.PostMatch))
}
