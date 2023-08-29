package invite

import (
	"context"
	"matchlog/internal/organization"
	"matchlog/internal/user"

	"github.com/pkg/errors"
)

type Service interface {
	CreateInvites(ctx context.Context, userID []uint, organizationID uint) error
	GetInvitesByUserID(ctx context.Context, userID uint) ([]Invite, error)
	GetInvitesByOrganizationID(ctx context.Context, organizationID uint) ([]Invite, error)
	DeclineInvite(ctx context.Context, id uint) error
	AcceptInvite(ctx context.Context, id uint) error
}

type ServiceImpl struct {
	repo        Repository
	userService user.Service
	orgService  organization.Service
}

func NewService(repo Repository, userService user.Service, orgService organization.Service) Service {
	return &ServiceImpl{
		repo:        repo,
		userService: userService,
		orgService:  orgService,
	}
}

func (s *ServiceImpl) CreateInvites(ctx context.Context, userIDs []uint, organizationID uint) error {
	if err := s.repo.CreateInvites(ctx, userIDs, organizationID); err != nil {
		return errors.Wrap(err, "failed to create invites")
	}

	return nil
}

func (s *ServiceImpl) GetInvitesByUserID(ctx context.Context, userID uint) ([]Invite, error) {
	invites, err := s.repo.GetInvitesByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invites by user id")
	}

	return invites, nil
}

func (s *ServiceImpl) GetInvitesByOrganizationID(ctx context.Context, organizationID uint) ([]Invite, error) {
	invites, err := s.repo.GetInvitesByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invites by organization id")
	}

	return invites, nil
}

func (s *ServiceImpl) DeclineInvite(ctx context.Context, id uint) error {
	if err := s.repo.DeleteInvite(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete invite")
	}

	return nil
}

func (s *ServiceImpl) AcceptInvite(ctx context.Context, id uint) error {
	invite, err := s.repo.GetInvite(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get invite")
	}

	u, err := s.userService.GetUser(ctx, invite.UserID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	if err := s.userService.UpdateUser(ctx, u.ID, u.Email, u.Name, u.Hash, &invite.OrganizationID, user.MemberRole); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	if err := s.repo.DeleteInvite(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete invite")
	}

	return nil
}
