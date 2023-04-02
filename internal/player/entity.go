package player

import (
	"time"
)

type Player struct {
	ID uint `gorm:"primaryKey"`

	Name   string `gorm:"index;not null;default:Unknown Player"`
	Rating int    `gorm:"default:1000;not null"`

	CreatedAt time.Time
}
