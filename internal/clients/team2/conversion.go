package team2

import (
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

func forcesToVectorConversion(force utils.Forces) ForceVector {
	xCoordinate := force.Pedal * float64(math.Cos(float64(math.Pi*force.Turning.SteeringForce)))
	yCoordinate := force.Pedal * float64(math.Sin(float64(math.Pi*force.Turning.SteeringForce)))

	newVector := ForceVector{X: xCoordinate, Y: yCoordinate}
	return newVector
}

func dotProduct(v1, v2 ForceVector) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func magnitude(v ForceVector) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func cosineSimilarity(v1, v2 ForceVector) float64 {
	return dotProduct(v1, v2) / (magnitude(v1) * magnitude(v2))
}

// Normalizes the values of a map to be between 0 and 1.
// Useful when calculating social capital from Networks, Reputation, and Institutions
func normalizeMapValues(capitals map[uuid.UUID]float64) map[uuid.UUID]float64 {
	if len(capitals) == 0 {
		return capitals
	}

	min := 0.0
	max := 0.0
	for _, value := range capitals {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	for key, value := range capitals {
		capitals[key] = (value - min) / (max - min)
	}

	return capitals
}
