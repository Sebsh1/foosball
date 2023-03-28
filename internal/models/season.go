package models

import (
	"time"
)

type Season struct {
	ID uint `gorm:"primaryKey"`

	Name    string `gorm:"index"`
	Matches []*Match

	Start time.Time `gorm:"not null"`
	End   time.Time `gorm:"not null"`

	CreatedAt time.Time
}
