package main

import (
	"database/sql"
	"fmt"
	"os"

	"./server"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "student:student@tcp(127.0.0.1)/cards")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(11)
	}
	defer db.Close()
	fmt.Println("Successfully connected to database!")

	server.StartHTTP(db)
}
