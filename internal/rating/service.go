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
	Elo      Method = "elo"
	Weighted Method = "weighted"
	RMS      Method = "rms"
	Glicko2  Method = "glicko2"
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

	winnerUsers, err := s.userService.GetUsers(ctx, winners)
	if err != nil {
		return errors.Wrap(err, "failed to get winner users")
	}

	loserUsers, err := s.userService.GetUsers(ctx, losers)
	if err != nil {
		return errors.Wrap(err, "failed to get loser users")
	}

	newRatingsWinners := make([]int, len(winners))
	newRatingsLosers := make([]int, len(losers))

	switch s.config.Method {
	case Elo:
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsElo(winnerUsers, loserUsers)
	case Weighted:
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsWeighted(winnerUsers, loserUsers)
	case RMS:
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsRMS(winnerUsers, loserUsers)
	case Glicko2:
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsGlicko2(winnerUsers, loserUsers)
	default:
		logrus.WithField("method", s.config.Method).Error("unrecognized rating method, defaulting to elo")
		newRatingsWinners, newRatingsLosers = s.calculateNewRatingsElo(winnerUsers, loserUsers)
	}

	userIDs := append(winners, losers...)
	ratings := append(newRatingsWinners, newRatingsLosers...)
	if err := s.userService.UpdateRatings(ctx, userIDs, ratings); err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
