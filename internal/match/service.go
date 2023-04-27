package match

import (
	"context"
	"encoding/json"
)

type Service interface {
	DetermineResult(ctx context.Context, teamA, teamB []uint, sets []Set) (draw bool, winner []uint, loser []uint)
	CreateMatch(ctx context.Context, teamA, teamB []uint, sets []Set) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) CreateMatch(ctx context.Context, teamA, teamB []uint, sets []Set) error {
	marshalledTeamA, err := json.Marshal(teamA)
	if err != nil {
		return err
	}

	marshalledTeamB, err := json.Marshal(teamB)
	if err != nil {
		return err
	}

	match := &Match{
		TeamA: marshalledTeamA,
		TeamB: marshalledTeamB,
		Sets:  sets,
	}

	if err := s.repo.CreateMatch(ctx, match); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) DetermineResult(ctx context.Context, teamA, teamB []uint, sets []Set) (bool, []uint, []uint) {
	teamASetWins := 0
	teamBSetWins := 0

	for _, set := range sets {
		if set.PointsA > set.PointsB {
			teamASetWins++
		} else if set.PointsB > set.PointsA {
			teamBSetWins++
		}
	}

	if teamASetWins > teamBSetWins {
		return false, teamA, teamB
	} else if teamBSetWins > teamASetWins {
		return false, teamB, teamA
	} else {
		return true, nil, nil
	}

}
