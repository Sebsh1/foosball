package models

import (
	"time"
)

type Team struct {
	ID uint `gorm:"primaryKey"`

	Players []*Player `gorm:"many2many:player_teams;"`

	CreatedAt time.Time
}
