//go:generate mockgen --source=service.go -destination=service_mock.go -package=season
package season

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	Length time.Duration `mapstructure:"length" default:"90d"`
}

type Service interface {
	GetSeason(ctx context.Context, id uint) (*Season, error)
	CreateSeason(ctx context.Context, start time.Time, name *string) error
	UpdateSeason(ctx context.Context, season *Season) error
	DeleteSeason(ctx context.Context, id uint) error
}

type ServiceImpl struct {
	config Config
	repo   Repository
}

func NewService(config Config, repo Repository) Service {
	return &ServiceImpl{
		config: config,
		repo:   repo,
	}
}

func (s *ServiceImpl) GetSeason(ctx context.Context, id uint) (*Season, error) {
	player, err := s.repo.GetSeason(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get season")
	}

	return player, nil
}

func (s *ServiceImpl) CreateSeason(ctx context.Context, start time.Time, name *string) error {
	season := &Season{
		Start: start,
		End:   start.Add(s.config.Length),
	}

	if name != nil {
		season.Name = *name
	}

	err := s.repo.CreateSeason(ctx, season)
	if err != nil {
		return errors.Wrap(err, "failed to create season")
	}

	return nil
}

func (s *ServiceImpl) DeleteSeason(ctx context.Context, id uint) error {
	season, err := s.repo.GetSeason(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get season")
	}

	err = s.repo.DeleteSeason(ctx, season)
	if err != nil {
		return errors.Wrap(err, "failed to delete season")
	}

	return nil
}

func (s *ServiceImpl) UpdateSeason(ctx context.Context, season *Season) error {
	err := s.repo.UpdateSeason(ctx, season)
	if err != nil {
		return errors.Wrap(err, "failed to update season")
	}

	return nil
}
