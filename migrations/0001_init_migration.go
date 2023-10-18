package migrations

import (
	"matchlog/internal/club"
	"matchlog/internal/match"
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

			Email   string `gorm:"index"`
			Name    string `gorm:"index;not null"`
			Hash    string
			Virtual bool `gorm:"default:false"`

			CreatedAt time.Time
		}

		type Statistic struct {
			ID uint `gorm:"primaryKey"`

			UserID uint `gorm:"not null"`
			GameID uint `gorm:"not null"`

			Wins   int
			Draws  int
			Losses int
			Streak int

			CreatedAt time.Time
		}

		type Rating struct {
			ID uint `gorm:"primaryKey"`

			UserID uint `gorm:"not null"`
			GameID uint `gorm:"not null"`

			Value      float64 `gorm:"default:1000.0"`
			Deviation  float64
			Volatility float64 `gorm:"default:0.06"`

			CreatedAt time.Time
		}

		type Game struct {
			ID   uint   `gorm:"primaryKey"`
			Name string `gorm:"not null"`
		}

		type Club struct {
			ID uint `gorm:"primaryKey"`

			Name string `gorm:"not null"`

			CreatedAt time.Time
		}

		type ClubsUsers struct {
			ID uint `gorm:"primaryKey"`

			ClubID uint `gorm:"primaryKey"`
			UserID uint `gorm:"primaryKey"`

			Accepted bool      `gorm:"default:false"`
			Role     club.Role `gorm:"default:member"`

			CreatedAt time.Time
		}

		type ClubsGames struct {
			ID uint `gorm:"primaryKey"`

			ClubID uint `gorm:"primaryKey"`
			GameID uint `gorm:"primaryKey"`

			CreatedAt time.Time
		}

		type ClubsTournaments struct {
			ID uint `gorm:"primaryKey"`

			ClubID       uint `gorm:"primaryKey"`
			TournamentID uint `gorm:"primaryKey"`

			CreatedAt time.Time
		}

		type ClubsLeagues struct {
			ID uint `gorm:"primaryKey"`

			ClubID   uint `gorm:"primaryKey"`
			LeagueID uint `gorm:"primaryKey"`

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

		return tx.AutoMigrate(
			&User{},
			&Statistic{},
			&Rating{},
			&Club{},
			&Match{},
			&Game{},
			&ClubsUsers{},
			&ClubsGames{},
			&ClubsTournaments{},
			&ClubsLeagues{},
		)
	},
}
