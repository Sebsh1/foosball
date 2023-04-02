package tournament

import (
	"foosball/internal/match"
	"time"
)

type Tournament struct {
	ID       uint `gorm:"primaryKey"`
	SeasonID uint `gorm:"index"`

	Name    string
	Matches *[]match.Match `gorm:"foreignKey:TournamentID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time
}
