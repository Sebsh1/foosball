package match

import (
	"time"
)

type Match struct {
	ID       uint `gorm:"primaryKey"`
	SeasonID uint `gorm:"index"`

	TeamAID uint `gorm:"index;not null"`
	TeamBID uint `gorm:"index;not null"`

	GoalsA int `gorm:"not null"`
	GoalsB int `gorm:"not null"`

	CreatedAt time.Time
}
