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
	Position int
	Value    float64
	UserID   uint
	Name     string
}

type Leaderboard struct {
	Type       LeaderboardType
	Placements []Placement
}
