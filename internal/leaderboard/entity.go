package leaderboard

type LeaderboardType string

const (
	TypeWins          LeaderboardType = "wins"
	TypeWinStreak     LeaderboardType = "win-streak"
	TypeLossStreak    LeaderboardType = "loss-streak"
	TypeWinLossRatio  LeaderboardType = "win-loss-ratio"
	TypeRating        LeaderboardType = "rating"
	TypeMatchesPlayed LeaderboardType = "matches-played"
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
