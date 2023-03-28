package match

import (
	"context"
	"foosball/internal/models"
	"foosball/internal/player"

	"github.com/pkg/errors"
)

type Config struct {
	Method string
}

type Service interface {
	CreateMatch(ctx context.Context, teamA, teamB models.Team, goalsA, goalsB int) error
	GetMatch(ctx context.Context, id uint) (*models.Match, error)
	DeleteMatch(ctx context.Context, match *models.Match) error
	GetMatchesWithPlayerID(ctx context.Context, id uint) ([]*models.Match, error)
}

type ServiceImpl struct {
	repo          Repository
	playerService player.Service
}

func NewService(repo Repository, playerService player.Service) Service {
	return &ServiceImpl{
		repo:          repo,
		playerService: playerService,
	}
}

func (s *ServiceImpl) CreateMatch(ctx context.Context, teamA, teamB models.Team, goalsA, goalsB int) error {
	match := &models.Match{
		TeamAID: teamA.ID,
		TeamBID: teamB.ID,
		TeamA:   teamA,
		TeamB:   teamA,
		GoalsA:  goalsA,
		GoalsB:  goalsB,
	}

	err := s.repo.CreateMatch(ctx, match)
	if err != nil {
		return errors.Wrap(err, "failed to create match")
	}

	return nil
}

func (s *ServiceImpl) GetMatch(ctx context.Context, id uint) (*models.Match, error) {
	match, err := s.repo.GetMatch(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get match")
	}

	return match, nil
}

func (s *ServiceImpl) DeleteMatch(ctx context.Context, match *models.Match) error {
	err := s.repo.DeleteMatch(ctx, match)
	if err != nil {
		return errors.Wrap(err, "failed to delete match")
	}

	return nil
}
func (s *ServiceImpl) GetMatchesWithPlayerID(ctx context.Context, id uint) ([]*models.Match, error) {
	matches, err := s.repo.GetMatchesWithPlayerID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get matches with player id %d", id)
	}

	return matches, nil
}
