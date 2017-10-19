package game

import (
	"time"

	"../../card"
	"../../user"
)

// Game - A cards game
type Game struct {
	Name         string
	MaxPlayers   int
	Players      []Player
	ownerID      int
	judgeID      int
	stage        int
	nextStage    time.Time
	whiteDraw    []card.WhiteCard
	whiteDiscard []card.WhiteCard
	BlackDraw    []card.BlackCard
	BlackDiscard []card.BlackCard
	BlackCurrent card.BlackCard
}

// CreateGame .
func CreateGame(name string, maxPlayers int, whiteCards []card.WhiteCard, blackCards []card.BlackCard) Game {
	return Game{Name: name, MaxPlayers: maxPlayers, Players: []Player{}}
}

// UserState - The state of a game for a particular user
type UserState struct {
	Name         string           `json:"name"`
	BlackCurrent card.BlackCard   `json:"blackCard"`
	WhiteCards   []card.WhiteCard `json:"whiteCards"`
	JudgeID      int              `json:"judgeId"`
	OwnerID      int              `json:"ownerId"`
	Players      []Player         `json:"players"`
	Hand         []card.WhiteCard `json:"hand"`
}

// GetState returns the game state for a particular user
func (g Game) GetState(uid int) UserState {
	return UserState{}
}

// Join .
func (g Game) Join(u user.User) {
	// Check UserID
	// TODO - Assign owner if first player
}

// Leave .
func (g Game) Leave(uid int) {
}

// PlayCard .
func (g Game) PlayCard(uid int, cid int) error {
	return nil
}
