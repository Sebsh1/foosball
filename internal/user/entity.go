package user

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`

	Email string `gorm:"index;not null"`
	Name  string `gorm:"index;not null"`
	Hash  string `gorm:"not null"`

	OrganizationID uint
	Rating         int `gorm:"default:1000"`
	Admin          bool

	CreatedAt time.Time
}
