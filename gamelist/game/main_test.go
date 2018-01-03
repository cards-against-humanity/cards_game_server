package game

import (
	"fmt"
	"strings"
	"testing"

	"../../card"
	"../../server/socket"
	"../../user"
)

func TestGame(t *testing.T) {
	var game *Game
	var err error

	// Game can be created without error
	game, err = createMockGame()
	if err != nil {
		t.Errorf("Failed to create game: %s", err.Error())
	}

	// Attempt to start game without enough players
	game, _ = createMockGame()
	game.Join(createMockUser(1))
	err = game.Start(1)
	if err == nil {
		t.Errorf("Expected an error to be returned when attempting to start a game without enough players")
	} else if !(strings.HasPrefix(err.Error(), "Need ") && strings.HasSuffix(err.Error(), " more players to start")) {
		t.Errorf("Received the wrong error")
	}

	// Starting game as owner and attempting to start not as owner
	game, _ = createMockGame()
	for i := 1; i <= 10; i++ {
		game.Join(createMockUser(i))
	}
	err = game.Start(10)
	if err == nil {
		t.Errorf("Expected an error to be returned when attempting to start a game while not being the owner")
	} else if err.Error() != "Only the owner can start the game" {
		t.Errorf("Received the wrong error")
	}
	err = game.Start(1)
	if err != nil {
		t.Errorf("Should not receive an error when starting a game as the owner")
	}

	// Judge ID should be assigned to an existing player when the game is started
	game, _ = createMockGame()
	game.Join(createMockUser(7))
	game.Join(createMockUser(8))
	game.Join(createMockUser(9))
	game.Join(createMockUser(10))
	game.Start(7)
	if game.judgeID < 7 || game.judgeID > 10 {
		t.Errorf("Judge ID should be assigned to a player that is in the game (expected user ID of between 7 and 10 but got %v)", game.judgeID)
	}

	// Game should stop when player count drops below a certain threshold
	game, _ = createMockGame()
	for i := 1; i <= minStartPlayers; i++ {
		game.Join(createMockUser(i))
	}
	game.Start(1)
	if !game.isRunning() {
		t.Errorf("Game should be running")
	}
	game.Leave(1)
	if game.isRunning() {
		t.Errorf("Game should not be running")
	}

	// Should give cards to players once the game is started
	game, _ = createMockGame()
	for i := 1; i <= minStartPlayers; i++ {
		game.Join(createMockUser(i))
	}
	game.Start(1)
	for _, p := range game.Players {
		if len(p.hand) != handSize {
			t.Errorf(fmt.Sprintf("Player should have %v cards in their hand, but instead had %v", handSize, len(p.hand)))
		}
	}

	// Game should progress to judge phase once all players have selected a card
	game, _ = createMockGame()
	for i := 1; i <= minStartPlayers; i++ {
		game.Join(createMockUser(i))
	}
	game.Start(1)
	for i := 2; i <= minStartPlayers; i++ {
		player, _ := game.getPrivatePlayer(i)
		err := game.PlayCard(i, player.hand[0].ID)
		fmt.Println(err)
		if err != nil {
			t.Errorf("User should be able to play a card")
		}
	}
	if game.stage != 2 {
		t.Errorf("Game should automatically enter judge phase once all players have selected cards")
	}

	// Judge should not be able to play a card
	game, _ = createMockGame()
	for i := 1; i <= minStartPlayers; i++ {
		game.Join(createMockUser(i))
	}
	game.Start(1)
	err = game.PlayCard(1, 1)
	if err.Error() != "Cannot play cards as the judge" {
		t.Errorf("Judge should not be able to play a card")
	}

	// Judge can select a card
	game, _ = createMockGame()
	for i := 1; i <= minStartPlayers; i++ {
		game.Join(createMockUser(i))
	}
	game.Start(1)
	for i := 2; i <= minStartPlayers; i++ {
		player, _ := game.getPrivatePlayer(i)
		err := game.PlayCard(i, player.hand[0].ID)
		fmt.Println(err)
		if err != nil {
			t.Errorf("User should be able to play a card")
		}
	}
	err = game.VoteCard(1, game.whitePlayed[2][0].ID)
	if err != nil {
		t.Errorf("Judge should be able to vote a card: " + err.Error())
	}
	if game.stage != 3 {
		t.Errorf("Game should progress to scoring stage after the judge has voted")
	}
}

func createMockUser(id int) user.User {
	return user.User{ID: id, Name: "Tommy Volk", Email: "tvolk131@gmail.com"}
}
func createMockGame() (*Game, error) {
	return CreateGame("test game", 4, getMockWhiteCards(100), getMockBlackCards(100), socket.CreateHandler())
}

func getMockBlackCards(length int) []card.BlackCard {
	cards := []card.BlackCard{}
	for i := 0; i < length; i++ {
		cards = append(cards, card.CreateBlackCard(i, string(i), 1, 1))
	}
	return cards
}
func getMockWhiteCards(length int) []card.WhiteCard {
	cards := []card.WhiteCard{}
	for i := 0; i < length; i++ {
		cards = append(cards, card.CreateWhiteCard(i, string(i), 1))
	}
	return cards
}
