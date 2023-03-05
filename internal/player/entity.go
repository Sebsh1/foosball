package player

import (
	"time"
)

type Player struct {
	ID uint `gorm:"primaryKey"`

	Name   string `gorm:"index;default:Unknown Player;not null"`
	Rating int    `gorm:"index;default:1000;not null"`

	CreatedAt time.Time
}
