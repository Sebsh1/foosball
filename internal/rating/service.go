//go:generate mockgen --source=service.go -destination=service_mock.go -package=rating
package rating

import (
	"context"
	"foosball/internal/user"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Method string

const (
	Elo     Method = "elo"
	RMS     Method = "rms"
	Glicko2 Method = "glicko2"
)

type Config struct {
	Method Method `mapstructure:"method" validate:"required" default:"elo"`
}

type Service interface {
	UpdateRatings(ctx context.Context, draw bool, winners, losers []uint) error
}

type ServiceImpl struct {
	config      Config
	userService user.Service
}

func NewService(config Config, userService user.Service) Service {
	return &ServiceImpl{
		config:      config,
		userService: userService,
	}
}

func (s *ServiceImpl) UpdateRatings(ctx context.Context, draw bool, winners, losers []uint) error {
	if draw {
		return nil // TODO implement draw in rating methods
	}

	winnerUsersStats, err := s.userService.GetUsersStats(ctx, winners)
	if err != nil {
		return errors.Wrap(err, "failed to get winner users")
	}

	loserUsersStats, err := s.userService.GetUsersStats(ctx, losers)
	if err != nil {
		return errors.Wrap(err, "failed to get loser users")
	}

	newRatingsWinners := make([]int, len(winners))
	newRatingsLosers := make([]int, len(losers))

	switch s.config.Method {
	case Elo:
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsElo(winnerUsersStats, loserUsersStats)
	case RMS:
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsRMS(winnerUsersStats, loserUsersStats)
	case Glicko2:
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsGlicko2(winnerUsersStats, loserUsersStats)
	default:
		logrus.WithField("method", s.config.Method).Error("unrecognized rating method, defaulting to elo")
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsElo(winnerUsersStats, loserUsersStats)
	}

	userIDs := append(winners, losers...)
	ratings := append(newRatingsWinners, newRatingsLosers...)
	if err := s.userService.UpdateRatings(ctx, userIDs, ratings); err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
