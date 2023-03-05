package season

import (
	"foosball/internal/match"
	"time"
)

type Season struct {
	ID uint `gorm:"primaryKey"`

	Name    string         `gorm:"index"`
	Matches []*match.Match `gorm:"index"`

	Start time.Time `gorm:"index;not null"`
	End   time.Time `gorm:"index;not null"`

	CreatedAt time.Time
}
