package card

import (
	"math/rand"
	"time"
)

// Shuffle randomizes a deck's order
func Shuffle(s *[]interface{}) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range *s {
		swap(i, randInt(i, len(*s)), s)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func swap(i1 int, i2 int, s *[]interface{}) {
	temp := (*s)[i1]
	(*s)[i1] = (*s)[i2]
	(*s)[i2] = temp
}
