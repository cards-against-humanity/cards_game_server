package user

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"../db"
)

// User .
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetByID fetches a user from the database with a given ID
func GetByID(id int) (User, error) {
	rows, err := db.GetInstance().Query("SELECT name, email FROM users WHERE id = " + strconv.Itoa(id))
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

// GetByRequest gets the user that sent the HTTP request
func GetByRequest(r *http.Request) (User, error) {
	cookie, e := r.Cookie("connect.sid")
	if e != nil {
		return User{}, e
	}
	return GetByCookie(cookie.Value)
}

// GetByCookie fetches a user from the database associated with a given cookie
func GetByCookie(c string) (User, error) {
	uid, e := getIDByCookie(c)
	if e != nil {
		return User{}, e
	}
	return GetByID(uid)
}

func parseUserID(data string) int {
	split := strings.Split(data, `}`)
	data = split[len(split)-3]
	data = strings.Split(data, `"user":`)[1]
	i, _ := strconv.ParseInt(data, 10, 64)
	return int(i)
}
