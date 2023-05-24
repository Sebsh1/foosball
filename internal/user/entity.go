package user

import (
	"foosball/internal/organization"
	"foosball/internal/rating"
	"foosball/internal/statistics"
	"time"
)

type Role string

const (
	AdminRole  Role = "admin"
	MemberRole Role = "member"
	NoRole     Role = ""
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Email string `gorm:"index;not null"`
	Name  string `gorm:"index;not null"`
	Hash  string `gorm:"not null"`

	OrganizationID *uint `gorm:"index"`
	Organization   organization.Organization
	Role           Role

	Rating     rating.Rating
	Statistics statistics.Statistics

	CreatedAt time.Time
}
