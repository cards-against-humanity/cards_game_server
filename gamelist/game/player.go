package game

import (
	"../../card"
	"../../user"
)

// Player a user that also contains other game-specific player data
type Player struct {
	user.User `json:"user"`
	Hand      []card.WhiteCard
	Score     int `json:"score"`
}

func newPlayer(user user.User) Player {
	return Player{User: user, Hand: []card.WhiteCard{}, Score: 0}
}
