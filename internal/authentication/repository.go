//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=authentication
package authentication

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetUser(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, username string) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetUser(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *RepositoryImpl) DeleteUser(ctx context.Context, username string) error {
	return r.db.WithContext(ctx).Delete(&User{}, "username = ?", username).Error
}
