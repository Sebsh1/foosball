package statistic

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	UpdateStatisticsByUserIDs(ctx context.Context, userIDs []uint, result MatchResult) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) UpdateStatisticsByUserIDs(ctx context.Context, userIDs []uint, result MatchResult) error {
	var updatedStatistics []Statistic

	for _, userID := range userIDs {
		stats, err := s.repo.GetStatisticByUserID(ctx, userID)
		if err != nil {
			return err
		}

		switch result {
		case ResultWin:
			stats.Wins++
			stats.WinStreak++
			stats.LossStreak = 0
		case ResultLoss:
			stats.Losses++
			stats.WinStreak = 0
			stats.LossStreak++
		case ResultDraw:
			stats.Draws++
			stats.WinStreak = 0
			stats.LossStreak = 0
		}

		updatedStatistics = append(updatedStatistics, *stats)
	}

	if err := s.repo.UpdateStatistics(ctx, updatedStatistics); err != nil {
		return errors.Wrap(err, "failed to update statistics")
	}

	return nil
}
