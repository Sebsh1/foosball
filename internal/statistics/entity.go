package statistics

import "time"

type Statistics struct {
	ID uint

	Wins   int
	Losses int
	Draws  int

	WinStreak  int
	LossStreak int
	DrawStreak int

	CreatedAt time.Time
}
