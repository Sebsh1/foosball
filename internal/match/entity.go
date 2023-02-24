package match

import (
	"foosball/internal/player"
	"time"
)

type Match struct {
	ID uint

	TeamA  []player.Player `gorm:"index;not null"`
	TeamB  []player.Player `gorm:"index;not null"`
	ScoreA int             `gorm:"index;not null"`
	ScoreB int             `gorm:"index;not null"`
	Winner string          `gorm:"index;not null"`

	CreatedAt time.Time
}
