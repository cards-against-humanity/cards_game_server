package game

// Game Stages
// -----------
// 0. Not running
// 1. Card play phase
// 2. Judge phase
// 3. Scoring phase

import (
	"errors"
	"strconv"
	"time"

	"../../card"
	"../../server/socket"
	"../../user"
)

const handSize = 8
const minStartPlayers = 4 // Minimum number of players in a game in order to start playing

// Game - A cards game
type Game struct {
	Name          string
	MaxPlayers    int
	Players       []player
	ownerID       int
	judgeID       int
	stage         int
	nextStage     *time.Time
	socketHandler *socket.Handler
	timer         *time.Timer
	whiteDraw     []card.WhiteCard
	whiteDiscard  []card.WhiteCard
	whitePlayed   map[int][]card.WhiteCard // Maps user IDs to an array of cards they played this round
	BlackDraw     []card.BlackCard
	BlackDiscard  []card.BlackCard
	BlackCurrent  *card.BlackCard
}

// UserState - The state of a game for a particular user
type UserState struct {
	Name              string                   `json:"name"`
	BlackCard         *card.BlackCard          `json:"blackCard"`
	WhiteCardsUnknown [][]card.WhiteCard       `json:"whiteCardsUnknown,omitempty"`
	WhiteCardsKnown   map[int][]card.WhiteCard `json:"whiteCardsKnown,omitempty"`
	JudgeID           int                      `json:"judgeId,omitempty"`
	OwnerID           int                      `json:"ownerId"`
	Players           []Player                 `json:"players"`
	Hand              []card.WhiteCard         `json:"hand"`
	CurrentStage      int                      `json:"currentStage,omitempty"`
	NextStage         *time.Time               `json:"nextStage"`
}

// GenericState - The state of a game for a user that is not in the game
type GenericState struct {
	Name  string    `json:"name"`
	Owner user.User `json:"owner"`
}

// CreateGame - Basic constructor for game struct
func CreateGame(name string, maxPlayers int, whiteCards []card.WhiteCard, blackCards []card.BlackCard, socketHandler *socket.Handler) (*Game, error) {
	if len(name) > 64 {
		return &Game{}, errors.New("Game name must not exceed 64 characters")
	}
	// TODO - Get min black card count from config file instead of hardcoding to 10
	if len(blackCards) < 10 {
		return &Game{}, errors.New("Insufficient number of black cards")
	}
	if len(whiteCards) < (maxPlayers * handSize) {
		return &Game{}, errors.New("Insufficient number of white cards")
	}
	if maxPlayers < 3 {
		return &Game{}, errors.New("Max players must be at least 3")
	}
	if maxPlayers > 20 {
		return &Game{}, errors.New("Max players must not exceed 20")
	}
	game := Game{
		Name:          name,
		MaxPlayers:    maxPlayers,
		socketHandler: socketHandler,
		whiteDraw:     whiteCards,
		BlackDraw:     blackCards,
	}
	return &game, nil
}

// GetState - Gets the game state for a particular user (will return spectator game state if user is not in the game)
func (g *Game) GetState(pID int) UserState {
	player, _ := g.getPrivatePlayer(pID)
	knownCards := make(map[int][]card.WhiteCard)
	unknownCards := [][]card.WhiteCard{}

	knownCards[pID] = g.whitePlayed[pID]
	if g.stage == 2 {
		for id, c := range g.whitePlayed {
			if id != pID {
				unknownCards = append(unknownCards, c)
			}
		}
	} else if g.stage == 3 {
		for id, c := range g.whitePlayed {
			knownCards[id] = c
		}
	}

	return UserState{
		Name:              g.Name,
		BlackCard:         g.BlackCurrent,
		WhiteCardsUnknown: unknownCards,
		WhiteCardsKnown:   knownCards,
		JudgeID:           g.judgeID,
		OwnerID:           g.ownerID,
		Players:           g.getPublicPlayers(),
		Hand:              player.hand,
		CurrentStage:      g.stage,
		NextStage:         g.nextStage,
	}
}

// Start - Kicks off the game loop if the user is the game owner
func (g *Game) Start(uID int) error {
	if g.ownerID != uID {
		return errors.New("Only the owner can start the game")
	}
	if g.isRunning() {
		return errors.New("Game is already running")
	}
	if len(g.Players) < minStartPlayers {
		return errors.New("Need " + strconv.Itoa(minStartPlayers-len(g.Players)) + " more players to start")
	}
	g.fillPlayerHands()
	g.next()
	return nil
}

// Stop - Ends the game loop and resets user hands, card decks, etc...
func (g *Game) Stop(uID int) error {
	if g.ownerID != uID {
		return errors.New("Only the owner can stop the game")
	}
	if !g.isRunning() {
		return errors.New("Game is not running")
	}
	g.stop()
	return nil
}

