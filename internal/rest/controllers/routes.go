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
	userGroup := e.Group("/user", authGuard)
	userGroup.DELETE("/user", authHandler(h.DeleteUser))
	userGroup.GET("/user/invites", authHandler(h.GetUserInvites))
	userGroup.POST("/user/invites/:inviteId", authHandler(h.RespondToInvite))

	// Clubs
	clubGroup := e.Group("/club", authGuard)
	clubGroup.POST("", authHandler(h.CreateClub))
	clubGroup.PUT("", authHandler(h.UpdateClub))
	clubGroup.DELETE("", authHandler(h.DeleteClub))
	clubGroup.GET("/users", authHandler(h.GetUsersInClub))
	clubGroup.POST("/invite", authHandler(h.InviteUsersToClub))
	clubGroup.POST("/users/virtual", authHandler(h.AddVirtualUserToClub))
	//clubGroup.POST("/users/:userId/virtual/:virtualUserId", authHandler(h.TransferVirtualUserToUser))
	clubGroup.DELETE("/users/:userId", authHandler(h.RemoveUserFromClub))
	clubGroup.PUT("/users/:userId", authHandler(h.UpdateUserRole))
	clubGroup.GET("/top/:topX/measures/:leaderboardType", authHandler(h.GetLeaderboard))
	clubGroup.POST("/matches", authHandler(h.PostMatch))
}
