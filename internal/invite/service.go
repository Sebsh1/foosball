package invite

import (
	"context"
	"matchlog/internal/organization"
	"matchlog/internal/user"
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
		return err
	}

	return nil
}

func (s *ServiceImpl) GetInvitesByUserID(ctx context.Context, userID uint) ([]Invite, error) {
	invites, err := s.repo.GetInvitesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return invites, nil
}

func (s *ServiceImpl) GetInvitesByOrganizationID(ctx context.Context, organizationID uint) ([]Invite, error) {
	invites, err := s.repo.GetInvitesByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	return invites, nil
}

func (s *ServiceImpl) DeclineInvite(ctx context.Context, id uint) error {
	if err := s.repo.DeleteInvite(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) AcceptInvite(ctx context.Context, id uint) error {
	invite, err := s.repo.GetInvite(ctx, id)
	if err != nil {
		return err
	}

	user, err := s.userService.GetUser(ctx, invite.UserID)
	if err != nil {
		return err
	}

	if err := s.userService.UpdateUser(ctx, user.ID, user.Email, user.Name, user.Hash, &invite.OrganizationID, user.Role); err != nil {
		return err
	}

	if err := s.repo.DeleteInvite(ctx, id); err != nil {
		return err
	}

	return nil
}
