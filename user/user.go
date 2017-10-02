package user

import (
	"database/sql"
	"log"
	"strconv"
)

// User .
type User struct {
	id       int
	googleID int
	name     string
	email    string
}

// Get .
func Get(id int, db *sql.DB) (User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE id = %s", strconv.Itoa(id))
	if err != nil {
		return User{}, err
	}
	var googleID int
	var name string
	var email string
	var themeID int
	var createdAt []uint8
	var updatedAt []uint8
	if err := rows.Scan(id, &googleID, &name, &email, &themeID, &createdAt, &updatedAt); err != nil {
		log.Fatal(err)
	}
	return User{id: id, googleID: googleID, name: name, email: email}, nil
}
