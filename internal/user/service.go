package user

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetUser(ctx context.Context, id uint) (*User, error)
	GetUsers(ctx context.Context, ids []uint) ([]*User, error)
	GetUserByEmail(ctx context.Context, email string) (exists bool, user *User, err error)
	CreateUser(ctx context.Context, email, name, hash string) error
	CreateVirtualUser(ctx context.Context, name string) error
	DeleteUser(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, id uint, email, name, hash string, virtual bool) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) GetUser(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user %d", id)
	}

	return user, nil
}

func (s *ServiceImpl) GetUsers(ctx context.Context, ids []uint) ([]*User, error) {
	users, err := s.repo.GetUsers(ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get users %v", ids)
	}

	return users, nil
}

func (s *ServiceImpl) CreateUser(ctx context.Context, email, name, hash string) error {
	user := &User{
		Email:   email,
		Name:    name,
		Hash:    hash,
		Virtual: false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (s *ServiceImpl) CreateVirtualUser(ctx context.Context, name string) error {
	user := &User{
		Name:    name,
		Virtual: true,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create virtual user")
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

func (s *ServiceImpl) UpdateUser(ctx context.Context, id uint, email, name, hash string, virtual bool) error {
	user := &User{
		ID:      id,
		Email:   email,
		Name:    name,
		Hash:    hash,
		Virtual: virtual,
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}
