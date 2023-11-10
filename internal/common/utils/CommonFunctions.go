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

// CheckCollision checks if two coordinates are within a certain distance of each other to determine if there is a collision.
func CheckCollision(c1 Coordinates, c2 Coordinates, epsilon float64) bool {
	// Manually calculate distance between the two coordinates, and if the distance is less than epsilon, then there is a collision.
	return (c1.X-c2.X)*(c1.X-c2.X)+(c1.Y-c2.Y)*(c1.Y-c2.Y) < epsilon*epsilon
}

func GenerateRandomFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
