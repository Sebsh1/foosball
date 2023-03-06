package models

import (
	"time"
)

type Tournament struct {
	ID       uint `gorm:"primaryKey"`
	SeasonID uint `gorm:"index"`

	Teams   []Team  `gorm:"not null"`
	Matches []Match `gorm:"not null"`

	CreatedAt time.Time
}
