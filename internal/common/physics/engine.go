package physics

import (
	utils "SOMAS2023/internal/common/utils"
	"math"
)

/*
The Engine is responsible for calculating physics for the environment
*/

func CalcAcceleration(f float64, m float64) float64 {
	return f / m
}

func CalcVelocity(acc float64, currVelocity float64) float64 {
	var newVelocity float64
	// dt is equal to one
	if (currVelocity + (acc * 1)) < 0 {
		newVelocity = 0.0
	} else {
		newVelocity = (acc * 1) + currVelocity
	}
	return newVelocity
}

func GetNewPosition(coordinates utils.Coordinates, velocity float64, orientation float64) utils.Coordinates {
	coordinates.X += velocity * float64(math.Cos(float64(math.Pi*orientation)))
	coordinates.Y += velocity * float64(math.Sin(float64(math.Pi*orientation)))
	return coordinates
}

