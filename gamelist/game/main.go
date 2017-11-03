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
	"../../user"
)

// Game - A cards game
type Game struct {
	Name         string
	MaxPlayers   int
	Players      []player
	ownerID      int
	judgeID      int
	stage        int
	nextStage    time.Time
	whiteDraw    []card.WhiteCard
	whiteDiscard []card.WhiteCard
	whitePlayed  map[int][]card.WhiteCard // Maps user IDs to an array of cards they played this round
	BlackDraw    []card.BlackCard
	BlackDiscard []card.BlackCard
	BlackCurrent card.BlackCard
}

// UserState - The state of a game for a particular user
type UserState struct {
	Name              string                   `json:"name"`
	BlackCard         card.BlackCard           `json:"blackCard,omitempty"`
	WhiteCardsUnknown [][]card.WhiteCard       `json:"whiteCardsUnknown,omitempty"`
	WhiteCardsKnown   map[int][]card.WhiteCard `json:"whiteCardsKnown,omitempty"`
	JudgeID           int                      `json:"judgeId,omitempty"`
	OwnerID           int                      `json:"ownerId"`
	Players           []Player                 `json:"players"`
	Hand              []card.WhiteCard         `json:"hand,omitempty"`
	CurrentStage      int                      `json:"currentStage,omitempty"`
	NextStage         time.Time                `json:"nextStage,omitempty"`
}

// GenericState - The state of a game for a user that is not in the game
type GenericState struct {
	Name  string    `json:"name"`
	Owner user.User `json:"owner"`
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
	return Game{Name: name, MaxPlayers: maxPlayers, Players: []player{}, stage: 0}, nil
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

// Join .
func (g *Game) Join(u user.User) {
	if !g.playerIsInGame(u.ID) {
		g.Players = append(g.Players, player{user: u, hand: []card.WhiteCard{}, score: 0})
	}
	// TODO - Assign owner if first player
}

// Leave .
func (g *Game) Leave(pID int) {
}

// PlayCard .
func (g *Game) PlayCard(pID int, cid int) error {
	return nil
}

// GetGenericState returns a simple generic state for a game
func (g *Game) GetGenericState() GenericState {
	owner, _ := g.getPrivatePlayer(g.ownerID)
	return GenericState{
		Name:  g.Name,
		Owner: owner.user,
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
