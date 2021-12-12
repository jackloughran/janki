package janki

import "time"

// Card represents a flashcard
type Card struct {
	ID             int64
	Front          string
	Back           string
	Repititions    int
	NextRepitition time.Time
}
