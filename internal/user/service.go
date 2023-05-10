package user

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetUsers(ctx context.Context, ids []uint) ([]*User, error)
	GetUsersStats(ctx context.Context, ids []uint) ([]*UserStats, error)
	GetUserByEmail(ctx context.Context, email string) (exists bool, user *User, err error)
	GetUsersInOrganization(ctx context.Context, organizationID uint) ([]User, error)
	CreateUser(ctx context.Context, email, name, hash string) error
	DeleteUser(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, id uint, email, name, hash string, organizationID uint, admin bool) error
	UpdateRatings(ctx context.Context, userIDs []uint, ratings []int) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetUsers(ctx context.Context, ids []uint) ([]*User, error) {
	users, err := s.repo.GetUsers(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	return users, nil
}

func (s *ServiceImpl) GetUsersStats(ctx context.Context, ids []uint) ([]*UserStats, error) {
	users, err := s.repo.GetUsersStats(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users stats")
	}

	return users, nil
}

func (s *ServiceImpl) GetUsersInOrganization(ctx context.Context, organizationID uint) ([]User, error) {
	users, err := s.repo.GetUsersInOrganization(ctx, organizationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users in organization")
	}

	return users, nil
}

func (s *ServiceImpl) CreateUser(ctx context.Context, email, name, hash string) error {
	user := &User{
		Email: email,
		Name:  name,
		Hash:  hash,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (s *ServiceImpl) GetUserByEmail(ctx context.Context, email string) (bool, *User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil, nil
		}

		return false, nil, errors.Wrap(err, "failed to get user by email")
	}

	return true, user, nil
}

func (s *ServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.DeleteUserByID(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete user by id")
	}

	return nil
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, id uint, email, name, hash string, organizationID uint, admin bool) error {
	user := &User{
		ID:             id,
		Email:          email,
		Name:           name,
		Hash:           hash,
		OrganizationID: &organizationID,
		Admin:          admin,
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

func (s *ServiceImpl) UpdateRatings(ctx context.Context, userIDs []uint, ratings []int) error {
	if len(userIDs) != len(ratings) {
		return errors.New("userIDs and ratings must have the same length")
	}

	if err := s.repo.UpdateRatings(ctx, userIDs, ratings); err != nil {
		return errors.Wrap(err, "failed to update ratings")
	}

	return nil
}
