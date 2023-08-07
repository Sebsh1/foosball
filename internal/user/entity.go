package user

import (
	"time"
)

type Role string

const (
	AdminRole  Role = "admin"
	Manager    Role = "manager"
	MemberRole Role = "member"
	NoneRole   Role = ""
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Email string `gorm:"index;not null,unique"`
	Name  string `gorm:"index;not null"`
	Hash  string `gorm:"not null"`

	OrganizationID *uint `gorm:"index"`
	Role           Role  `default:"member"`

	CreatedAt time.Time
}
