package user

import (
	"time"
)

type Role string

const (
	AdminRole   Role = "admin"
	ManagerRole Role = "manager"
	MemberRole  Role = "member"
	VirtualRole Role = "virtual"
	NoneRole    Role = "none"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Email string `gorm:"index"`
	Name  string `gorm:"index;not null"`
	Hash  string

	OrganizationID *uint `gorm:"index"`
	Role           Role  `default:"none"`

	CreatedAt time.Time
}
