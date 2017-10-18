package card

import (
	"fmt"
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

// Print prints card to console
func (c BlackCard) Print() {
	fmt.Println("------------")
	fmt.Printf("ID: %v\n", c.ID)
	fmt.Printf("Text: '%v'\n", c.Text)
	fmt.Printf("Answer Fields: %v\n", c.AnswerFields)
	fmt.Printf("Created on %v\n", c.CreatedAt.Format(time.UnixDate))
	fmt.Printf("Updated on %v\n", c.UpdatedAt.Format(time.UnixDate))
	fmt.Printf("Cardpack ID: %v\n", c.CardpackID)
	fmt.Println("------------")
}
