package match

import (
	"context"
	"foosball/internal/player"

	"github.com/pkg/errors"
)

type Config struct {
	Method string
}

type Service interface {
	CreateMatch(ctx context.Context, teamA, teamB []*player.Player, scoreA, scoreB int, winner string) error
	DeleteMatch(ctx context.Context, match *Match) error
	GetMatchesWithPlayer(ctx context.Context, player *player.Player) ([]*Match, error)
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

func (s *ServiceImpl) CreateMatch(ctx context.Context, teamA []*player.Player, teamB []*player.Player, scoreA int, scoreB int, winner string) error {
	match := &Match{
		TeamA:  []player.Player{},
		TeamB:  []player.Player{},
		ScoreA: scoreA,
		ScoreB: scoreB,
		Winner: winner,
	}

	err := s.repo.CreateMatch(ctx, match)
	if err != nil {
		return errors.Wrap(err, "failed to create match")
	}

	return nil
}

func (s *ServiceImpl) DeleteMatch(ctx context.Context, match *Match) error {
	err := s.repo.DeleteMatch(ctx, match)
	if err != nil {
		return errors.Wrap(err, "failed to delete match")
	}

	return nil
}
func (s *ServiceImpl) GetMatchesWithPlayer(ctx context.Context, player *player.Player) ([]*Match, error) {
	//TODO
	return []*Match{}, nil
}