// Join - Adds a user to the game
func (g *Game) Join(u user.User) {
	if !g.playerIsInGame(u.ID) {
		g.Players = append(g.Players, player{user: u, hand: []card.WhiteCard{}, score: 0})
		if len(g.Players) == 1 {
			g.ownerID = u.ID
		}
		// TODO - Give player a hand if game is in progress
	}
	g.updateUserStates()
}

// Leave - Removes a user from the game
func (g *Game) Leave(pID int) {
	for i, p := range g.Players {
		if p.user.ID == pID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			if pID == g.ownerID {
				if len(g.Players) == 0 {
					g.ownerID = 0
				} else {
					g.ownerID = g.Players[0].user.ID
				}
			}
			if cards, exists := g.whitePlayed[pID]; exists {
				g.whiteDiscard = append(g.whiteDiscard, cards...)
				delete(g.whitePlayed, pID)
			}
			if pID == g.judgeID {
				g.judgeID = g.Players[0].user.ID
				if g.stage == 1 || g.stage == 2 {
					g.resetPlayedCards()
					g.fillPlayerHands()
					g.setNextBlackCard()
					g.stage = 1
					g.next()
				}
			}
			if len(g.Players) < minStartPlayers && g.isRunning() {
				g.stop()
			}
			break
		}
	}
	g.updateUserStates()
}

// KickUser - Allows the game owner to boot users from the game
func (g *Game) KickUser(ownerID int, userID int) {
	if ownerID == g.ownerID && ownerID != userID {
		g.Leave(userID)
	}
}

// PlayCard - Plays a card for a particular user
func (g *Game) PlayCard(pID int, cID int) error {
	if g.stage != 1 {
		return errors.New("Cannot play cards at this time")
	}
	if g.judgeID == pID {
		return errors.New("Cannot play cards as the judge")
	}
	if len(g.whitePlayed[pID]) >= g.BlackCurrent.AnswerFields {
		if g.BlackCurrent.AnswerFields == 1 {
			return errors.New("You have already played a card this round")
		}
		return errors.New("You have already played all cards for this round")
	}
	card := g.popCardFromHand(pID, cID)
	if card == nil {
		return errors.New("Cannot play a card that you do not have")
	}

	g.whitePlayed[pID] = append(g.whitePlayed[pID], *card)
	if g.allUsersHavePlayed() {
		g.next()
	}
	g.updateUserStates()
	return nil
}

// VoteCard - Allows the game judge to pick their favorite card
func (g *Game) VoteCard(judgeID int, cardID int) error {
	if g.judgeID != judgeID {
		return errors.New("Must be the judge to upvote a card")
	}
	if g.stage != 2 {
		return errors.New("Cannot upvote a card at this time")
	}
	playerID := g.cardPlayerID(cardID)
	if playerID == -1 {
		return errors.New("This card was not played by anyone")
	}

	player, err := g.getPrivatePlayer(playerID)
	if err != nil {
		return err
	}

	player.score++
	g.next()
	return nil
}

// GetGenericState - Returns a simple generic state for a game
func (g *Game) GetGenericState() GenericState {
	owner, _ := g.getPrivatePlayer(g.ownerID)
	return GenericState{
		Name:  g.Name,
		Owner: owner.user,
	}
}

func (g *Game) stop() {
	if g.isRunning() {
		g.timer.Stop()
	}

	for _, p := range g.Players {
		g.whiteDraw = append(g.whiteDraw, p.hand...)
		p.hand = []card.WhiteCard{}
	}

	g.judgeID = 0
	g.stage = 0
	g.nextStage = nil
	g.timer = nil

	g.resetPlayedCards()
	g.whiteDraw = append(g.whiteDraw, g.whiteDiscard...)
	g.whiteDiscard = []card.WhiteCard{}
	g.whitePlayed = make(map[int][]card.WhiteCard)

	g.resetBlackDeck()
	g.updateUserStates()
}

