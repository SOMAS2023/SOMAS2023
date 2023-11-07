package utils

import (
	"math/rand"
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
