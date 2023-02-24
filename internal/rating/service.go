package rating

import (
	"context"
	"fmt"
	"foosball/internal/player"

	"github.com/pkg/errors"
)

type Team int

const (
	teamA Team = iota
	teamB
)

type Method int

const (
	Elo Method = iota
	Weighted
	RMS
)

func (m Method) String() string {
	return []string{"elo", "weighted", "rms"}[m]
}

type Config struct {
	Method string `validate:"required" default:"elo"` // One of elo, weighted or rms
}

type Service interface {
	UpdateRatings(ctx context.Context, teamA []*player.Player, teamB []*player.Player, winner Team) error
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

func (s *ServiceImpl) UpdateRatings(ctx context.Context, teamA []*player.Player, teamB []*player.Player, winner Team) error {
	teamAChanges := make([]int, len(teamA))
	teamBChanges := make([]int, len(teamB))

	switch s.config.Method {
	case Elo.String():
		teamAChanges, teamBChanges = s.calculateRatingChangesElo(teamA, teamB)
	case Weighted.String():
		teamAChanges, teamBChanges = s.calculateRatingChangesWeighted(teamA, teamB)
	case RMS.String():
		teamAChanges, teamBChanges = s.calculateRatingChangesRMS(teamA, teamB)
	default:
		return errors.New(fmt.Sprintf("unrecognized rating method: %s", s.config.Method))
	}

	players := append(teamA, teamB...)
	ratings := append(teamAChanges, teamBChanges...)
	err := s.playerService.UpdatePlayers(ctx, players, ratings)
	if err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
