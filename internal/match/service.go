package match

import (
	"context"
	"foosball/internal/models"
	"foosball/internal/player"

	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

type Config struct {
	Method string
}

type Service interface {
	CreateMatch(ctx context.Context, teamA, teamB []*models.Player, goalsA, goalsB int) error
	GetMatch(ctx context.Context, id uint) (*models.Match, error)
	DeleteMatch(ctx context.Context, match *models.Match) error
	GetMatchesWithPlayer(ctx context.Context, player *models.Player) ([]*models.Match, error)
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

func (s *ServiceImpl) CreateMatch(ctx context.Context, teamA, teamB []*models.Player, goalsA, goalsB int) error {
	teamAIDs := make([]uint, len(teamA))
	for i, p := range teamA {
		teamAIDs[i] = p.ID
	}

	teamBIDs := make([]uint, len(teamB))
	for i, p := range teamB {
		teamAIDs[i] = p.ID
	}

	match := &models.Match{
		TeamA:  datatypes.JSONType[[]uint]{Data: teamAIDs},
		TeamB:  datatypes.JSONType[[]uint]{Data: teamBIDs},
		GoalsA: goalsA,
		GoalsB: goalsB,
	}

	err := s.repo.CreateMatch(ctx, match)
	if err != nil {
		return errors.Wrap(err, "failed to create match")
	}

	//TODO add match to players

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
func (s *ServiceImpl) GetMatchesWithPlayer(ctx context.Context, player *models.Player) ([]*models.Match, error) {
	matches, err := s.repo.GetMatchesWithPlayer(ctx, player)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get matches with player")
	}

	return matches, nil
}
