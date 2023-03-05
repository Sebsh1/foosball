package rating

import (
	"context"
	"fmt"
	"foosball/internal/player"
	"foosball/internal/team"

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
	UpdateRatings(ctx context.Context, teamA, teamB *team.Team, winner Team) error
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

func (s *ServiceImpl) UpdateRatings(ctx context.Context, teamA, teamB *team.Team, winner Team) error {
	playersTeamA := teamA.GetPlayers()
	playersTeamB := teamB.GetPlayers()
	newRatingsTeamA := make([]int, len(playersTeamA))
	newRatingsTeamB := make([]int, len(playersTeamB))

	switch s.config.Method {
	case Elo.String():
		newRatingsTeamA, newRatingsTeamB = s.calculateRatingChangesElo(teamA, teamB, winner)
	case RMS.String():
		newRatingsTeamA, newRatingsTeamB = s.calculateRatingChangesRMS(teamA, teamB, winner)
	default:
		return errors.New(fmt.Sprintf("unrecognized rating method: %s", s.config.Method))
	}

	players := append(playersTeamA, playersTeamB...)
	ratings := append(newRatingsTeamA, newRatingsTeamB...)
	err := s.playerService.UpdatePlayerRatings(ctx, players, ratings)
	if err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
