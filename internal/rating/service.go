package rating

import (
	"context"
	"fmt"
	"foosball/internal/models"
	"foosball/internal/player"

	"github.com/pkg/errors"
)

type Team int

const (
	TeamA Team = iota
	TeamB
)

type Method int

const (
	Elo Method = iota
	Weighted
	RMS
)

func (m Method) String() string {
	return []string{"elo", "rms"}[m]
}

type Config struct {
	Method string `mapstructure:"method" validate:"required" default:"elo"`
}

type Service interface {
	UpdateRatings(ctx context.Context, teamA, teamB *models.Team, winner Team) error
}

type ServiceImpl struct {
	config        Config
	playerService player.Service
}

func NewService(config Config, playerService player.Service) Service {
	return &ServiceImpl{
		config:        config,
		playerService: playerService,
	}
}

func (s *ServiceImpl) UpdateRatings(ctx context.Context, teamA, teamB *models.Team, winner Team) error {
	newRatingsTeamA := make([]int, len(teamA.Players))
	newRatingsTeamB := make([]int, len(teamB.Players))

	switch s.config.Method {
	case Elo.String():
		newRatingsTeamA, newRatingsTeamB = s.calculateRatingChangesElo(teamA, teamB, winner)
	case RMS.String():
		newRatingsTeamA, newRatingsTeamB = s.calculateRatingChangesRMS(teamA, teamB, winner)
	default:
		return errors.New(fmt.Sprintf("unrecognized rating method: %s", s.config.Method))
	}

	players := append(teamA.Players, teamB.Players...)
	ratings := append(newRatingsTeamA, newRatingsTeamB...)
	err := s.playerService.UpdatePlayers(ctx, players, ratings)
	if err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
