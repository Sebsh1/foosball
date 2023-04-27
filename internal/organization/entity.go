package organization

import (
	"foosball/internal/user"
	"time"
)

type Organization struct {
	ID uint `gorm:"primaryKey"`

	Name         string `gorm:"index;not null"`
	Users        []user.User
	RatingMethod string `gorm:"not null"`

	CreatedAt time.Time
}
