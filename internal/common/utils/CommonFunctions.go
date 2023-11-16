package utils

import (
	"math/rand"

	"github.com/google/uuid"
)

// GenerateRandomCoordinates creates random X and Y coordinates within the grid boundaries.
func GenerateRandomCoordinates() Coordinates {
	// Generate random coordinates
	return Coordinates{
		X: rand.Float64() * GridWidth,
		Y: rand.Float64() * GridHeight,
	}
}

// GenerateRandomCoordinates creates random X and Y coordinates within the grid boundaries.
func GenerateRandomColour() Colour {
	// Generate a random index between 0 and the number of colours - 1.
	randomIndex := rand.Intn(int(NumOfColours))
	return Colour(randomIndex)
}

func GenerateRandomFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// this function will take in a list of maps from ids to their corresponding vote (yes/ no in the case of acceptance)
// and retunr a list of ids that can be accepted according to some metric (ie more than half voted yes)
// ranked according to a metric (ie overall number of yes's)
func GetAcceptanceRanking([]map[uuid.UUID]bool) []uuid.UUID {
	// TODO implement
	return make([]uuid.UUID, 0)
}

// returns the winner accoring to chosen voting strategy (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func WinnerFromDist([]map[uuid.UUID]float64) uuid.UUID {
	panic("not implemented")
}
