package season

import (
	"foosball/internal/match"
	"time"
)

type Season struct {
	ID uint `gorm:"primaryKey"`

	Name    string         `gorm:"index"`
	Matches []*match.Match `gorm:"foreignKey:SeasonID;constraint:OnDelete:CASCADE;"`

	Start time.Time `gorm:"not null"`
	End   time.Time `gorm:"not null"`

	CreatedAt time.Time
}
