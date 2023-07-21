package leaderboard

import (
	"context"
	"matchlog/internal/rating"
	"matchlog/internal/statistic"
	"matchlog/internal/user"

	"github.com/pkg/errors"
)

type Service interface {
	GetLeaderboard(ctx context.Context, organizationID uint, topX int, leaderboardType LeaderboardType) (*Leaderboard, error)
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

func (s *ServiceImpl) GetLeaderboard(ctx context.Context, organizationID uint, topX int, leaderboardType LeaderboardType) (*Leaderboard, error) {
	var userIDs []uint
	var values []float64

	usersInOrg, err := s.userService.GetUsersInOrganization(ctx, organizationID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get users in organization %d", organizationID)
	}

	userIDsInOrg := s.getUserIDsFromUsers(usersInOrg)

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

	placements := make([]Placement, len(users))
	for i, u := range users {
		placements[i] = Placement{
			Value:  values[i],
			UserID: u.ID,
			Name:   u.Name,
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

func (s *ServiceImpl) getUserIDsFromUsers(users []user.User) []uint {
	ids := make([]uint, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}
	return ids
}

func (s *ServiceImpl) convertIntToFloat64(values []int) []float64 {
	floats := make([]float64, len(values))
	for i, v := range values {
		floats[i] = float64(v)
	}
	return floats
}
