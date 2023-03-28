package models

import (
	"time"
)

type Player struct {
	ID uint `gorm:"primaryKey"`

	Name   string  `gorm:"index;default:Unknown Player;not null"`
	Rating int     `gorm:"default:1000;not null"`
	Teams  []*Team `gorm:"many2many:player_teams;"`

	CreatedAt time.Time
}
