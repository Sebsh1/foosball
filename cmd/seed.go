package cmd

import (
	"context"
	"foosball/internal/authentication"
	"foosball/internal/match"
	"foosball/internal/mysql"
	"foosball/internal/player"
	"foosball/internal/season"
	"foosball/internal/team"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:  "seed",
	Long: "Seeds the database with some initial data. Should only be used for development.",
	Run:  seed,
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func seed(cmd *cobra.Command, args []string) {
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

	teamRepo := team.NewRepository(db)
	teamService := team.NewService(teamRepo)

	seasonRepo := season.NewRepository(db)
	seasonService := season.NewService(config.Season, seasonRepo)

	matchRepo := match.NewRepository(db)
	matchService := match.NewService(matchRepo, playerService, seasonService)

	log.Info("Seeding database...")

	if err := seedUsers(ctx, authService); err != nil {
		log.WithError(err).Warn("failed to seed users")
	}

	if err := seedPlayers(ctx, playerService); err != nil {
		log.WithError(err).Warn("failed to seed players")
	}

	if err := seedTeams(ctx, teamService, playerService); err != nil {
		log.WithError(err).Warn("failed to seed teams")
	}

	if err := seedSeasons(ctx, seasonService); err != nil {
		log.WithError(err).Warn("failed to seed seasons")
	}

	if err := seedMatches(ctx, matchService); err != nil {
		log.WithError(err).Warn("failed to seed matches")
	}

	log.Info("Finished seeding database.")
}

func seedUsers(ctx context.Context, authService authentication.Service) error {
	usernames := []string{"admin", "user1", "user2", "user3", "user4", "user5", "user6", "user7", "user8", "user9", "user10"}
	passwords := []string{"admin", "password1", "password2", "password3", "password4", "password5", "password6", "password7", "password8", "password9", "password10"}
	for i, username := range usernames {
		if err := authService.CreateUser(ctx, username, passwords[i]); err != nil {
			return err
		}
	}

	return nil
}

func seedPlayers(ctx context.Context, playerService player.Service) error {
	names := []string{"admin", "Hans", "Peter", "Max", "Moritz", "Fritz", "Karl", "Mikkel", "Mads", "Morten", "Sebastian"}
	for _, name := range names {
		if err := playerService.CreatePlayer(ctx, name); err != nil {
			return err
		}
	}
	return nil
}

func seedTeams(ctx context.Context, teamService team.Service, playerService player.Service) error {
	for i := 1; i <= 5; i++ {
		ids := []uint{uint(i), uint(i + 1)}

		players, err := playerService.GetPlayers(ctx, ids)
		if err != nil {
			return err
		}

		if err := teamService.CreateTeam(ctx, players); err != nil {
			return err
		}
	}

	for i := 4; i <= 8; i++ {
		ids := []uint{uint(i), uint(i + 1), uint(i + 2)}

		players, err := playerService.GetPlayers(ctx, ids)
		if err != nil {
			return err
		}

		if err := teamService.CreateTeam(ctx, players); err != nil {
			return err
		}
	}

	return nil
}

func seedSeasons(ctx context.Context, seasonService season.Service) error {
	names := []string{"Alpha Season", "Beta Season", "Season 1"}
	starts := []string{"2021-01-01", "2022-01-01", "2023-01-01"}
	for i, name := range names {
		start, err := time.Parse("2006-01-02", starts[i])
		if err != nil {
			return err
		}
		if err := seasonService.CreateSeason(ctx, &name, start); err != nil {
			return err
		}
	}
	return nil
}

func seedMatches(ctx context.Context, matchService match.Service) error {
	teamAids := []uint{1, 2, 3, 4}
	teamBids := []uint{5, 6, 7, 8}
	goalsA := []int{10, 8, 5, 10}
	goalsB := []int{9, 10, 10, 4}

	for i := 0; i < 4; i++ {
		if err := matchService.CreateMatch(ctx, teamAids[i], teamBids[i], goalsA[i], goalsB[i]); err != nil {
			return err
		}
	}
	return nil
}
