package main

import (
	"math/rand"
	"time"

	"./server"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	server.StartHTTP()
}
