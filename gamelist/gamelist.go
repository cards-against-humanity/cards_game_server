package gamelist

import (
	"errors"

	"../card"
	"../server/socket"
	"../user"
	"./game"
)

// GameList a group of games where each game has a unique name
type GameList struct {
	socketHandler *socket.Handler
	gamesByName   map[string]*game.Game
	gamesByUserID map[int]*game.Game
}

// CreateGameList constructor, generates an empty game list
func CreateGameList(socketHandler *socket.Handler) GameList {
	return GameList{
		socketHandler: socketHandler,
		gamesByName:   make(map[string]*game.Game),
		gamesByUserID: make(map[int]*game.Game),
	}
}

// CreateGame creates a new game with the given name and cards
func (gl *GameList) CreateGame(u user.User, name string, maxPlayers int, bc []card.BlackCard, wc []card.WhiteCard) error {
	if _, exists := gl.gamesByName[name]; exists {
		return errors.New("Game name is taken")
	}
	gl.LeaveGame(u)
	game, err := game.CreateGame(name, maxPlayers, wc, bc, gl.socketHandler)
	if err != nil {
		return err
	}
	game.Join(u)
	gl.gamesByName[name] = game
	gl.gamesByUserID[u.ID] = game
	return nil
}

// StartGame starts the game that the user is in (if they are the game owner)
func (gl *GameList) StartGame(uID int) error {
	if userGame, exists := gl.gamesByUserID[uID]; exists {
		err := userGame.Start(uID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("User is not in a game")
}

// StopGame starts the game that the user is in (if they are the game owner)
func (gl *GameList) StopGame(uID int) error {
	if userGame, exists := gl.gamesByUserID[uID]; exists {
		err := userGame.Stop(uID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("User is not in a game")
}

// GetStateForUser returns a game state from the perspective of a particular user
func (gl *GameList) GetStateForUser(u user.User) *game.UserState {
	userGame, exists := gl.gamesByUserID[u.ID]
	if !exists {
		return nil
	}
	state := userGame.GetState(u.ID)
	return &state
}

// JoinGame adds a user to a particular game
func (gl *GameList) JoinGame(u user.User, gn string) error {
	oldGame, _ := gl.gamesByUserID[u.ID]
	newGame, _ := gl.gamesByName[gn]
	if newGame == nil {
		return errors.New("Game does not exist")
	}
	if oldGame != nil && oldGame.Name == newGame.Name {
		return errors.New("You are already in this game")
	}
	if len(newGame.Players) >= newGame.MaxPlayers {
		return errors.New("Game is full")
	}
	if oldGame == nil {
		gl.LeaveGame(u)
	}
	newGame.Join(u)
	gl.gamesByUserID[u.ID] = newGame
	return nil
}

// LeaveGame removes a user from a particular game
func (gl *GameList) LeaveGame(u user.User) {
	if game, inGame := gl.gamesByUserID[u.ID]; inGame {
		game.Leave(u.ID)
		delete(gl.gamesByUserID, u.ID)
		if len(game.Players) == 0 {
			delete(gl.gamesByName, game.Name)
		}
	}
}

// KickUser kicks a user from the game if the kicker is the game owner
func (gl *GameList) KickUser(owner user.User, uID int) {
	if game, inGame := gl.gamesByUserID[owner.ID]; inGame {
		game.KickUser(owner.ID, uID)
	}
}

// PlayCard allows user to play a card if they are not the judge
func (gl *GameList) PlayCard(u user.User, cID int) error {
	if game, inGame := gl.gamesByUserID[u.ID]; inGame {
		err := game.PlayCard(u.ID, cID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Cannot play a card when you are not in a game")
}

// VoteCard allows user to pick a favorite card
func (gl *GameList) VoteCard(judge user.User, cID int) {
	if game, inGame := gl.gamesByUserID[judge.ID]; inGame {
		game.VoteCard(judge.ID, cID)
	}
}

// GetList fetches a list of all current games
func (gl *GameList) GetList() []game.GenericState {
	list := []game.GenericState{}
	for _, game := range gl.gamesByName {
		list = append(list, game.GetGenericState())
	}
	return list
}
