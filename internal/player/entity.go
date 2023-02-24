package player

import (
	"time"
)

type Player struct {
	ID uint

	Name   string `gorm:"index;not null" default:"Unknown Player"`
	Rating int    `gorm:"index;not null" default:"1000"`

	CreatedAt time.Time
}
