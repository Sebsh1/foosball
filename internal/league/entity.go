package league

import "time"

type League struct {
	Id uint `gorm:"primaryKey"`

	Name string `gorm:"not null"`

	CreatedAt time.Time
}

type LeaguesMatches struct {
	Id uint `gorm:"primaryKey"`

	LeagueId uint `gorm:"primaryKey"`
	MatchId  uint `gorm:"primaryKey"`

	CreatedAt time.Time
}
