package statistic

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetTopXUserIDsByWins(ctx context.Context, topX int) (userIDs []uint, wins []int, err error)
	GetTopXUserIDsByWinStreak(ctx context.Context, topX int) (userIDs []uint, winStreaks []int, err error)
	GetTopXUserIDsByLossStreak(ctx context.Context, topX int) (userIDs []uint, lossStreaks []int, err error)
	GetTopXUserIDsByWinLossRatio(ctx context.Context, topX int) (userIDs []uint, winLossRatios []float64, err error)
	GetTopXUserIDsByMatchesPlayed(ctx context.Context, topX int) (userIDs []uint, matchesPlayed []int, err error)
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

func (s *ServiceImpl) GetTopXUserIDsByWins(ctx context.Context, topX int) ([]uint, []int, error) {
	userIDs, wins, err := s.repo.GetTopXUserIDsByWins(ctx, topX)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get top %d userIDs by wins", topX)
	}

	return userIDs, wins, nil
}

func (s *ServiceImpl) GetTopXUserIDsByWinStreak(ctx context.Context, topX int) ([]uint, []int, error) {
	userIDs, winStreaks, err := s.repo.GetTopXUserIDsByWinStreaks(ctx, topX)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get top %d userIDs by win streaks", topX)
	}

	return userIDs, winStreaks, nil
}

func (s *ServiceImpl) GetTopXUserIDsByLossStreak(ctx context.Context, topX int) ([]uint, []int, error) {
	userIDs, lossStreaks, err := s.repo.GetTopXUserIDsByLossStreaks(ctx, topX)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get top %d userIDs by loss streaks", topX)
	}

	return userIDs, lossStreaks, nil
}

func (s *ServiceImpl) GetTopXUserIDsByWinLossRatio(ctx context.Context, topX int) ([]uint, []float64, error) {
	userIDs, winLossRatios, err := s.repo.GetTopXUserIDsByWinLossRatios(ctx, topX)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get top %d userIDs by win loss ratio", topX)
	}

	return userIDs, winLossRatios, nil
}

func (s *ServiceImpl) GetTopXUserIDsByMatchesPlayed(ctx context.Context, topX int) ([]uint, []int, error) {
	userIDs, matchesPlayed, err := s.repo.GetTopXUserIDsByMatchesPlayed(ctx, topX)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get top %d userIDs by matches played", topX)
	}

	return userIDs, matchesPlayed, nil
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
			stats.LossStreak = 0
			stats.Wins++
			stats.WinStreak++
			stats.MatchesPlayed++
			stats.WinLossRatio = float64(stats.Wins) / float64(stats.Losses)
		case ResultLoss:
			stats.WinStreak = 0
			stats.Losses++
			stats.LossStreak++
			stats.MatchesPlayed++
			stats.WinLossRatio = float64(stats.Wins) / float64(stats.Losses)
		case ResultDraw:
			stats.WinStreak = 0
			stats.LossStreak = 0
			stats.Draws++
			stats.MatchesPlayed++
		}

		updatedStatistics = append(updatedStatistics, *stats)
	}

	if err := s.repo.UpdateStatistics(ctx, updatedStatistics); err != nil {
		return errors.Wrap(err, "failed to update statistics")
	}

	return nil
}
