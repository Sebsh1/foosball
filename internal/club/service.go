package club

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetClub(ctx context.Context, id uint) (*Club, error)
	GetClubs(ctx context.Context, ids []uint) ([]Club, error)
	GetUserIDsInClub(ctx context.Context, id uint) ([]uint, error)
	GetInvitesByUserID(ctx context.Context, userID uint) ([]ClubsUsers, error)
	CreateClub(ctx context.Context, name string, adminUserID uint) (orgID uint, err error)
	RemoveUserFromClub(ctx context.Context, userID uint, ClubID uint) error
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
	UpdateUserRole(ctx context.Context, userID uint, ClubID uint, role Role) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetClub(ctx context.Context, id uint) (*Club, error) {
	org, err := s.repo.GetClub(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Club")
	}

	return org, nil
}

func (s *service) GetClubs(ctx context.Context, ids []uint) ([]Club, error) {
	orgs, err := s.repo.GetClubs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Clubs")
	}

	return orgs, nil
}

func (s *service) GetUserIDsInClub(ctx context.Context, id uint) ([]uint, error) {
	userIDs, err := s.repo.GetUserIDsInClub(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get userIDs in Club")
	}

	return userIDs, nil
}

func (s *service) GetInvitesByUserID(ctx context.Context, userID uint) ([]ClubsUsers, error) {
	invites, err := s.repo.GetInvitesByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user invites")
	}

	return invites, nil
}

func (s *service) CreateClub(ctx context.Context, name string, adminUserID uint) (uint, error) {
	org := &Club{
		Name: name,
	}

	orgID, err := s.repo.CreateClub(ctx, org)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create Club")
	}

	if err := s.repo.AddUserToClub(ctx, adminUserID, orgID, AdminRole); err != nil {
		return 0, errors.Wrap(err, "failed to add creating user to Club")
	}

	return orgID, nil
}

func (s *service) RemoveUserFromClub(ctx context.Context, userID uint, ClubID uint) error {
	if err := s.repo.RemoveUserFromClub(ctx, userID, ClubID); err != nil {
		return errors.Wrap(err, "failed to remove user from Club")
	}

	return nil
}

func (s *service) DeleteClub(ctx context.Context, id uint) error {
	if err := s.repo.DeleteClub(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete Club")
	}

	return nil
}

func (s *service) UpdateClub(ctx context.Context, id uint, name string) error {
	if err := s.repo.UpdateClub(ctx, id, name); err != nil {
		return errors.Wrap(err, "failed to update Club")
	}

	return nil
}

func (s *service) UpdateUserRole(ctx context.Context, userID uint, ClubID uint, role Role) error {
	if err := s.repo.UpdateUserRole(ctx, userID, ClubID, role); err != nil {
		return errors.Wrap(err, "failed to update user role")
	}

	return nil
}
