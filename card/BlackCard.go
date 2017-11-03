package card

// BlackCard .
type BlackCard struct {
	Card
}

// CreateBlackCard generates a card struct
func CreateBlackCard(id int, text string, answerFields int, cardpackID int) BlackCard {
	return BlackCard{Card: Card{ID: id, Type: "black", Text: text, AnswerFields: answerFields, CardpackID: cardpackID}}
}
