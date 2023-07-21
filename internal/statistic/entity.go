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
	ID uint

	Wins   int
	Draws  int
	Losses int
	Streak int

	CreatedAt time.Time
}
