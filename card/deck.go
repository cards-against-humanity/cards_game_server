package card

import (
	"math/rand"
	"time"
)

// Deck .
type Deck []Card

// Shuffle randomizes a deck's order
func (d Deck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range d {
		d.swap(i, randInt(i, len(d)))
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (d Deck) swap(i1 int, i2 int) {
	temp := d[i1]
	d[i1] = d[i2]
	d[i2] = temp
}
