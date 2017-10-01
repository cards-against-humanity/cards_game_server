package card

import (
	"fmt"
	"time"
)

// WhiteCard .
type WhiteCard struct {
	id         int
	text       string
	createdAt  time.Time
	updatedAt  time.Time
	cardpackID int
}

// Print prints card to console
func (c WhiteCard) Print() {
	fmt.Println("------------")
	fmt.Printf("ID: %v\n", c.id)
	fmt.Printf("Text: '%v'\n", c.text)
	fmt.Printf("Created on %v\n", c.createdAt.Format(time.UnixDate))
	fmt.Printf("Updated on %v\n", c.updatedAt.Format(time.UnixDate))
	fmt.Printf("Cardpack ID: %v\n", c.cardpackID)
	fmt.Println("------------")
}
