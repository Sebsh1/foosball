package match

import (
	"time"

	"gorm.io/datatypes"
)

type Match struct {
	ID uint `gorm:"primaryKey"`

	OrganiziationID uint `gorm:"not null"`

	TeamA datatypes.JSON `gorm:"not null"` // Marshalled JSON list of user IDs
	TeamB datatypes.JSON `gorm:"not null"` // Marshalled JSON list of user IDs

	Sets []string `gorm:"serializer:json;not null"`

	CreatedAt time.Time
}
