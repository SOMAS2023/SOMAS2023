package physics

import (
	utils "SOMAS2023/internal/common/utils"
	"math"
)

/*
The Engine is responsible for calculating physics for the environment
*/

func CalcAcceleration(f float64, m float64, v float64) float64 {
	if m == 0 {
		panic("zero mass")
	}
	return (f - CalcDrag(v)) / m
}

func CalcDrag(velocity float64) float64 {
	return utils.DragCoefficient * math.Pow(velocity, 2)
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

// ComputeOrientation is to compute the orientation from source coordinate to target coordinate
func ComputeOrientation(src utils.Coordinates, target utils.Coordinates) float64 {
	xDiff := target.X - src.X
	yDiff := target.Y - src.Y
	return math.Atan2(yDiff, xDiff) / math.Pi
}

// ComputeDistance is to compute the L2 distance from source to target
func ComputeDistance(src utils.Coordinates, target utils.Coordinates) float64 {
	return math.Pow(src.X-target.X, 2) + math.Pow(src.Y-target.Y, 2)
}

// This function is to be called from the server only
func GenerateNewState(initialState utils.PhysicalState, force float64, orientation float64) utils.PhysicalState {
	acceleration := CalcAcceleration(force, initialState.Mass, initialState.Velocity)
	velocity := CalcVelocity(acceleration, initialState.Velocity)
	coordinates := GetNewPosition(initialState.Position, velocity, orientation)

	finalState := utils.PhysicalState{
		Position:     coordinates,
		Acceleration: acceleration,
		Velocity:     velocity,
		Mass:         initialState.Mass,
	}

	return finalState
}
