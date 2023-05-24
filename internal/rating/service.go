//go:generate mockgen --source=service.go -destination=service_mock.go -package=rating
package rating

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service interface {
	UpdateRatings(ctx context.Context, method Method, draw bool, winners, losers []uint) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) UpdateRatings(ctx context.Context, method Method, draw bool, winners, losers []uint) error {
	if draw {
		return nil // TODO implement draw in rating methods
	}

	winnerRatings, err := s.repo.GetRatingsByUserIDs(ctx, winners)
	if err != nil {
		return errors.Wrap(err, "failed to get winner users")
	}

	loserRatings, err := s.repo.GetRatingsByUserIDs(ctx, losers)
	if err != nil {
		return errors.Wrap(err, "failed to get loser users")
	}

	updatedRatings := make([]Rating, len(winners)+len(losers))

	switch method {
	case Elo:
		updatedRatings = s.calculateNewRatingsElo(winnerRatings, loserRatings)
	case RMS:
		updatedRatings = s.calculateNewRatingsRMS(winnerRatings, loserRatings)
	case Glicko2:
		updatedRatings = s.calculateNewRatingsGlicko2(winnerRatings, loserRatings)
	default:
		logrus.WithField("method", method).Error("unrecognized rating method, defaulting to elo")
		updatedRatings = s.calculateNewRatingsElo(winnerRatings, loserRatings)
	}

	if err := s.repo.UpdateRatings(ctx, updatedRatings); err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
