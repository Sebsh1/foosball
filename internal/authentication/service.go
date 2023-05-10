//go:generate mockgen --source=service.go -destination=service_mock.go -package=authentication
package authentication

import (
	"context"
	"foosball/internal/user"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, email, password string) (valid bool, token string, err error)
	VerifyJWT(ctx context.Context, token string) (valid bool, claims *Claims, err error)
	Signup(ctx context.Context, email, username, password string) (success bool, err error)
}

type ServiceImpl struct {
	secret      string
	userService user.Service
}

func NewService(secret string, userService user.Service) Service {
	return &ServiceImpl{
		secret:      secret,
		userService: userService,
	}
}

func (s *ServiceImpl) Login(ctx context.Context, email string, password string) (bool, string, error) {
	exists, user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, "", err
	}

	if !exists {
		logrus.WithField("email", email).Debug("user does not exist")
		return false, "", nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return false, "", nil
	}

	token, err := s.generateJWT(user.Name, user.ID, user.OrganizationID, user.Admin)
	if err != nil {
		return false, "", err
	}

	return true, token, nil
}

func (s *ServiceImpl) Signup(ctx context.Context, email string, username string, password string) (bool, error) {
	exists, _, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return false, err
	}

	if err = s.userService.CreateUser(ctx, email, username, string(hashedPasswordBytes)); err != nil {
		return false, err
	}

	return true, nil
}

func (s *ServiceImpl) VerifyJWT(ctx context.Context, token string) (bool, *Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return false, nil, err
	}

	if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
		return false, nil, errors.New("unexpected signing method")
	}

	if !parsedToken.Valid {
		return false, nil, errors.New("jwt is invalid")
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return false, nil, errors.New("failed to parse claims")
	}

	return true, claims, nil
}

func (s *ServiceImpl) generateJWT(name string, userID uint, organizationID *uint, admin bool) (string, error) {
	now := time.Now()

	standardClaims := jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(6 * time.Hour).Unix(),
		Issuer:    "matchlogger",
	}

	orgID := uint(0)
	if organizationID != nil {
		orgID = *organizationID
	}

	claims := Claims{
		StandardClaims: standardClaims,
		Name:           name,
		UserID:         userID,
		OrganizationID: orgID,
		Admin:          admin,
	}

	tokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := tokenUnsigned.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to sign access token")
	}

	return tokenSigned, nil
}
