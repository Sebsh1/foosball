package tournament

import (
	"foosball/internal/match"
	"foosball/internal/team"
	"time"
)

type Tournament struct {
	ID uint `gorm:"primaryKey"`

	Teams   []team.Team   `gorm:"not null"`
	Matches []match.Match `gorm:"not null"`

	CreatedAt time.Time
}
