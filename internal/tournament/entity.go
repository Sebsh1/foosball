package tournament

import "time"

type Tournament struct {
	Id uint `gorm:"primaryKey"`

	Name string `gorm:"not null"`

	CreatedAt time.Time
}
