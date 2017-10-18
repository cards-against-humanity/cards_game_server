package user

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

// User .
type User struct {
	id       int
	googleID string
	name     string
	email    string
}

// GetID .
func (u User) GetID() int {
	return u.id
}

// GetByID .
func GetByID(id int, db *sql.DB) (User, error) {
	rows, err := db.Query("SELECT googleId, name, email FROM users WHERE id = " + strconv.Itoa(id))
	if err != nil {
		return User{}, err
	}
	var googleID string
	var name string
	var email string
	rows.Next()
	if err := rows.Scan(&googleID, &name, &email); err != nil {
		log.Fatal(err)
	}
	return User{id: id, googleID: googleID, name: name, email: email}, nil
}

// GetByCookie .
func GetByCookie(c string, db *sql.DB) (User, error) {
	uid := GetIDByCookie(c, db)
	return GetByID(uid, db)
}

// GetIDByCookie fetches the user ID associated with a cookie from database
func GetIDByCookie(c string, db *sql.DB) int {
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
