//go:generate mockgen --source=service.go -destination=service_mock.go -package=authentication
package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Secret string
}

type Service interface {
	Login(ctx context.Context, username, password string) (valid bool, token string, err error)
	CreateUser(ctx context.Context, username, password string) error
	DeleteUser(ctx context.Context, username string) error
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

func (s *ServiceImpl) Login(ctx context.Context, username, password string) (bool, string, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return false, "", errors.Wrap(err, "failed to get user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return false, "", nil
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return false, "", errors.Wrap(err, "failed to generate token")
	}

	return true, token, nil
}

func (s *ServiceImpl) generateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	return token.SignedString([]byte(s.config.Secret))
}

func (s *ServiceImpl) CreateUser(ctx context.Context, username, password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	user := &User{
		Username:     username,
		PasswordHash: string(hashBytes),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (s *ServiceImpl) DeleteUser(ctx context.Context, username string) error {
	if err := s.repo.DeleteUser(ctx, username); err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}
