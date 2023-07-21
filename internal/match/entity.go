package match

import (
	"time"
)

type Result rune

const (
	TeamAWins Result = 'A'
	TeamBWins Result = 'B'
	Draw      Result = 'D'
)

type Match struct {
	ID uint `gorm:"primaryKey"`

	TeamA  []uint   `gorm:"serializer:json;not null"`
	TeamB  []uint   `gorm:"serializer:json;not null"`
	Sets   []string `gorm:"serializer:json;not null"`
	Result rune     `gorm:"not null"`

	CreatedAt time.Time
}
