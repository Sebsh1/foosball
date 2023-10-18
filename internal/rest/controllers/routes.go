package controllers

import (
	"matchlog/internal/authentication"
	"matchlog/internal/club"
	"matchlog/internal/leaderboard"
	"matchlog/internal/match"
	"matchlog/internal/rating"
	"matchlog/internal/rest/handlers"
	"matchlog/internal/rest/middleware"
	"matchlog/internal/statistic"
	"matchlog/internal/user"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handlers struct {
	logger             *zap.SugaredLogger
	authService        authentication.Service
	userService        user.Service
	clubService        club.Service
	matchService       match.Service
	ratingService      rating.Service
	statisticService   statistic.Service
	leaderboardService leaderboard.Service
}

func Register(
	e *echo.Group,
	logger *zap.SugaredLogger,
	authService authentication.Service,
	userService user.Service,
	clubService club.Service,
	matchService match.Service,
	ratingService rating.Service,
	statisticService statistic.Service,
	leaderboardService leaderboard.Service,
) {
	h := &Handlers{
		logger:             logger,
		authService:        authService,
		userService:        userService,
		clubService:        clubService,
		matchService:       matchService,
		ratingService:      ratingService,
		statisticService:   statisticService,
		leaderboardService: leaderboardService,
	}

	authHandler := handlers.AuthenticatedHandlerFactory(logger)

	authGuard := middleware.AuthGuard(authService)

	// Authentication
	e.POST("/login", h.Login)
	e.POST("/signup", h.Signup)

	// Users
	e.DELETE("/user", authHandler(h.DeleteUser), authGuard)
	e.GET("/user/invites", authHandler(h.GetUserInvites), authGuard)
	e.POST("/user/invites/:inviteId", authHandler(h.RespondToInvite), authGuard)

	// Clubs
	e.POST("/club", authHandler(h.CreateClub), authGuard)
	e.PUT("/club", authHandler(h.UpdateClub), authGuard)
	e.DELETE("/club", authHandler(h.DeleteClub), authGuard)
	e.GET("/club/users", authHandler(h.GetUsersInClub), authGuard)
	e.POST("/club/invite", authHandler(h.InviteUsersToClub), authGuard)
	e.POST("/club/users/virtual", authHandler(h.AddVirtualUserToClub), authGuard)
	e.POST("/club/users/:userId/virtual/:virtualUserId", authHandler(h.TransferVirtualUserToUser), authGuard)
	e.DELETE("/club/users/:userId", authHandler(h.RemoveUserFromClub), authGuard)
	e.PUT("/club/users/:userId", authHandler(h.UpdateUserRole), authGuard)
	e.GET("/club/top/:topX/measures/:leaderboardType", authHandler(h.GetLeaderboard), authGuard)
	e.POST("/club/matches", authHandler(h.PostMatch), authGuard)
}
