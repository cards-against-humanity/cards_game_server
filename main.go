package main

import (
	"./server"

	_ "github.com/lib/pq"
)

func main() {
	server.StartHTTP()
}
