package card

import (
	"fmt"
	"time"
)

// BlackCard .
type BlackCard struct {
	id           int
	text         string
	answerFields int
	createdAt    time.Time
	updatedAt    time.Time
	cardpackID   int
}

// Print prints card to console
func (c BlackCard) Print() {
	fmt.Println("------------")
	fmt.Printf("ID: %v\n", c.id)
	fmt.Printf("Text: '%v'\n", c.text)
	fmt.Printf("Answer Fields: %v\n", c.answerFields)
	fmt.Printf("Created on %v\n", c.createdAt.Format(time.UnixDate))
	fmt.Printf("Updated on %v\n", c.updatedAt.Format(time.UnixDate))
	fmt.Printf("Cardpack ID: %v\n", c.cardpackID)
	fmt.Println("------------")
}
