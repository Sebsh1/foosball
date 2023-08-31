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
	"go.uber.org/zap"
)

type Handlers struct {
	logger              *zap.SugaredLogger
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
	logger *zap.SugaredLogger,
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

	// Authentication
	e.POST("/login", h.Login)
	e.POST("/signup", h.Signup)

	// Users
	e.DELETE("/user", authHandler(h.DeleteUser), authGuard)
	e.GET("/user/invites", authHandler(h.GetUserInvites), authGuard)
	e.POST("/user/invites/:inviteId", authHandler(h.RespondToInvite), authGuard)

	// Organizations
	e.POST("/organization", authHandler(h.CreateOrganization), authGuard)
	e.PUT("/organization", authHandler(h.UpdateOrganization), authGuard, adminGuard)
	e.DELETE("/organization", authHandler(h.DeleteOrganization), authGuard, adminGuard)
	e.GET("/organization/users", authHandler(h.GetUsersInOrganization), authGuard)
	e.POST("/organization/invite", authHandler(h.InviteUsersToOrganization), authGuard, adminGuard)
	e.POST("/organization/users/virtual", authHandler(h.AddVirtualUserToOrganization), authGuard, adminGuard)
	e.POST("/organization/users/:userId/virtual/:virtualUserId", authHandler(h.TransferVirtualUserToUser), authGuard, adminGuard)
	e.DELETE("/organization/users/:userId", authHandler(h.RemoveUserFromOrganization), authGuard, adminGuard)
	e.PUT("/organization/users/:userId", authHandler(h.UpdateUserRole), authGuard, adminGuard)
	e.GET("/organization/top/:topX/measures/:leaderboardType", authHandler(h.GetLeaderboard), authGuard)
	e.POST("/organization/matches", authHandler(h.PostMatch), authGuard)
}
