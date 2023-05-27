package statistic

import (
	"context"

	"github.com/pkg/errors"
)

type Measure string

const (
	MeasureWins          Measure = "wins"
	MeasureWinStreak     Measure = "win-streak"
	MeasureLossStreak    Measure = "loss-streak"
	MeasureWinLossRatio  Measure = "win-loss-ratio"
	MeasureMatchesPlayed Measure = "matches-played"
)

type Service interface {
	GetTopXAmongUserIDsByMeasure(ctx context.Context, topX int, usersIDs []uint, measure Measure) (userIDs []uint, values []float64, err error)
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

func (s *ServiceImpl) GetTopXAmongUserIDsByMeasure(ctx context.Context, topX int, userIDs []uint, measure Measure) ([]uint, []float64, error) {
	var topXUserIDs []uint
	var values []float64

	switch measure {
	case MeasureWins:
		ids, wins, err := s.repo.GetTopXAmongUserIDsByWins(ctx, topX, userIDs)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by wins", topX)
		}

		topXUserIDs = ids
		values = s.convertIntsToFloat64s(wins)
	case MeasureWinStreak:
		ids, winStreaks, err := s.repo.GetTopXAmongUserIDsByWinStreaks(ctx, topX, userIDs)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by win streaks", topX)
		}

		topXUserIDs = ids
		values = s.convertIntsToFloat64s(winStreaks)
	case MeasureLossStreak:
		ids, lossStreaks, err := s.repo.GetTopXAmongUserIDsByLossStreaks(ctx, topX, userIDs)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by loss streaks", topX)
		}

		topXUserIDs = ids
		values = s.convertIntsToFloat64s(lossStreaks)
	case MeasureWinLossRatio:
		ids, winLossRatios, err := s.repo.GetTopXAmongUserIDsByWinLossRatios(ctx, topX, userIDs)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by win loss ratio", topX)
		}

		topXUserIDs = ids
		values = winLossRatios
	case MeasureMatchesPlayed:
		ids, matchesPlayed, err := s.repo.GetTopXAmongUserIDsByMatchesPlayed(ctx, topX, userIDs)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get top %d users by matches played", topX)
		}

		topXUserIDs = ids
		values = s.convertIntsToFloat64s(matchesPlayed)
	}

	return topXUserIDs, values, nil
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

func (s *ServiceImpl) convertIntsToFloat64s(ints []int) []float64 {
	floats := make([]float64, len(ints))
	for i, v := range ints {
		floats[i] = float64(v)
	}
	return floats
}
