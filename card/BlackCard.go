package card

import (
	"time"
)

// BlackCard .
type BlackCard struct {
	Card
}

// CreateBlackCard generates a card struct
func CreateBlackCard(id int, text string, answerFields int, createdAt time.Time, updatedAt time.Time, cardpackID int) BlackCard {
	return BlackCard{Card: Card{ID: id, Type: "black", Text: text, AnswerFields: answerFields, CreatedAt: createdAt, UpdatedAt: updatedAt, CardpackID: cardpackID}}
}
