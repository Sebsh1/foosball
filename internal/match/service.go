package match

import (
	"context"
	"foosball/internal/player"
	"foosball/internal/season"

	"github.com/pkg/errors"
)

type Config struct {
	Method string
}

type Service interface {
	GetMatch(ctx context.Context, id uint) (*Match, error)
	GetMatchesWithPlayerID(ctx context.Context, playerID uint) ([]*Match, error)
	CreateMatch(ctx context.Context, teamAID, teamBID uint, goalsA, goalsB int) error
	DeleteMatch(ctx context.Context, id uint) error
}

type ServiceImpl struct {
	repo          Repository
	playerService player.Service
	seasonService season.Service
}

func NewService(repo Repository, playerService player.Service, seasonService season.Service) Service {
	return &ServiceImpl{
		repo:          repo,
		playerService: playerService,
		seasonService: seasonService,
	}
}

func (s *ServiceImpl) CreateMatch(ctx context.Context, teamAID, teamBID uint, goalsA, goalsB int) error {
	seasonID, err := s.seasonService.GetCurrentSeasonID(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get current season id")
	}

	match := &Match{
		SeasonID: seasonID,
		TeamAID:  teamAID,
		TeamBID:  teamBID,
		GoalsA:   goalsA,
		GoalsB:   goalsB,
	}

	if err = s.repo.CreateMatch(ctx, match); err != nil {
		if err == ErrDuplicateEntry {
			return err
		}
		return errors.Wrap(err, "failed to create match")
	}

	return nil
}

func (s *ServiceImpl) GetMatch(ctx context.Context, id uint) (*Match, error) {
	match, err := s.repo.GetMatch(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get match")
	}

	return match, nil
}

func (s *ServiceImpl) DeleteMatch(ctx context.Context, id uint) error {
	match, err := s.repo.GetMatch(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get match")
	}

	err = s.repo.DeleteMatch(ctx, match)
	if err != nil {
		return errors.Wrap(err, "failed to delete match")
	}

	return nil
}
func (s *ServiceImpl) GetMatchesWithPlayerID(ctx context.Context, playerID uint) ([]*Match, error) {
	matches, err := s.repo.GetMatchesWithPlayerID(ctx, playerID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get matches with player id %d", playerID)
	}

	return matches, nil
}
