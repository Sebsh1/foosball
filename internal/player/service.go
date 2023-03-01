package player

import (
	"context"
	"foosball/internal/models"

	"github.com/pkg/errors"
)

type Service interface {
	GetPlayer(ctx context.Context, id uint) (*models.Player, error)
	GetPlayers(ctx context.Context, ids []uint) (*[]models.Player, error)
	CreatePlayer(ctx context.Context, name string) error
	DeletePlayer(ctx context.Context, id uint) error
	UpdatePlayers(ctx context.Context, players []*models.Player, ratings []int) error
	GetTopPlayersByRating(ctx context.Context, topX int) ([]*models.Player, error)
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetPlayer(ctx context.Context, id uint) (*models.Player, error) {
	player, err := s.repo.GetPlayer(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get player")
	}

	return player, nil
}

func (s *ServiceImpl) GetPlayers(ctx context.Context, ids []uint) (*[]models.Player, error) {
	player, err := s.repo.GetPlayers(ctx, ids)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get players")
	}

	return player, nil
}

func (s *ServiceImpl) CreatePlayer(ctx context.Context, name string) error {
	player := &models.Player{
		Name:   name,
		Rating: 1000,
	}

	err := s.repo.CreatePlayer(ctx, player)
	if err != nil {
		return errors.Wrap(err, "failed to create player")
	}

	return nil
}

func (s *ServiceImpl) DeletePlayer(ctx context.Context, id uint) error {
	player, err := s.repo.GetPlayer(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get player")
	}

	err = s.repo.DeletePlayer(ctx, player)
	if err != nil {
		return errors.Wrap(err, "failed to delete player")
	}

	return nil
}

func (s *ServiceImpl) UpdatePlayers(ctx context.Context, players []*models.Player, ratings []int) error {
	for i, p := range players {
		p.Rating += ratings[i]
	}

	err := s.repo.UpdatePlayers(ctx, players)
	if err != nil {
		return errors.Wrap(err, "failed to update players")
	}

	return nil
}

func (s *ServiceImpl) GetTopPlayersByRating(ctx context.Context, topX int) ([]*models.Player, error) {
	players, err := s.repo.GetTopPlayersByRating(ctx, topX)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get players")
	}

	return players, nil
}
