package authentication

import "time"

type User struct {
	ID           uint
	Username     string `gorm:"primaryKey;index;not null"`
	PasswordHash string `gorm:"not null"`

	CreatedAt time.Time
}
