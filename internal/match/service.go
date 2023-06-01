package match

import (
	"context"
	"encoding/json"
	"fmt"
)

type Service interface {
	DetermineResult(ctx context.Context, teamA, teamB []uint, scoresA, scoresB []int) (draw bool, winners []uint, losers []uint)
	CreateMatch(ctx context.Context, organizationID uint, teamA, teamB []uint, scoresA, scoresB []int) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) CreateMatch(ctx context.Context, organizationID uint, teamA, teamB []uint, scoresA, scoresB []int) error {
	marshalledTeamA, err := json.Marshal(teamA)
	if err != nil {
		return err
	}

	marshalledTeamB, err := json.Marshal(teamB)
	if err != nil {
		return err
	}

	sets := make([]string, len(scoresA))
	for i, scoreA := range scoresA {
		sets[i] = fmt.Sprintf("%d-%d", scoreA, scoresB[i])
	}

	match := &Match{
		OrganiziationID: organizationID,
		TeamA:           marshalledTeamA,
		TeamB:           marshalledTeamB,
		Sets:            sets,
	}

	if err := s.repo.CreateMatch(ctx, match); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) DetermineResult(ctx context.Context, teamA, teamB []uint, scoresA, scoresB []int) (bool, []uint, []uint) {
	teamASetWins := 0
	teamBSetWins := 0

	for i, scoreA := range scoresA {
		if scoreA > scoresB[i] {
			teamASetWins++
		} else if scoreA < scoresB[i] {
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
