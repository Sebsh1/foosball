package team

import (
	"context"
	"foosball/internal/models"

	"github.com/pkg/errors"
)

type Service interface {
	GetTeam(ctx context.Context, id uint) (*models.Team, error)
	CreateTeam(ctx context.Context, players []*models.Player) error
	DeleteTeam(ctx context.Context, id uint) error
	UpdateTeam(ctx context.Context, team *models.Team) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetTeam(ctx context.Context, id uint) (*models.Team, error) {
	player, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get team")
	}

	return player, nil
}

func (s *ServiceImpl) CreateTeam(ctx context.Context, players []*models.Player) error {
	team := &models.Team{
		Players: players,
	}
	err := s.repo.CreateTeam(ctx, team)
	if err != nil {
		return errors.Wrap(err, "failed to create team")
	}

	return nil
}

func (s *ServiceImpl) DeleteTeam(ctx context.Context, id uint) error {
	team, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get team")
	}

	err = s.repo.DeleteTeam(ctx, team)
	if err != nil {
		return errors.Wrap(err, "failed to delete team")
	}

	return nil
}

func (s *ServiceImpl) UpdateTeam(ctx context.Context, team *models.Team) error {
	err := s.repo.UpdateTeam(ctx, team)
	if err != nil {
		return errors.Wrap(err, "failed to update team")
	}

	return nil
}
