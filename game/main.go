package game

import (
	"time"
)

// Game - A cards game
type Game struct {
	name         string
	maxPlayers   int
	players      []player
	stage        int
	nextStage    time.Time
	whiteDraw    deck
	whiteDiscard deck
}

// UserState - The state of a game for a particular user
type UserState struct {
}

// GetState returns the game state for a particular user
func (g Game) GetState(uid int) UserState {
	return UserState{}
}

func (g Game) addUser(uid int) {
}

func (g Game) removeUser(uid int) {
}

func (g Game) playCard(uid int) error {
	return nil
}
