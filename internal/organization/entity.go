package organization

import (
	"time"
)

type Organization struct {
	ID uint `gorm:"primaryKey"`

	Name         string `gorm:"index;not null"`
	RatingMethod string `gorm:"not null"`

	CreatedAt time.Time
}
