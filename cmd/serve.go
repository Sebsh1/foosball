package cmd

import (
	"context"
	"foosball/internal/match"
	"foosball/internal/mysql"
	"foosball/internal/player"
	"foosball/internal/rating"
	"foosball/internal/rest"
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
		log.WithError(err).Fatal("failed to connect to Monitor database")
	}

	if err := db.AutoMigrate(
		&player.Player{},
		&match.Match{},
	); err != nil {
		log.WithError(err).Fatal("failed to auto migrate database")
	}

	playerRepo := player.NewRepository(db)
	playerService := player.NewService(playerRepo)

	matchRepo := match.NewRepository(db)
	matchService := match.NewService(matchRepo, playerService)

	ratingService := rating.NewService(config.Rating, playerService)

	httpServer, err := rest.NewServer(
		config.Rest,
		log,
		playerService,
		matchService,
		ratingService,
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
		log.WithError(err).Error("Failed to shutdown internal http server")
	}

}
