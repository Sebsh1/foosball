package tournament

import "time"

type Tournament struct {
	ID uint `gorm:"primaryKey"`

	Name string `gorm:"not null"`

	CreatedAt time.Time
}
