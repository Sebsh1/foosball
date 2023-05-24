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
	adminGuard := middleware.AdminGuard(logger)
	orginizationGuard := middleware.OrganizationGuard(logger)
	userGuard := middleware.UserGuard(logger)

	// Authentication
	e.POST("/login", h.Login)
	e.POST("/signup", h.Signup)

	// Users
	userGroup := e.Group("/user", authGuard, userGuard)
	userGroup.DELETE("/:userId", authHandler(h.DeleteUser))
	userGroup.GET("/:userId/invites", authHandler(h.GetUserInvites))
	userGroup.POST("/:userId/invite/:inviteId/accept", authHandler(h.AcceptInvite))
	userGroup.POST("/:userId/invite/:inviteId/decline", authHandler(h.DeclineInvite))

	// Organizations
	orgGroup := e.Group("/organization", authGuard)
	orgGroup.GET("/:orgId/users", authHandler(h.GetUsersInOrganization), orginizationGuard)
	orgGroup.DELETE("/:orgId", authHandler(h.DeleteOrganization), orginizationGuard, adminGuard)
	orgGroup.POST("", authHandler(h.CreateOrganization))
	orgGroup.POST("/:orgId", authHandler(h.UpdateOrganization), orginizationGuard, adminGuard)
	orgGroup.POST("/:orgId/invite/", authHandler(h.InviteUserToOrganization), orginizationGuard, adminGuard)
	orgGroup.POST("/:orgId/user/:userId/remove", authHandler(h.RemoveUserFromOrganization), orginizationGuard, adminGuard)
	orgGroup.POST("/:orgId/user/:userId/admin", authHandler(h.SetUserAsAdmin), orginizationGuard, adminGuard)
	orgGroup.POST("/:orgId/user/:userId/admin/remove", authHandler(h.RemoveAdminFromUser), orginizationGuard, adminGuard)

	// Matches
	matchGroup := e.Group("/match", authGuard)
	matchGroup.POST("", authHandler(h.PostMatch), orginizationGuard)
}
