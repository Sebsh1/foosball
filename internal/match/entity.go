package match

import (
	"time"

	"gorm.io/datatypes"
)

type Set struct {
	ID uint `gorm:"primaryKey"`

	MatchID uint `gorm:"not null"`
	PointsA int  `gorm:"not null" json:"pointsA"`
	PointsB int  `gorm:"not null" json:"pointsB"`
}

type Match struct {
	ID uint `gorm:"primaryKey"`

	TeamA datatypes.JSON `gorm:"not null"` // Marshalled JSON list of user IDs
	TeamB datatypes.JSON `gorm:"not null"` // Marshalled JSON list of user IDs

	Sets []Set `gorm:"not null"`

	CreatedAt time.Time
}
