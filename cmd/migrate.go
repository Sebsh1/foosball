package cmd

import (
	"context"
	"matchlog/migrations"
	"matchlog/pkg/database"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		logrus.WithError(err).Fatal("failed to load config")
	}

	l := GetLogger(config.Log)

	db, err := database.NewClient(ctx, config.DB.DSN)
	if err != nil {
		l.WithError(err).Fatal("failed to connect to database")
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.Migration00001Init,
	})

	if err = m.Migrate(); err != nil {
		logrus.WithError(err).Fatal("Migration failed")
	}

	logrus.Println("Migration finished successfully")
}
