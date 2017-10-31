package game

import (
	"../../card"
	"../../user"
)

// Player a user that also contains other game-specific player data
type Player struct {
	user.User `json:"user"`
	Score     int `json:"score"`
}

type player struct {
	user  user.User
	hand  []card.WhiteCard
	score int
}
