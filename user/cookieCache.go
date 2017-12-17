package user

import (
	"errors"

	"../db"
)

var cookieCache = make(map[string]int)

func getIDByCookie(c string) (int, error) {
	if len(c) > 32 {
		c = c[4:36]
	}
	if _, exists := cookieCache[c]; !exists {
		rows, err := db.GetInstance().Query(`SELECT data FROM "Sessions" WHERE sid = '` + c + `'`)
		if err != nil || !rows.Next() {
			return -1, errors.New("Cookie is not valid")
		}
		var data string
		rows.Scan(&data)
		cookieCache[c] = parseUserID(data)
	}
	return cookieCache[c], nil
}
