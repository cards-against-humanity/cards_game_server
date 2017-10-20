package card

import (
	"time"
)

// WhiteCard .
type WhiteCard struct {
	Card
}

// CreateWhiteCard generates a card struct
func CreateWhiteCard(id int, text string, createdAt time.Time, updatedAt time.Time, cardpackID int) WhiteCard {
	return WhiteCard{Card: Card{ID: id, Type: "white", Text: text, CreatedAt: createdAt, UpdatedAt: updatedAt, CardpackID: cardpackID}}
}
