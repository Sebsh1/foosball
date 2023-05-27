package tournament

import (
	"matchlog/internal/match"
	"matchlog/internal/user"
	"time"

	"gorm.io/datatypes"
)

type Format string

const (
	FormatSingleElimination Format = "single_elimination"
	FormatDoubleElimination Format = "double_elimination"
	FormatRoundRobin        Format = "round_robin"
	FormatSwiss             Format = "swiss"
	FormatBuchholz          Format = "buchholz"
)

type Tournament struct {
	ID             uint
	OrganizationID uint

	Name   string
	Format Format
	Rated  bool

	Users     []user.User
	Matches   []match.Match
	Structure datatypes.JSON

	CreatedAt time.Time
}
