//go:generate mockgen --source=service.go -destination=service_mock.go -package=rating
package rating

import (
	"context"
	"fmt"
	"foosball/internal/player"
	"foosball/internal/team"

	"github.com/pkg/errors"
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
	UpdateRatings(ctx context.Context, winners, losers *team.Team) error
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

func (s *ServiceImpl) UpdateRatings(ctx context.Context, winners, losers *team.Team) error {
	newRatingsWinners := make([]int, len(winners.Players))
	newRatingsLosers := make([]int, len(losers.Players))

	switch s.config.Method {
	case Elo.String():
		newRatingsWinners, newRatingsLosers = s.calculateRatingChangesElo(winners, losers)
	case RMS.String():
		newRatingsWinners, newRatingsLosers = s.calculateRatingChangesRMS(winners, losers)
	default:
		return errors.New(fmt.Sprintf("unrecognized rating method: %s", s.config.Method))
	}

	players := append(winners.Players, losers.Players...)
	ratings := append(newRatingsWinners, newRatingsLosers...)
	err := s.playerService.UpdatePlayers(ctx, players, ratings)
	if err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
