package physics

import (
	//"fmt"
	utils "SOMAS2023/internal/common/utils"
	"math"
	//baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	//"github.com/google/uuid"
)

// since in go you cannot define methods on types from other packages
type Coordinates utils.Coordinates
type Forces utils.Forces

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

func (eng *Engine) Update_loc(coordinates Coordinates, velocity float64, Orientation float64) Coordinates {
	coordinates.X += velocity * float64(math.Cos(float64(math.Pi*Orientation)))
	coordinates.Y += velocity * float64(math.Sin(float64(math.Pi*Orientation)))
	return coordinates
}
