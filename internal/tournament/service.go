//go:generate mockgen --source=service.go -destination=service_mock.go -package=tournament
package tournament

import (
	"context"
	"foosball/internal/team"
)

type Service interface {
	GetTournament(ctx context.Context, id uint) (*Tournament, error)
	CreateTournament(ctx context.Context, teams []*team.Team) error
	UpdateTournament(ctx context.Context, tournament *Tournament) error
	DeleteTournament(ctx context.Context, id uint) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (*ServiceImpl) CreateTournament(ctx context.Context, teams []*team.Team) error {
	// TODO
	panic("unimplemented")
}

func (*ServiceImpl) DeleteTournament(ctx context.Context, id uint) error {
	// TODO
	panic("unimplemented")
}

func (*ServiceImpl) GetTournament(ctx context.Context, id uint) (*Tournament, error) {
	// TODO
	panic("unimplemented")
}

func (*ServiceImpl) UpdateTournament(ctx context.Context, tournament *Tournament) error {
	// TODO
	panic("unimplemented")
}
