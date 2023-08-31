package cmd

import (
	"context"
	"matchlog/internal/authentication"
	"matchlog/internal/invite"
	"matchlog/internal/leaderboard"
	"matchlog/internal/match"
	"matchlog/internal/organization"
	"matchlog/internal/rating"
	"matchlog/internal/rest"
	"matchlog/internal/statistic"
	"matchlog/internal/user"
	"matchlog/pkg/database"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:  "serve",
	Long: "Start the service",
	Run:  serve,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config, err := loadConfig()
	if err != nil {
		zap.L().Fatal("Failed to load config", zap.Error(err))
	}

	l := GetLogger(config.Log)

	// Initialize database connection
	db, err := database.NewClient(ctx, config.DB.DSN)
	if err != nil {
		l.Fatal("Failed to connect to database",
			"error", err)
	}

	// Initialize User service
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// Initialize Organization service
	organizationRepository := organization.NewRepository(db)
	organizationService := organization.NewService(organizationRepository)

	// Initialize Invite service
	inviteRepository := invite.NewRepository(db)
	inviteService := invite.NewService(inviteRepository, userService, organizationService)

	// Initialize Authentication service
	authenticationService := authentication.NewService(config.Authentication.Secret, userService)

	// Initialize Match service
	matchRepository := match.NewRepository(db)
	matchService := match.NewService(matchRepository)

	// Initialize Rating service
	ratingRepository := rating.NewRepository(db)
	ratingService := rating.NewService(ratingRepository)

	// Initialize Statistics service
	statisticRepository := statistic.NewRepository(db)
	statisticService := statistic.NewService(statisticRepository)

	// Initialize Leaderboard service
	leaderboardService := leaderboard.NewService(userService, ratingService, statisticService)

	// Initialize REST server
	restServer, err := rest.NewServer(
		config.Rest.Port,
		l,
		authenticationService,
		userService,
		organizationService,
		inviteService,
		matchService,
		ratingService,
		statisticService,
		leaderboardService,
	)
	if err != nil {
		l.Fatal("Failed to create rest server",
			"error", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the REST server
	l.Info("REST server starting", zap.Int("port", config.Rest.Port))
	go func() {
		if err := restServer.Start(); err != nil {
			l.Fatal("Failed to start rest server",
				"error", err)
			cancel()
		}
	}()

	l.Info("Ready")
	// Wait for shutdown signal
	<-ctx.Done()

	// Stop the servers
	l.Info("Shutting down")

	shutdownctx, stop := context.WithTimeout(context.Background(), shutdownPeriod)
	defer stop()

	if err := restServer.Shutdown(shutdownctx); err != nil {
		l.Error("Failed to shutdown rest server",
			"error", err)
	}
}
