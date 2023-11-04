//go:generate mockgen --source=service.go -destination=service_mock.go -package=authentication
package authentication

import (
	"context"
	"matchlog/internal/user"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Secret string `mapstructure:"secret"`
}

type Service interface {
	Login(ctx context.Context, email, password string) (valid bool, accessToken, refreshToken string, err error)
	VerifyAccessToken(ctx context.Context, token string) (valid bool, claims *AccessClaims, err error)
	RefreshTokenPair(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)
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

func (s *ServiceImpl) Login(ctx context.Context, email string, password string) (bool, string, string, error) {
	exists, user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, "", "", errors.Wrap(err, "failed to get user by email")
	}

	if !exists {
		return false, "", "", nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return false, "", "", nil
	}

	accessToken, refreshToken, err := s.generateTokenPair(user.Name, user.Id)
	if err != nil {
		return false, "", "", errors.Wrap(err, "failed to generate jwt")
	}

	return true, accessToken, refreshToken, nil
}

func (s *ServiceImpl) Signup(ctx context.Context, email string, username string, password string) (bool, error) {
	exists, _, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user by email")
	}

	if exists {
		return false, nil
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return false, errors.Wrap(err, "failed to hash password")
	}

	if err = s.userService.CreateUser(ctx, email, username, string(hashedPasswordBytes)); err != nil {
		return false, errors.Wrap(err, "failed to create user")
	}

	return true, nil
}

func (s *ServiceImpl) VerifyAccessToken(ctx context.Context, token string) (bool, *AccessClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := parsedToken.Claims.(*AccessClaims)
	if !ok {
		return false, nil, errors.New("failed to parse claims")
	}

	return true, claims, nil
}

func (s *ServiceImpl) verifyRefreshToken(ctx context.Context, token string) (bool, *RefreshClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := parsedToken.Claims.(*RefreshClaims)
	if !ok {
		return false, nil, errors.New("failed to parse claims")
	}

	return true, claims, nil
}

func (s *ServiceImpl) RefreshTokenPair(ctx context.Context, refreshToken string) (string, string, error) {
	valid, claims, err := s.verifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to verify refresh token")
	}

	if !valid {
		return "", "", errors.New("refresh token is invalid")
	}

	u, err := s.userService.GetUser(ctx, claims.UserId)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get user by id")
	}

	accessToken, newRefreshToken, err := s.generateTokenPair(u.Name, claims.UserId)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to generate jwts")
	}

	return accessToken, newRefreshToken, nil
}

func (s *ServiceImpl) generateTokenPair(name string, userId uint) (string, string, error) {
	now := time.Now()

	accessStandardClaims := jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(15 * time.Minute).Unix(),
		Issuer:    "matchlog",
	}

	accessclaims := AccessClaims{
		StandardClaims: accessStandardClaims,
		Name:           name,
		UserId:         userId,
	}

	accessTokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, accessclaims)
	accessTokenSigned, err := accessTokenUnsigned.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to sign access token")
	}

	refreshStandardClaims := jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(15 * time.Minute).Unix(),
		Issuer:    "matchlog",
	}

	refreshTokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshStandardClaims)
	refreshTokenSigned, err := refreshTokenUnsigned.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to sign refresh token")
	}

	return accessTokenSigned, refreshTokenSigned, nil
}
