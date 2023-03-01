package models

import (
	"time"

	"gorm.io/datatypes"
)

type Player struct {
	ID uint `gorm:"primaryKey"`

	Name    string `gorm:"index;not null" default:"Unknown Player"`
	Rating  int    `gorm:"index;not null" default:"1000"`
	Matches datatypes.JSON

	CreatedAt time.Time
}
