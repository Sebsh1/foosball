package leaderboard

import (
	"context"
	"matchlog/internal/club"
	"matchlog/internal/rating"
	"matchlog/internal/statistic"
	"matchlog/internal/user"

	"github.com/pkg/errors"
)

type Service interface {
	GetLeaderboard(ctx context.Context, ClubID uint, topX int, leaderboardType LeaderboardType) (*Leaderboard, error)
}

type ServiceImpl struct {
	clubService      club.Service
	userService      user.Service
	ratingService    rating.Service
	statisticService statistic.Service
}

func NewService(clubService club.Service, userService user.Service, ratingService rating.Service, statisticService statistic.Service) Service {
	return &ServiceImpl{
		clubService:      clubService,
		userService:      userService,
		ratingService:    ratingService,
		statisticService: statisticService,
	}
}

func (s *ServiceImpl) GetLeaderboard(ctx context.Context, ClubID uint, topX int, leaderboardType LeaderboardType) (*Leaderboard, error) {
	var userIDs []uint
	var values []float64

	userIDsInOrg, err := s.clubService.GetUserIDsInClub(ctx, ClubID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get users in Club %d", ClubID)
	}

	switch leaderboardType {
	case TypeWins:
		ids, wins, err := s.statisticService.GetTopXAmongUserIDsByMeasure(ctx, topX, userIDsInOrg, statistic.MeasureWins)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by wins", topX)
		}

		userIDs = ids
		values = s.convertIntToFloat64(wins)
	case TypeStreak:
		ids, winstreaks, err := s.statisticService.GetTopXAmongUserIDsByMeasure(ctx, topX, userIDsInOrg, statistic.MeasureStreak)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by streak", topX)
		}

		userIDs = ids
		values = s.convertIntToFloat64(winstreaks)
	case TypeRating:
		ids, ratings, err := s.ratingService.GetTopXAmongUserIDsByRating(ctx, topX, userIDsInOrg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get top %d userIDs by rating", topX)
		}

		userIDs = ids
		values = s.convertIntToFloat64(ratings)
	default:
		return nil, errors.Errorf("unknown leaderboard type: %s", leaderboardType)
	}

	users, err := s.userService.GetUsers(ctx, userIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	entries := make([]Entry, len(users))
	for i, u := range users {
		entries[i] = Entry{
			Value:  values[i],
			UserID: u.ID,
			Name:   u.Name,
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get leaderboard")
	}

	lboard := &Leaderboard{
		Type:    leaderboardType,
		Entries: entries,
	}

	return lboard, nil
}

func (s *ServiceImpl) convertIntToFloat64(values []int) []float64 {
	floats := make([]float64, len(values))
	for i, v := range values {
		floats[i] = float64(v)
	}
	return floats
}
