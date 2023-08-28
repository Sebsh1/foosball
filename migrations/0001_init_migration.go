package migrations

import (
	"matchlog/internal/match"
	"matchlog/internal/user"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migration00001Init contains the initial migration.
var Migration00001Init = &gormigrate.Migration{
	ID: "init_00001",
	Migrate: func(tx *gorm.DB) error {
		type User struct {
			ID uint `gorm:"primaryKey"`

			Email string `gorm:"index"`
			Name  string `gorm:"index;not null"`
			Hash  string

			OrganizationID *uint     `gorm:"index"`
			Role           user.Role `default:"none"`

			CreatedAt time.Time
		}

		type Statistic struct {
			ID uint `gorm:"primaryKey"`

			UserID uint `gorm:"not null"`

			Wins   int
			Draws  int
			Losses int
			Streak int

			CreatedAt time.Time
		}

		type Rating struct {
			ID uint `gorm:"primaryKey"`

			UserID         uint `gorm:"not null"`
			OrganizationID uint `gorm:"not null"`

			Value      float64 `gorm:"default:1000.0"`
			Deviation  float64
			Volatility float64 `gorm:"default:0.06"`

			CreatedAt time.Time
		}

		type Organization struct {
			ID uint `gorm:"primaryKey"`

			Name string `gorm:"not null"`

			CreatedAt time.Time
		}

		type Match struct {
			ID uint `gorm:"primaryKey"`

			TeamA  []uint       `gorm:"serializer:json;not null"`
			TeamB  []uint       `gorm:"serializer:json;not null"`
			Sets   []string     `gorm:"serializer:json;not null"`
			Result match.Result `gorm:"not null"`

			CreatedAt time.Time
		}

		type Invite struct {
			ID uint `gorm:"primaryKey"`

			OrganizationID uint `gorm:"not null"`
			UserID         uint `gorm:"not null"`

			CreatedAt time.Time
		}

		return tx.AutoMigrate(
			&User{},
			&Statistic{},
			&Rating{},
			&Organization{},
			&Match{},
			&Invite{},
		)
	},
}
