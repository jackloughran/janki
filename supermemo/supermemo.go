package supermemo

import (
	"math"
	"time"

	"github.com/jackloughran/janki/janki"
)

func AdjustCard(card janki.Card, quality int) janki.Card {
	newEfficiencyFactor := card.EfficiencyFactor + (0.1 - (5-float64(quality))*(0.08+(5-float64(quality))*0.02))

	if newEfficiencyFactor < 1.3 {
		newEfficiencyFactor = 1.3
	}

	var newRepititions int
	if quality < 3 {
		newRepititions = 1
	} else {
		newRepititions = card.Repititions + 1
	}
	interval := calculateInterval(newRepititions, newEfficiencyFactor)

	return janki.Card{
		ID:               card.ID,
		Front:            card.Front,
		Back:             card.Back,
		EfficiencyFactor: newEfficiencyFactor,
		Repititions:      newRepititions,
		NextRepitition:   time.Now().Add(time.Duration(interval*24) * time.Hour),
	}
}

func calculateInterval(repititions int, efficiencyFactor float64) int {
	if repititions == 1 {
		return 1
	} else if repititions == 2 {
		return 6
	} else {
		return int(math.Ceil(float64(calculateInterval(repititions-1, efficiencyFactor)) * efficiencyFactor))
	}
}
