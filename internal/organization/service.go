package organization

import (
	"context"
)

type Service interface {
	GetOrganization(ctx context.Context, id uint) (*Organization, error)
	CreateOrganization(ctx context.Context, name, ratingMethod string) error
	DeleteOrganization(ctx context.Context, id uint) error
	UpdateOrganization(ctx context.Context, id uint, name, ratingMethod string) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetOrganization(ctx context.Context, id uint) (*Organization, error) {
	org, err := s.repo.GetOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *ServiceImpl) CreateOrganization(ctx context.Context, name, ratingMethod string) error {
	org := &Organization{
		Name:         name,
		RatingMethod: ratingMethod,
	}

	if err := s.repo.CreateOrganization(ctx, org); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) DeleteOrganization(ctx context.Context, id uint) error {
	if err := s.repo.DeleteOrganization(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) UpdateOrganization(ctx context.Context, id uint, name, ratingMethod string) error {
	if err := s.repo.UpdateOrganization(ctx, id, name, ratingMethod); err != nil {
		return err
	}

	return nil
}
