package models

import (
	"time"
)

type Match struct {
	ID       uint `gorm:"primaryKey"`
	SeasonID uint `gorm:"index;not null"`

	TeamAID uint
	TeamBID uint
	TeamA   Team `gorm:"index;not null;foreignKey:TeamAID"`
	TeamB   Team `gorm:"index;not null;foreignKey:TeamBID"`
	GoalsA  int  `gorm:"index;not null"`
	GoalsB  int  `gorm:"index;not null"`

	CreatedAt time.Time
}
