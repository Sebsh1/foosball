package cmd

import (
	"context"
	"foosball/internal/authentication"
	"foosball/internal/match"
	"foosball/internal/mysql"
	"foosball/internal/player"
	"foosball/internal/rating"
	"foosball/internal/rest"
	"foosball/internal/season"
	"foosball/internal/team"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:  "serve",
	Long: "Start the service",
	Run:  serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	log := GetLogger(config.Log)

	db, err := mysql.NewClient(ctx, config.DB)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}

	if err := db.AutoMigrate(
		&authentication.User{},
		&player.Player{},
		&team.Team{},
		&season.Season{},
		&match.Match{},
	); err != nil {
		log.WithError(err).Fatal("failed to auto migrate database")
	}

	authRepo := authentication.NewRepository(db)
	authService := authentication.NewService(config.Auth, authRepo)

	playerRepo := player.NewRepository(db)
	playerService := player.NewService(playerRepo)

	ratingService := rating.NewService(config.Rating, playerService)

	teamRepo := team.NewRepository(db)
	teamService := team.NewService(teamRepo)

	seasonRepo := season.NewRepository(db)
	seasonService := season.NewService(config.Season, seasonRepo)

	matchRepo := match.NewRepository(db)
	matchService := match.NewService(matchRepo, playerService, seasonService)

	httpServer, err := rest.NewServer(
		config.Rest,
		log,
		authService,
		playerService,
		matchService,
		ratingService,
		teamService,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to create http server")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.WithField("port", config.Rest.Port).Info("REST server starting")
	go func() {
		if err := httpServer.Start(); err != nil {
			log.WithError(err).Error("failed to start http server")
			cancel()
		}
	}()

	log.Info("Ready")

	<-ctx.Done()
	log.Info("shutting down")
	shutdownCtx, stop := context.WithTimeout(context.Background(), shutdownPeriod)
	defer stop()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("failed to shutdown internal http server")
	}

}
