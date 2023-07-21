package leaderboard

type LeaderboardType string

const (
	TypeWins   LeaderboardType = "wins"
	TypeStreak LeaderboardType = "streak"
	TypeRating LeaderboardType = "rating"
)

type Placement struct {
	Value  float64 `json:"value"`
	UserID uint    `json:"user_id"`
	Name   string  `json:"name"`
}

type Leaderboard struct {
	Type       LeaderboardType `json:"type"`
	Placements []Placement     `json:"placements"`
}
