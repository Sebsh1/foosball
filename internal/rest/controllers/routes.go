package controllers

import (
	"foosball/internal/authentication"
	"foosball/internal/match"
	"foosball/internal/player"
	"foosball/internal/rating"
	"foosball/internal/team"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	logger        *logrus.Entry
	authService   authentication.Service
	playerService player.Service
	matchService  match.Service
	ratingService rating.Service
	teamService   team.Service
}

func Register(
	e *echo.Group,
	logger *logrus.Entry,
	authService authentication.Service,
	playerService player.Service,
	matchService match.Service,
	ratingService rating.Service,
	teamService team.Service,
) {
	h := &Handlers{
		logger:        logger,
		authService:   authService,
		playerService: playerService,
		matchService:  matchService,
		ratingService: ratingService,
		teamService:   teamService,
	}

	// Authentication
	e.POST("/login", h.Login)

	// Players
	e.GET("/player/:id", h.GetPlayer)
	e.POST("/player", h.PostPlayer)
	e.DELETE("/player/:id", h.DeletePlayer)
	e.GET("/player/:id/stats", h.GetPlayerStatistics)

	// Matches
	e.GET("/match/:id", h.GetMatch)
	e.POST("/match", h.PostMatch)

	// Seasons
	e.GET("/season/:id", h.GetSeason)
	e.POST("/season", h.CreateSeason)
	e.DELETE("/season/:id", h.DeleteSeason)

	// Leaderboards
	e.GET("/leaderboard/rating", h.GetLeaderboardRating)
}
