package leaderboard

import (
	"context"
	"foosball/internal/rating"
	"foosball/internal/statistic"
	"foosball/internal/user"

	"github.com/pkg/errors"
)

type Service interface {
	GetLeaderboard(ctx context.Context, leaderboardType LeaderboardType, topX int) (*Leaderboard, error)
}

type ServiceImpl struct {
	userService      user.Service
	ratingService    rating.Service
	statisticService statistic.Service
}

func NewService(userService user.Service, ratingService rating.Service, statisticService statistic.Service) Service {
	return &ServiceImpl{
		userService:      userService,
		ratingService:    ratingService,
		statisticService: statisticService,
	}
}

func (s *ServiceImpl) GetLeaderboard(ctx context.Context, leaderboardType LeaderboardType, topX int) (*Leaderboard, error) {
	var userIDs []uint
	var values []float64

	switch leaderboardType {
	case TypeWins:
		ids, wins, err := s.statisticService.GetTopXUserIDsByWins(ctx, topX)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by wins", topX)
		}

		userIDs = ids
		values = convertIntsToFloat64s(wins)
	case TypeWinStreak:
		ids, winstreaks, err := s.statisticService.GetTopXUserIDsByWinStreak(ctx, topX)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by win streaks", topX)
		}

		userIDs = ids
		values = convertIntsToFloat64s(winstreaks)
	case TypeLossStreak:
		ids, lossStreaks, err := s.statisticService.GetTopXUserIDsByLossStreak(ctx, topX)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by loss streaks", topX)
		}

		userIDs = ids
		values = convertIntsToFloat64s(lossStreaks)
	case TypeWinLossRatio:
		ids, ratios, err := s.statisticService.GetTopXUserIDsByWinLossRatio(ctx, topX)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by win-loss ratio", topX)
		}

		userIDs = ids
		values = ratios
	case TypeRating:
		ids, ratings, err := s.ratingService.GetTopXUserIDsByRating(ctx, topX)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by rating", topX)
		}

		userIDs = ids
		values = convertIntsToFloat64s(ratings)
	case TypeMatchesPlayed:
		ids, matches, err := s.statisticService.GetTopXUserIDsByMatchesPlayed(ctx, topX)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by matches played", topX)
		}

		userIDs = ids
		values = convertIntsToFloat64s(matches)
	default:
		return nil, errors.Errorf("unknown leaderboard type: %s", leaderboardType)
	}

	users, err := s.userService.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	placements := make([]Placement, len(users))
	for i, u := range users {
		placements[i] = Placement{
			Position: i + 1,
			Value:    values[i],
			UserID:   u.ID,
			Name:     u.Name,
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get leaderboard")
	}

	lboard := &Leaderboard{
		Type:       leaderboardType,
		Placements: placements,
	}

	return lboard, nil
}

func convertIntsToFloat64s(ints []int) []float64 {
	floats := make([]float64, len(ints))
	for i, v := range ints {
		floats[i] = float64(v)
	}
	return floats
}
