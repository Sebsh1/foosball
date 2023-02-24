package migrations

import (
	"foosball/internal/player"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migration00001Init contains the initial migration.
var Migration00001Init = &gormigrate.Migration{
	ID: "init_00001",
	Migrate: func(tx *gorm.DB) error {

		type Player struct {
			ID uint

			Name   string `gorm:"index;not null" default:"Unknown Player"`
			Rating int    `gorm:"index;not null" default:"1000"`

			CreatedAt time.Time
		}

		type Match struct {
			ID uint

			TeamA  []player.Player `gorm:"index;not null"`
			TeamB  []player.Player `gorm:"index;not null"`
			ScoreA int             `gorm:"index;not null"`
			ScoreB int             `gorm:"index;not null"`
			Winner string          `gorm:"index;not null"`

			CreatedAt time.Time
		}

		return tx.AutoMigrate(
			&Player{},
			&Match{},
		)
	},
}
