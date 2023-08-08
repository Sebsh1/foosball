package organization

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetOrganization(ctx context.Context, id uint) (*Organization, error)
	CreateOrganization(ctx context.Context, name string) error
	DeleteOrganization(ctx context.Context, id uint) error
	UpdateOrganization(ctx context.Context, id uint, name string) error
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
		return nil, errors.Wrap(err, "failed to get organization")
	}

	return org, nil
}

func (s *ServiceImpl) CreateOrganization(ctx context.Context, name string) error {
	org := &Organization{
		Name: name,
	}

	if err := s.repo.CreateOrganization(ctx, org); err != nil {
		return errors.Wrap(err, "failed to create organization")
	}

	return nil
}

func (s *ServiceImpl) DeleteOrganization(ctx context.Context, id uint) error {
	if err := s.repo.DeleteOrganization(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete organization")
	}

	return nil
}

func (s *ServiceImpl) UpdateOrganization(ctx context.Context, id uint, name string) error {
	if err := s.repo.UpdateOrganization(ctx, id, name); err != nil {
		return errors.Wrap(err, "failed to update organization")
	}

	return nil
}
