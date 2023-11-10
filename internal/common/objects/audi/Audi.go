package audi

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"
)

type IAudi interface {
	objects.IBaseBiker
}

// MegaBike will have the following forces
type Audi struct {
	*objects.BaseBiker
}

func (audi *Audi) DecideForce(gameState utils.IGameState) {
	pedal := 0.0
	brake := 0.0
	turning := 0.0

	audiBike, e := getAudiBike(gameState)
	if e != nil || audiBike == nil {
		panic("unexpected error happened when Audi get audiBike")
	}

	targetBike, e := searchForTarget(gameState)
	if e != nil {
		panic("unexpected error happened when Audi searching for target")
	}

	// the speed of audi should be constant, use force 1 to represent
	if targetBike == nil { // no target, audi will stop
		if audiBike.GetVelocity() > 0.0 {
			brake = 1.0
		}
		audi.SetForces(utils.Forces{pedal, brake, turning})
	}

	turning, e = computeForTurning(audiBike.GetPosition(), targetBike.GetPosition())
	if e != nil {
		panic("unexpected error happened when Audi compute for turing")
	}
	if audiBike.GetVelocity() == 0.0 { // acc only when audi is not moving
		pedal = 1.0
	}
	audi.SetForces(utils.Forces{pedal, brake, turning})
}

func searchForTarget(state utils.IGameState) (*objects.MegaBike, error) {
	// TODO waiting for func GetGameState() to be implemented
	panic("Unimplemented func")
}

func getAudiBike(state utils.IGameState) (*objects.MegaBike, error) {
	// TODO waiting for func GetGameState() to be implemented
	panic("Unimplemented func")
}

func computeForTurning(src utils.Coordinates, target utils.Coordinates) (float64, error) {
	xDiff := target.X - src.X
	yDiff := target.Y - src.Y
	return math.Atan(yDiff/xDiff) / math.Pi, nil
}
