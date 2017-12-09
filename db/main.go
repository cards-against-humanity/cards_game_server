package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

const (
	dbUser     = "student"
	dbPassword = "student"
	dbName     = "cards"
)

var instance *sql.DB
var once sync.Once

// GetInstance - Accesses the database singleton driver
func GetInstance() *sql.DB {
	once.Do(func() {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
		db, err := sql.Open("postgres", dbinfo)
		if err != nil {
			fmt.Println("Error connecting to database:", err)
			os.Exit(11)
		}
		fmt.Println("Successfully connected to database!")

		instance = db
	})
	return instance
}
