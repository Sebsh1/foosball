package rest

import (
	"context"
	"fmt"
	"foosball/internal/authentication"
	"foosball/internal/match"
	"foosball/internal/player"
	"foosball/internal/rating"
	"foosball/internal/rest/controllers"
	"foosball/internal/rest/helpers"
	"foosball/internal/team"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port                 uint32 `mapstructure:"port" default:"8001"`
	GZIPCompressionLevel int    `mapstructure:"gzip_compression_level" default:"5"`
}

type Server struct {
	echo   *echo.Echo
	config Config
}

func NewServer(
	conf Config,
	logger *logrus.Logger,
	authService authentication.Service,
	playerService player.Service,
	matchService match.Service,
	ratingService rating.Service,
	teamService team.Service,
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
			Level:   conf.GZIPCompressionLevel,
			Skipper: middleware.DefaultGzipConfig.Skipper,
		}),
	)

	root := e.Group("/foosball")

	controllers.Register(
		root,
		logger.WithField("module", "rest"),
		authService,
		playerService,
		matchService,
		ratingService,
		teamService,
	)

	return &Server{
		echo:   e,
		config: conf,
	}, nil
}

func (s *Server) Start() error {
	return errors.Wrap(s.echo.Start(fmt.Sprintf(":%d", s.config.Port)), "Failed to start server")
}

func (s *Server) Shutdown(ctx context.Context) error {
	return errors.Wrap(s.echo.Shutdown(ctx), "Failed to shutdown server")
}
