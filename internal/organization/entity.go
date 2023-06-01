package organization

import (
	"matchlog/internal/rating"
	"time"
)

type Organization struct {
	ID uint `gorm:"primaryKey"`

	Name         string        `gorm:"index;not null"`
	RatingMethod rating.Method `gorm:"not null"`

	CreatedAt time.Time
}
