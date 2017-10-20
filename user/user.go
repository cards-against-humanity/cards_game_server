package user

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

// User .
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetByID fetches a user from the database with a given ID
func GetByID(id int, db *sql.DB) (User, error) {
	rows, err := db.Query("SELECT name, email FROM users WHERE id = " + strconv.Itoa(id))
	if err != nil {
		return User{}, err
	}
	var name string
	var email string
	rows.Next()
	if err := rows.Scan(&name, &email); err != nil {
		log.Fatal(err)
	}
	return User{ID: id, Name: name, Email: email}, nil
}

// GetByCookie fetches a user from the database associated with a given cookie
func GetByCookie(c string, db *sql.DB) (User, error) {
	uid := getIDByCookie(c, db)
	return GetByID(uid, db)
}

func getIDByCookie(c string, db *sql.DB) int {
	if len(c) > 32 {
		c = c[4:36]
	}
	rows, err := db.Query(`SELECT data FROM sessions WHERE sid = "` + c + `"`)
	if err != nil {
		return -1
	}
	var data string
	rows.Next()
	rows.Scan(&data)
	return parseUserID(data)
}

func parseUserID(data string) int {
	data = strings.Split(data, `"passport":`)[1]
	data = strings.Split(data, `"user":`)[1]
	data = strings.Split(data, `}`)[0]
	i, _ := strconv.ParseInt(data, 10, 64)
	return int(i)
}
