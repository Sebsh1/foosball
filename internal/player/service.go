package player

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetPlayerByName(ctx context.Context, name string) (*Player, error)
	CreatePlayer(ctx context.Context, name string) error
	DeletePlayer(ctx context.Context, name string) error
	UpdatePlayers(ctx context.Context, players []*Player, ratings []int) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetPlayerByName(ctx context.Context, name string) (*Player, error) {
	player, err := s.repo.GetPlayerByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get player")
	}

	return player, nil
}

func (s *ServiceImpl) CreatePlayer(ctx context.Context, name string) error {
	player := &Player{
		Name:   name,
		Rating: 1000,
	}

	err := s.repo.CreatePlayer(ctx, player)
	if err != nil {
		return errors.Wrap(err, "failed to create player")
	}

	return nil
}

func (s *ServiceImpl) DeletePlayer(ctx context.Context, name string) error {
	player, err := s.repo.GetPlayerByName(ctx, name)
	if err != nil {
		return errors.Wrap(err, "failed to get player")
	}

	err = s.repo.DeletePlayer(ctx, player)
	if err != nil {
		return errors.Wrap(err, "failed to delete player")
	}

	return nil
}

func (s *ServiceImpl) UpdatePlayers(ctx context.Context, players []*Player, ratings []int) error {
	for i, p := range players {
		p.Rating += ratings[i]
	}

	err := s.repo.UpdatePlayers(ctx, players)
	if err != nil {
		return errors.Wrap(err, "failed to update players")
	}

	return nil
}
