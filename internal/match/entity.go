package match

import (
	"foosball/internal/team"
	"time"
)

type Match struct {
	ID uint `gorm:"primaryKey"`

	TeamA  team.Team `gorm:"index;not null"`
	TeamB  team.Team `gorm:"index;not null"`
	GoalsA int       `gorm:"index;not null"`
	GoalsB int       `gorm:"index;not null"`

	CreatedAt time.Time
}
