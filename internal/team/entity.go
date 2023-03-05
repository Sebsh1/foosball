package team

import (
	"foosball/internal/player"
	"time"
)

type Team struct {
	ID uint `gorm:"primaryKey"`

	PlayerA *player.Player `gorm:"index;not null"`
	PlayerB *player.Player `gorm:"index"`
	PlayerC *player.Player `gorm:"index"`

	CreatedAt time.Time
}

func (t *Team) GetPlayers() []*player.Player {
	var players []*player.Player

	players = append(players, t.PlayerA)
	if t.PlayerB == nil {
		return players
	}
	players = append(players, t.PlayerB)
	if t.PlayerC == nil {
		return players
	}
	players = append(players, t.PlayerC)
	return players
}
