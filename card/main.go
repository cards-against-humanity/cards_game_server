package card

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// GetCards .
func GetCards(cpids []int, db *sql.DB) ([]BlackCard, []WhiteCard) {
	_, err := db.Query("SELECT id FROM cardpacks WHERE" + generateOrQuery(`"id"`, intSliceToStringSlice(cpids)))
	if err != nil {
		fmt.Printf("Error: one or more cardpack ID is invalid - %v", cpids)
		return nil, nil
	}
	rows, err := db.Query("SELECT * FROM cards WHERE" + generateOrQuery(`"cardpackId"`, intSliceToStringSlice(cpids)))
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
		var createdAt string
		var updatedAt string
		var cardpackID int
		if err := rows.Scan(&id, &text, &ctype, &answerFields, &createdAt, &updatedAt, &cardpackID); err != nil {
			log.Fatal(err)
		}

		if ctype == "black" {
			bc = append(bc, CreateBlackCard(id, text, int(answerFields.Int64), cardpackID))
		} else {
			wc = append(wc, CreateWhiteCard(id, text, cardpackID))
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

func generateOrQuery(fieldName string, elems []string) string {
	query := ""
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
