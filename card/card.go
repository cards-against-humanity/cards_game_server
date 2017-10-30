package card

import (
	"time"
)

// Card .
type Card struct {
	ID           int       `json:"id"`
	Type         string    `json:"type"`
	Text         string    `json:"text"`
	AnswerFields int       `json:"answerFields,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CardpackID   int       `json:"cardpackId"`
}
