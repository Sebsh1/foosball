package user

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`

	Email string `gorm:"index;not null"`
	Name  string `gorm:"index;not null"`
	Hash  string `gorm:"not null"`

	OrganizationID *uint
	Admin          bool

	CreatedAt time.Time
}

type UserStats struct {
	ID uint `gorm:"primaryKey"`

	UserID uint
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Rating int `gorm:"default:1000"`
	Wins   int `gorm:"default:0"`
	Losses int `gorm:"default:0"`
	Draws  int `gorm:"default:0"`

	CreatedAt time.Time
}
