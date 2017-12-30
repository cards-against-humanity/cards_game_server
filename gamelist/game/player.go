package game

import (
	"errors"

	"../../card"
	"../../user"
)

// Player a user that also contains other game-specific player data
type Player struct {
	User      user.User `json:"user"`
	Score     int       `json:"score"`
	HasPlayed bool      `json:"hasPlayed"`
}

type player struct {
	user  user.User
	hand  []card.WhiteCard
	score int
}

func (g *Game) getPrivatePlayer(pID int) (player, error) {
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

func (g Game) getPublicPlayers() []Player {
	pl := []Player{}
	for _, p := range g.Players {
		pl = append(pl, g.getPublicPlayerFromPrivate(p))
	}
	return pl
}
