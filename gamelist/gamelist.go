package gamelist

import (
	"errors"

	"../card"
	"../user"
	"./game"
)

// GameList a group of games where each game has a unique name
type GameList struct {
	gamesByName   map[string]game.Game
	gamesByUserID map[int]game.Game
}

// CreateGameList constructor, generates an empty game list
func CreateGameList() GameList {
	return GameList{gamesByName: make(map[string]game.Game), gamesByUserID: make(map[int]game.Game)}
}

// CreateGame creates a new game with the given name and cards
func (gl GameList) CreateGame(u user.User, name string, maxPlayers int, bc []card.BlackCard, wc []card.WhiteCard) error {
	if _, exists := gl.gamesByName[name]; exists {
		return errors.New("Game name is taken")
	}
	gl.LeaveGame(u)
	game := game.CreateGame(name, maxPlayers, wc, bc)
	game.Join(u)
	gl.gamesByName[name] = game
	gl.gamesByUserID[u.GetID()] = game
	return nil
}

// GetStateForUser returns a game state from the perspective of a particular user
func (gl GameList) GetStateForUser(u user.User) game.UserState {
	userGame, exists := gl.gamesByUserID[u.GetID()]
	if !exists {
		// TODO - Properly configure blank game state
		return game.UserState{}
	}
	return userGame.GetState(u.GetID())
}

// JoinGame adds a user to a particular game
func (gl GameList) JoinGame(u user.User, gn string) error {
	game, exists := gl.gamesByName[gn]
	if !exists {
		return errors.New("Game does not exist")
	}
	if len(game.Players) >= game.MaxPlayers {
		return errors.New("Game is full")
	}
	gl.LeaveGame(u)
	game.Join(u)
	gl.gamesByUserID[u.GetID()] = game
	return nil
}

// LeaveGame removes a user from a particular game
func (gl GameList) LeaveGame(u user.User) {
	// TODO - Look at this method and finish implementing it
	// TODO - Pause game if player count drops below a certain threshold
	// TODO - If user was the owner/judge, reassign
	if game, inGame := gl.gamesByUserID[u.GetID()]; inGame {
		game.Leave(u.GetID())
		delete(gl.gamesByUserID, u.GetID())
		if len(game.Players) == 0 {
			delete(gl.gamesByName, game.Name)
		}
	}
}

// KickUser kicks a user from the game if the kicker is the game owner
func (gl GameList) KickUser(owner user.User, u user.User) {
	// TODO - Implement this method
}

// PlayCard .
func (gl GameList) PlayCard(u user.User, c card.WhiteCard) {
	// TODO - Implement this method
}

// VotePlayer allows user to pick a favorite card
func (gl GameList) VotePlayer(judge user.User, u user.User) {
	// TODO - Implement this method
}
