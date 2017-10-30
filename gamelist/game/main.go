package game

import (
	"errors"
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

// UserState - The state of a game for a particular user
// TODO - Add nextStage to this struct
type UserState struct {
	Name         string           `json:"name"`
	BlackCurrent card.BlackCard   `json:"blackCard"`
	WhiteCards   []card.WhiteCard `json:"whiteCards"`
	JudgeID      int              `json:"judgeId"`
	OwnerID      int              `json:"ownerId"`
	Players      []Player         `json:"players"`
	Hand         []card.WhiteCard `json:"hand"`
}

// GenericState - The state of a game for a user that is not in the game
type GenericState struct {
	Name    string      `json:"name"`
	OwnerID int         `json:"ownerId"`
	Players []user.User `json:"players"`
}

// CreateGame .
func CreateGame(name string, maxPlayers int, whiteCards []card.WhiteCard, blackCards []card.BlackCard) (Game, error) {
	if len(name) > 64 {
		return Game{}, errors.New("Game name must not exceed 64 characters")
	}
	// TODO - Get min black card count from config file instead of hardcoding to 10
	if len(blackCards) < 10 {
		return Game{}, errors.New("Insufficient number of black cards")
	}
	// TODO - Get min white card count from config file instead of hardcoding
	if len(whiteCards) < (maxPlayers * 10) {
		return Game{}, errors.New("Insufficient number of white cards")
	}
	if maxPlayers < 3 {
		return Game{}, errors.New("Max players must be at least 3")
	}
	if maxPlayers > 20 {
		return Game{}, errors.New("Max players must not exceed 20")
	}
	return Game{Name: name, MaxPlayers: maxPlayers, Players: []Player{}}, nil
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
