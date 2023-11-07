package physics

import (
	utils "SOMAS2023/internal/common/utils"
	"math"
)

/*
The Engine struct is responsible for calculating physics for the environment
*/
type Engine struct {
}

func (eng *Engine) CalcAcceleration(f float64, m float64) float64 {
	return f / m
}

func (eng *Engine) CalcVelocity(acc float64, dt float64, currVelocity float64) float64 {
	var newVelocity float64
	if (currVelocity + (acc * dt)) < 0 {
		newVelocity = 0.0
	} else {
		newVelocity = (acc * dt) + currVelocity
	}
	return newVelocity
}

func (eng *Engine) UpdateLoc(coordinates utils.Coordinates, velocity float64, orientation float64) utils.Coordinates {
	coordinates.X += velocity * float64(math.Cos(float64(math.Pi*orientation)))
	coordinates.Y += velocity * float64(math.Sin(float64(math.Pi*orientation)))
	return coordinates
}
