package main

import (
	"database/sql"
	"fmt"
	"os"

	"./server"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "student"
	DB_PASSWORD = "student"
	DB_NAME     = "cards"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(11)
	}
	defer db.Close()
	fmt.Println("Successfully connected to database!")

	server.StartHTTP(db)
}
