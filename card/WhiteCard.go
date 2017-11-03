package card

// WhiteCard .
type WhiteCard struct {
	Card
}

// CreateWhiteCard generates a card struct
func CreateWhiteCard(id int, text string, cardpackID int) WhiteCard {
	return WhiteCard{Card: Card{ID: id, Type: "white", Text: text, CardpackID: cardpackID}}
}
