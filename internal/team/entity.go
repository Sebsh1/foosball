package team

import (
	"foosball/internal/player"
	"time"
)

type Team struct {
	ID uint `gorm:"primaryKey"`

	Players []*player.Player `gorm:"not null;many2many:players_teams"`

	CreatedAt time.Time
}
