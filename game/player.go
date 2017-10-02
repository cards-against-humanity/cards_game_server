package game

import "../user"

type player struct {
	user  user.User
	hand  deck
	score int
}
