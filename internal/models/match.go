package models

import (
	"time"

	"gorm.io/datatypes"
)

type Match struct {
	ID uint `gorm:"primaryKey"`

	TeamA  datatypes.JSONType[[]uint] `gorm:"not null"`
	TeamB  datatypes.JSONType[[]uint] `gorm:"not null"`
	GoalsA int                        `gorm:"index;not null"`
	GoalsB int                        `gorm:"index;not null"`

	CreatedAt time.Time
}
