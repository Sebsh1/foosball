package cmd

import (
	"context"
	"matchlog/migrations"
	"matchlog/pkg/database"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate the database",
	Run:  migrate,
}

func init() { //nolint:gochecknoinits
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	config, err := loadConfig()
	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}

	l := GetLogger(config.Log)

	db, err := database.NewClient(ctx, config.DB.DSN)
	if err != nil {
		l.Fatal("failed to connect to database",
			"error", err)
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.Migration00001Init,
	})

	if err = m.Migrate(); err != nil {
		l.Fatal("Migration failed",
			"error", err)
	}

	zap.L().Info("Migration finished successfully")
}
