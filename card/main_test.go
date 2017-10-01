package card

import "testing"

func TestGenerateOrQuery(t *testing.T) {
	var q1, q2 string

	q1 = generateOrQuery("SELECT * FROM users WHERE", "ID", []string{"1", "2", "3"})
	q2 = "SELECT * FROM users WHERE ID = 1 OR ID = 2 OR ID = 3"
	if q1 != q2 {
		t.Errorf("Failed: Expected...\n%s\nto equal...\n%s", q1, q2)
	}

	q1 = generateOrQuery("SELECT id FROM card WHERE", "cardpackId", []string{"4"})
	q2 = "SELECT id FROM card WHERE cardpackId = 4"
	if q1 != q2 {
		t.Errorf("Failed: Expected...\n%s\nto equal...\n%s", q1, q2)
	}
}
