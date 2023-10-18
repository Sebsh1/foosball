package league

import (
	"context"
)

type Service interface {
	CreateLeague(ctx context.Context, name string) error
	GetLeague(ctx context.Context, id uint) (*League, error)
	UpdateLeague(ctx context.Context, id uint, name string) error
	DeleteLeague(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateLeague(ctx context.Context, name string) error {
	league := &League{
		Name: name,
	}

	return s.repo.CreateLeague(ctx, league)
}

func (s *service) GetLeague(ctx context.Context, id uint) (*League, error) {
	return s.repo.GetLeague(ctx, id)
}

func (s *service) UpdateLeague(ctx context.Context, id uint, name string) error {
	league, err := s.GetLeague(ctx, id)
	if err != nil {
		return err
	}

	league.Name = name

	return s.repo.UpdateLeague(ctx, league)
}

func (s *service) DeleteLeague(ctx context.Context, id uint) error {
	return s.repo.DeleteLeague(ctx, id)
}
