package user

import (
	"foosball/internal/organization"
	"foosball/internal/rating"
	"foosball/internal/statistic"
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

	RatingID uint
	Rating   rating.Rating

	StatisticsID uint
	Statistics   statistic.Statistic

	CreatedAt time.Time
}
