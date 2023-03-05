package controllers

import (
	"foosball/internal/match"
	"foosball/internal/player"
	"foosball/internal/rating"
	"foosball/internal/team"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	logger        *logrus.Entry
	playerService player.Service
	matchService  match.Service
	ratingService rating.Service
	teamService   team.Service
}

func Register(
	e *echo.Group,
	logger *logrus.Entry,
	playerService player.Service,
	matchService match.Service,
	ratingService rating.Service,
	teamService team.Service,
) {
	h := &Handlers{
		logger:        logger,
		playerService: playerService,
		matchService:  matchService,
		ratingService: ratingService,
		teamService:   teamService,
	}

	e.GET("/player/:id", h.GetPlayer)
	e.POST("/player", h.PostPlayer)
	e.DELETE("/player/:id", h.DeletePlayer)
	e.GET("/player/:id/stats", h.GetPlayerStatistics)
	e.GET("/player/:id/matches", h.GetMatchesWithPlayer)

	e.GET("/match/:id", h.GetMatch)
	e.POST("/match", h.PostMatch)

	e.GET("/leaderboard/rating", h.GetLeaderboardRating)
}
