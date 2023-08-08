package controllers

import (
	"matchlog/internal/authentication"
	"matchlog/internal/invite"
	"matchlog/internal/leaderboard"
	"matchlog/internal/match"
	"matchlog/internal/organization"
	"matchlog/internal/rating"
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/middleware"
	"matchlog/internal/statistic"
	"matchlog/internal/user"

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
	statisticService    statistic.Service
	leaderboardService  leaderboard.Service
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
	statisticService statistic.Service,
	leaderboardService leaderboard.Service,
) {
	h := &Handlers{
		logger:              logger,
		authService:         authService,
		userService:         userService,
		organizationService: organizationService,
		inviteService:       inviteService,
		matchService:        matchService,
		ratingService:       ratingService,
		statisticService:    statisticService,
		leaderboardService:  leaderboardService,
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
	e.DELETE("/user", authHandler(h.DeleteUser), authGuard, userGuard)
	e.GET("/user/invites", authHandler(h.GetUserInvites), authGuard, userGuard)
	e.POST("/user/invites/:inviteId/", authHandler(h.RespondToInvite), authGuard, userGuard)

	// Organizations
	e.POST("/organization", authHandler(h.CreateOrganization), authGuard)
	e.PUT("/organization", authHandler(h.UpdateOrganization), authGuard, orginizationGuard, adminGuard)
	e.DELETE("/organization", authHandler(h.DeleteOrganization), authGuard, orginizationGuard, adminGuard)
	e.GET("/organization/users", authHandler(h.GetUsersInOrganization), authGuard, orginizationGuard)
	e.POST("/organization/invite", authHandler(h.InviteUsersToOrganization), authGuard, orginizationGuard, adminGuard)
	e.POST("/organization/users/virtual", authHandler(h.AddVirtualUserToOrganization), authGuard, orginizationGuard, adminGuard)
	e.POST("/organization/users/:userId/virtual/:virtualUserId", authHandler(h.TransferVirtualUserToUser), authGuard, orginizationGuard, adminGuard)
	e.DELETE("/organization/users/:userId", authHandler(h.RemoveUserFromOrganization), authGuard, orginizationGuard, adminGuard)
	e.PUT("/organization/users/:userId", authHandler(h.UpdateUserRole), authGuard, orginizationGuard, adminGuard)
	e.GET("/organization/top/:topX/measures/:leaderboardType", authHandler(h.GetLeaderboard), authGuard, orginizationGuard)
	e.POST("/organization/matches", authHandler(h.PostMatch), authGuard, orginizationGuard)
}
