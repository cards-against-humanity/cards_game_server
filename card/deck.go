package card

import (
	"math/rand"
)

// TODO - Reduce to single functions for both white and black cards

// ShuffleBlackDeck randomizes a deck's order
func ShuffleBlackDeck(s *[]BlackCard) {
	for i := range *s {
		swapB(i, randInt(i, len(*s)), s)
	}
}

// ShuffleWhiteDeck randomizes a deck's order
func ShuffleWhiteDeck(s *[]WhiteCard) {
	for i := range *s {
		swapW(i, randInt(i, len(*s)), s)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func swapB(i1 int, i2 int, s *[]BlackCard) {
	temp := (*s)[i1]
	(*s)[i1] = (*s)[i2]
	(*s)[i2] = temp
}

func swapW(i1 int, i2 int, s *[]WhiteCard) {
	temp := (*s)[i1]
	(*s)[i1] = (*s)[i2]
	(*s)[i2] = temp
}