func (g *Game) next() {
	if g.timer != nil {
		g.timer.Stop()
	}
	switch g.stage {
	case 0:
		g.resetBlackDeck()
		g.resetWhiteDeck()
		g.judgeID = g.Players[0].user.ID

		interval := time.Duration(30) * time.Second
		g.timer = time.AfterFunc(interval, g.next)
		time := time.Now().Add(interval)
		g.nextStage = &time

		g.BlackCurrent = &(g.BlackDraw[len(g.BlackDraw)-1])
		g.BlackDraw = g.BlackDraw[:len(g.BlackDraw)-1]
		g.fillPlayerHands()
		break
	case 1:
		interval := time.Duration(30) * time.Second
		g.timer = time.AfterFunc(interval, g.next)
		time := time.Now().Add(interval)
		g.nextStage = &time
		break
	case 2:
		interval := time.Duration(30) * time.Second
		g.timer = time.AfterFunc(interval, g.next)
		time := time.Now().Add(interval)
		g.nextStage = &time
		break
	case 3:
		interval := time.Duration(30) * time.Second
		g.timer = time.AfterFunc(interval, g.next)
		time := time.Now().Add(interval)
		g.nextStage = &time

		g.resetPlayedCards()
		g.fillPlayerHands()
		g.setNextBlackCard()
		// Set judgeID
		break
	}

	g.timer = time.AfterFunc(time.Duration(5)*time.Second, g.next)

	g.stage++
	if g.stage == 4 {
		g.stage = 1
	}

	g.updateUserStates()
}

///////////////////////
//// -- Helpers -- ////
///////////////////////

// userHasPlayed returns whether a user has played the correct number of cards for this round
func (g Game) userHasPlayed(pID int) bool {
	if g.BlackCurrent == nil {
		return false
	}
	return g.BlackCurrent.AnswerFields == len(g.whitePlayed[pID])
}

func (g Game) playerIsInGame(pID int) bool {
	for _, p := range g.Players {
		if p.user.ID == pID {
			return true
		}
	}
	return false
}

func (g *Game) isRunning() bool {
	return g.timer != nil
}

func (g *Game) updateUserStates() {
	for _, u := range g.Players {
		g.socketHandler.SendActionToUser(u.user.ID, socket.Action{Type: "game/SET_GAME_STATE", Payload: g.GetState(u.user.ID)})
	}
}

func (g *Game) fillPlayerHands() {
	for i := 0; i < len(g.Players); i++ {
		for len(g.Players[i].hand) < handSize {
			if len(g.whiteDraw) == 0 {
				g.shuffleWhiteDeck()
			}
			card := g.whiteDraw[len(g.whiteDraw)-1]
			g.whiteDraw = g.whiteDraw[:len(g.whiteDraw)-1]
			g.Players[i].hand = append(g.Players[i].hand, card)
		}
	}
}

func (g *Game) resetPlayedCards() {
	for i, l := range g.whitePlayed {
		g.whiteDiscard = append(g.whiteDiscard, l...)
		delete(g.whitePlayed, i)
	}
}

func (g *Game) resetBlackDeck() {
	g.BlackDraw = append(g.BlackDraw, g.BlackDiscard...)
	g.BlackDiscard = []card.BlackCard{}
	if g.BlackCurrent != nil {
		g.BlackDraw = append(g.BlackDraw, *g.BlackCurrent)
		g.BlackCurrent = nil
	}
	card.ShuffleBlackDeck(&g.BlackDraw)
}

func (g *Game) resetWhiteDeck() {
	g.resetPlayedCards()
	g.shuffleWhiteDeck()
}

func (g *Game) shuffleWhiteDeck() {
	g.whiteDraw = append(g.whiteDraw, g.whiteDiscard...)
	g.whiteDiscard = []card.WhiteCard{}
	card.ShuffleWhiteDeck(&g.whiteDraw)
}

func (g *Game) setNextBlackCard() {
	if len(g.BlackDraw) == 0 {
		g.BlackDraw = g.BlackDiscard
		g.BlackDiscard = []card.BlackCard{}
		card.ShuffleBlackDeck(&g.BlackDraw)
	}
	if g.BlackCurrent != nil {
		g.BlackDiscard = append(g.BlackDiscard, *g.BlackCurrent)
	}
	g.BlackCurrent = &g.BlackDraw[len(g.BlackDraw)-1]
	g.BlackDraw = g.BlackDraw[:len(g.BlackDraw)-1]
}

func (g *Game) popCardFromHand(pID int, cID int) *card.WhiteCard {
	player, err := g.getPrivatePlayer(pID)
	if err != nil {
		return nil
	}

	for i, c := range player.hand {
		if c.ID == cID {
			card := player.hand[i]
			player.hand = append(player.hand[:i], player.hand[i+1:]...)
			return &card
		}
	}
	return nil
}

func (g Game) allUsersHavePlayed() bool {
	for _, c := range g.whitePlayed {
		if len(c) < g.BlackCurrent.AnswerFields {
			return false
		}
	}
	return true
}

// Returns the user ID of the user who has played a particular card this round (returns -1 if no one played that card)
func (g *Game) cardPlayerID(cID int) int {
	for _, cl := range g.whitePlayed {
		for _, c := range cl {
			if c.ID == cID {
			}
		}
	}
	return -1
}
