package statistic

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error)
	GetTopXAmongUserIDsByMeasure(ctx context.Context, topX int, userIDs []uint, measure Measure) (topXUserIDs []uint, values []int, err error)
	UpdateStatisticsByUserIDs(ctx context.Context, userIDs []uint, result MatchResult) error
	TransferStatistics(ctx context.Context, fromUserID, toUserID uint) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error) {
	stats, err := s.repo.GetStatisticByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statistics for user %d", userID)
	}

	return stats, nil
}

func (s *ServiceImpl) GetTopXAmongUserIDsByMeasure(ctx context.Context, topX int, userIDs []uint, measure Measure) ([]uint, []int, error) {
	var topXUserIDs []uint
	var values []int

	switch measure {
	case MeasureWins:
		ids, wins, err := s.repo.GetTopXAmongUserIDsByWins(ctx, topX, userIDs)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by wins", topX)
		}

		topXUserIDs = ids
		values = wins
	case MeasureStreak:
		ids, streaks, err := s.repo.GetTopXAmongUserIDsByStreak(ctx, topX, userIDs)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by win streaks", topX)
		}

		topXUserIDs = ids
		values = streaks
	default:
		return nil, nil, errors.Errorf("unsupported measure: %s", measure)
	}

	return topXUserIDs, values, nil
}

func (s *ServiceImpl) UpdateStatisticsByUserIDs(ctx context.Context, userIDs []uint, result MatchResult) error {
	oldStatistics, err := s.repo.GetStatisticsByUserIDs(ctx, userIDs)
	if err != nil {
		return errors.Wrapf(err, "failed to get statistics for users %v", userIDs)
	}

	updatedStatistics := make([]Statistic, len(userIDs))

	for i := range userIDs {
		stats := oldStatistics[i]

		switch result {
		case ResultWin:
			stats.Wins++
			if stats.Streak >= 0 {
				stats.Streak++
			} else {
				stats.Streak = 1
			}
		case ResultLoss:
			stats.Losses++
			if stats.Streak <= 0 {
				stats.Streak--
			} else {
				stats.Streak = -1
			}
		case ResultDraw:
			stats.Draws++
			stats.Streak = 0
		}

		updatedStatistics[i] = *stats
	}

	if err := s.repo.UpdateStatistics(ctx, updatedStatistics); err != nil {
		return errors.Wrap(err, "failed to update statistics")
	}

	return nil
}

func (s *ServiceImpl) TransferStatistics(ctx context.Context, fromUserID, toUserID uint) error {
	fromStats, err := s.repo.GetStatisticByUserID(ctx, fromUserID)
	if err != nil {
		return errors.Wrapf(err, "failed to get statistics for user %d", fromUserID)
	}

	toStats, err := s.repo.GetStatisticByUserID(ctx, toUserID)
	if err != nil {
		return errors.Wrapf(err, "failed to get statistics for user %d", toUserID)
	}

	fromStats.UserID = toUserID
	toStats.UserID = fromUserID

	if err := s.repo.UpdateStatistics(ctx, []Statistic{*fromStats, *toStats}); err != nil {
		return errors.Wrap(err, "failed to update statistics")
	}

	return nil
}
