package cmd

import (
	"context"
	"foosball/internal/authentication"
	"foosball/internal/connectors/mysql"
	"foosball/internal/invite"
	"foosball/internal/match"
	"foosball/internal/organization"
	"foosball/internal/user"

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
		log.WithError(err).Fatal("failed to connect to database")
	}

	if err := db.AutoMigrate(
		&match.Match{},
		&match.Set{},
		&organization.Organization{},
		&user.User{},
		&invite.Invite{},
	); err != nil {
		log.WithError(err).Fatal("failed to auto migrate database")
	}

	userReposotory := user.NewRepository(db)
	userService := user.NewService(userReposotory)

	ratingService := rating.NewService(config.Rating, userService)

	organizationRepository := organization.NewRepository(db)
	organizationService := organization.NewService(organizationRepository)

	authService := authentication.NewService(config.Auth, userService)

	inviteRepo := invite.NewRepository(db)
	inviteService := invite.NewService(inviteRepo, userService, organizationService)

	matchRepo := match.NewRepository(db)
	matchService := match.NewService(matchRepo)

	httpServer, err := rest.NewServer(
		config.Rest,
		log,
		authService,
		userService,
		organizationService,
		inviteService,
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
		log.WithError(err).Error("failed to shutdown internal http server")
	}

}
