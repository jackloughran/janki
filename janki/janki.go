package janki

import "time"

// Card represents a flashcard
type Card struct {
	ID               string
	Front            string
	Back             string
	EfficiencyFactor float64
	Repititions      int
	NextRepitition   time.Time
}
