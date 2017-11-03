package card

// Card .
type Card struct {
	ID           int    `json:"id"`
	Type         string `json:"type"`
	Text         string `json:"text"`
	AnswerFields int    `json:"answerFields,omitempty"`
	CardpackID   int    `json:"cardpackId"`
}
