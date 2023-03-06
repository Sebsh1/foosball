package tournament

import (
	"context"
	"foosball/internal/models"
)

type Service interface {
	GetTournament(ctx context.Context, id uint) (*models.Tournament, error)
	CreateTournament(ctx context.Context, teams []*models.Team) error
	DeleteTournament(ctx context.Context, id uint) error
	UpdateTournament(ctx context.Context, tournament *models.Tournament) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (*ServiceImpl) CreateTournament(ctx context.Context, teams []*models.Team) error {
	// TODO
	panic("unimplemented")
}

func (*ServiceImpl) DeleteTournament(ctx context.Context, id uint) error {
	// TODO
	panic("unimplemented")
}

func (*ServiceImpl) GetTournament(ctx context.Context, id uint) (*models.Tournament, error) {
	// TODO
	panic("unimplemented")
}

func (*ServiceImpl) UpdateTournament(ctx context.Context, tournament *models.Tournament) error {
	// TODO
	panic("unimplemented")
}
