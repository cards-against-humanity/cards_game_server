package card

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

// GetCards .
func GetCards(cpids []int, db *sql.DB) ([]BlackCard, []WhiteCard) {
	_, err := db.Query(generateOrQuery("SELECT id FROM cardpacks WHERE", "id", intSliceToStringSlice(cpids)))
	if err != nil {
		fmt.Printf("Error: one or more cardpack ID is invalid - %v", cpids)
		return nil, nil
	}
	rows, err := db.Query(generateOrQuery("SELECT * FROM cards WHERE", "cardpackId", intSliceToStringSlice(cpids)))
	if err != nil {
		fmt.Println("Error reading cards from database:", err)
		return nil, nil
	}

	defer rows.Close()
	bc := []BlackCard{}
	wc := []WhiteCard{}
	for rows.Next() {
		var id int
		var text string
		var ctype string
		var answerFields sql.NullInt64
		var createdAt []uint8
		var updatedAt []uint8
		var cardpackID int
		if err := rows.Scan(&id, &text, &ctype, &answerFields, &createdAt, &updatedAt, &cardpackID); err != nil {
			log.Fatal(err)
		}

		if ctype == "black" {
			bc = append(bc, CreateBlackCard(id, text, int(answerFields.Int64), intsToTime(createdAt), intsToTime(updatedAt), cardpackID))
		} else {
			wc = append(wc, CreateWhiteCard(id, text, intsToTime(createdAt), intsToTime(updatedAt), cardpackID))
		}
	}
	return bc, wc
}

func intsToBytes(nl []uint8) []byte {
	b := make([]byte, len(nl))
	for i, n := range nl {
		b[i] = n
	}
	return b
}

func bytesToTime(b []byte) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05 -0700 MST", string(b)+" +0000 UTC")
}

func intsToTime(nl []uint8) time.Time {
	t, e := bytesToTime(intsToBytes(nl))
	if e != nil {
		log.Fatal(e)
	}
	return t
}

func generateOrQuery(baseQuery string, fieldName string, elems []string) string {
	query := baseQuery
	for i, e := range elems {
		query += " " + fieldName + " = " + e
		if i < len(elems)-1 {
			query += " OR"
		}
	}
	return query
}

func intSliceToStringSlice(ints []int) []string {
	s := make([]string, len(ints))
	for i, n := range ints {
		s[i] = strconv.Itoa(n)
	}
	return s
}
