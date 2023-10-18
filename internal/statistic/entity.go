package statistic

import "time"

type MatchResult int

const (
	ResultWin MatchResult = iota
	ResultLoss
	ResultDraw
)

type Measure string

const (
	MeasureWins   Measure = "wins"
	MeasureStreak Measure = "streak"
)

type Statistic struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null"`
	GameID uint `gorm:"not null"`

	Wins   int
	Draws  int
	Losses int
	Streak int

	CreatedAt time.Time
}
