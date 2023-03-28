package models

import (
	"time"
)

type Match struct {
	ID       uint `gorm:"primaryKey"`
	SeasonID uint `gorm:"index;not null"`

	TeamAID uint
	TeamBID uint
	TeamA   Team `gorm:"index;not null;foreignKey:ID;references:TeamAID"`
	TeamB   Team `gorm:"index;not null;foreignKey:ID;references:TeamBID"`
	GoalsA  int  `gorm:"not null"`
	GoalsB  int  `gorm:"not null"`

	CreatedAt time.Time
}
