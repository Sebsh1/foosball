package team

import (
	"context"
	"fmt"
	"foosball/internal/player"

	"github.com/pkg/errors"
)

type Service interface {
	GetTeam(ctx context.Context, id uint) (*Team, error)
	CreateTeam(ctx context.Context, players []*player.Player) error
	DeleteTeam(ctx context.Context, id uint) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetTeam(ctx context.Context, id uint) (*Team, error) {
	player, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get team")
	}

	return player, nil
}

func (s *ServiceImpl) CreateTeam(ctx context.Context, players []*player.Player) error {
	team := &Team{}
	switch len(players) {
	case 3:
		team.PlayerC = players[2]
		fallthrough
	case 2:
		team.PlayerB = players[1]
		fallthrough
	case 1:
		team.PlayerA = players[0]
	default:
		return errors.New(fmt.Sprintf("unsupported number of players on team: %d", len(players)))
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
