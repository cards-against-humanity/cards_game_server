package card

import (
	"fmt"
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

// Print prints card to console
func (c WhiteCard) Print() {
	fmt.Println("------------")
	fmt.Printf("ID: %v\n", c.ID)
	fmt.Printf("Text: '%v'\n", c.Text)
	fmt.Printf("Created on %v\n", c.CreatedAt.Format(time.UnixDate))
	fmt.Printf("Updated on %v\n", c.UpdatedAt.Format(time.UnixDate))
	fmt.Printf("Cardpack ID: %v\n", c.CardpackID)
	fmt.Println("------------")
}
