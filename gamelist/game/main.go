package game

// Game Stages
// -----------
// 0. Not running
// 1. Card play phase
// 2. Judge phase
// 3. Scoring phase

import (
	"errors"
	"time"

	"../../card"
	"../../server/socket"
	"../../user"
)

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

// CreateGame .
func CreateGame(name string, maxPlayers int, whiteCards []card.WhiteCard, blackCards []card.BlackCard, socketHandler *socket.Handler) (Game, error) {
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
	return Game{
		Name:          name,
		MaxPlayers:    maxPlayers,
		Players:       []player{},
		socketHandler: socketHandler,
		whiteDraw:     whiteCards,
		whiteDiscard:  []card.WhiteCard{},
		whitePlayed:   make(map[int][]card.WhiteCard),
		BlackDraw:     blackCards,
		BlackDiscard:  []card.BlackCard{},
	}, nil
}

// GetState returns the game state for a particular user (will return generic game state if user is not in the game)
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

// Start .
func (g *Game) Start(uID int) error {
	if g.ownerID != uID {
		return errors.New("Only the owner can start the game")
	}
	if g.timer != nil {
		return errors.New("Game is already running")
	}
	g.next()
	return nil
}

// Join .
func (g *Game) Join(u user.User) {
	if !g.playerIsInGame(u.ID) {
		g.Players = append(g.Players, player{user: u, hand: []card.WhiteCard{}, score: 0})
		if len(g.Players) == 1 {
			g.ownerID = u.ID
		}
	}
}

// Leave .
func (g *Game) Leave(pID int) {
	for i, p := range g.Players {
		if p.user.ID == pID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			if pID == g.ownerID {
				g.ownerID = g.Players[0].user.ID
			}
			if pID == g.judgeID {
				// TODO - Properly reassign judge
			}
			break
		}
	}
	if len(g.Players) < 4 {
		g.stop()
	}
}

// KickUser allows the game owner to boot users from the game
func (g *Game) KickUser(ownerID int, userID int) {
	if ownerID == g.ownerID && ownerID != userID {
		g.Leave(userID)
	}
}

// PlayCard .
func (g *Game) PlayCard(pID int, cID int) error {
	return nil
}

// VoteCard allows the game judge to pick their favorite card
func (g *Game) VoteCard(judgeID int, cardID int) {
}

// GetGenericState returns a simple generic state for a game
func (g *Game) GetGenericState() GenericState {
	owner, _ := g.getPrivatePlayer(g.ownerID)
	return GenericState{
		Name:  g.Name,
		Owner: owner.user,
	}
}

func (g *Game) stop() {
	g.timer.Stop()

	for _, p := range g.Players {
		g.whiteDraw = append(g.whiteDraw, p.hand...)
		p.hand = []card.WhiteCard{}
	}

	g.judgeID = 0
	g.stage = 0
	g.nextStage = nil
	g.timer = nil

	g.whiteDraw = append(g.whiteDraw, g.whiteDiscard...)
	g.whiteDiscard = []card.WhiteCard{}
	for _, list := range g.whitePlayed {
		g.whiteDraw = append(g.whiteDraw, list...)
	}
	g.whitePlayed = make(map[int][]card.WhiteCard)

	g.BlackDraw = append(g.BlackDraw, g.BlackDiscard...)
	g.BlackDraw = append(g.BlackDraw, *g.BlackCurrent)
	g.BlackDiscard = []card.BlackCard{}
	g.BlackCurrent = nil
}

func (g *Game) next() {
	g.timer = time.AfterFunc(time.Duration(5)*time.Second, func() {
	})
	// Update the states of all users in the game
	for _, u := range g.Players {
		g.socketHandler.SendActionToUser(u.user.ID, socket.Action{Type: "game/SET_GAME_STATE", Payload: g.GetState(u.user.ID)})
	}
}

///////////////////////
//// -- Helpers -- ////
///////////////////////

func (g Game) getPrivatePlayer(pID int) (player, error) {
	for _, p := range g.Players {
		if p.user.ID == pID {
			return p, nil
		}
	}
	return player{}, errors.New("User is not in this game")
}

func (g Game) getPublicPlayer(pID int) (Player, error) {
	pPriv, err := g.getPrivatePlayer(pID)
	if err != nil {
		return Player{}, err
	}
	return Player{User: pPriv.user, Score: pPriv.score, HasPlayed: g.userHasPlayed(pID)}, nil
}

func (g Game) getPublicPlayerFromPrivate(pPriv player) Player {
	return Player{User: pPriv.user, Score: pPriv.score, HasPlayed: g.userHasPlayed(pPriv.user.ID)}
}

// userHasPlayed returns whether a user has played the correct number of cards for this round
func (g Game) userHasPlayed(pID int) bool {
	return g.BlackCurrent.AnswerFields == len(g.whitePlayed[pID])
}

func (g Game) getPublicPlayers() []Player {
	pl := []Player{}
	for _, p := range g.Players {
		pl = append(pl, g.getPublicPlayerFromPrivate(p))
	}
	return pl
}

func (g Game) playerIsInGame(pID int) bool {
	for _, p := range g.Players {
		if p.user.ID == pID {
			return true
		}
	}
	return false
}
