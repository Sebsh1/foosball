package rest

import (
	"context"
	"fmt"
	"matchlog/internal/authentication"
	"matchlog/internal/club"
	"matchlog/internal/leaderboard"
	"matchlog/internal/match"
	"matchlog/internal/rating"
	"matchlog/internal/rest/controllers"
	"matchlog/internal/rest/helpers"
	"matchlog/internal/statistic"
	"matchlog/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	Port int `mapstructure:"port" default:"8001"`
}

type Server struct {
	echo *echo.Echo
	port int
}

func NewServer(
	port int,
	logger *zap.SugaredLogger,
	authService authentication.Service,
	userService user.Service,
	ClubService club.Service,
	matchService match.Service,
	ratingService rating.Service,
	statisticService statistic.Service,
	leaderboardService leaderboard.Service,
) (*Server, error) {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Validator = helpers.NewValidator()

	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: middleware.DefaultGzipConfig.Skipper,
		}),
	)

	root := e.Group("/api")

	controllers.Register(
		root,
		logger.With("module", "rest"),
		authService,
		userService,
		ClubService,
		matchService,
		ratingService,
		statisticService,
		leaderboardService,
	)

	return &Server{
		echo: e,
		port: port,
	}, nil
}

func (s *Server) Start() error {
	err := s.echo.Start(fmt.Sprintf("0.0.0.0:%d", s.port))
	if err != nil {
		return errors.Wrap(err, "Failed to start server")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to shutdown server")
	}

	return nil
}
