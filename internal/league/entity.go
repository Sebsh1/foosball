package league

import "time"

type League struct {
	ID uint `gorm:"primaryKey"`

	Name string `gorm:"not null"`

	CreatedAt time.Time
}

type LeaguesMatches struct {
	ID uint `gorm:"primaryKey"`

	LeagueID uint `gorm:"primaryKey"`
	MatchID  uint `gorm:"primaryKey"`

	CreatedAt time.Time
}
