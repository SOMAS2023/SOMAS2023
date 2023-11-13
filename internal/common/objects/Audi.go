package objects

import (
	phy "SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"math"
)

type IAudi interface {
	IPhysicsObject
	UpdateGameState(state IGameState)
}

type Audi struct {
	*PhysicsObject
	target *MegaBike
}

// GetAudi is a constructor for Audi that initializes it with a new UUID and default position.
func GetAudi() *Audi {
	return &Audi{
		PhysicsObject: GetPhysicsObject(utils.MassAudi),
	}
}

func GetIAudi() IAudi {
	return &Audi{
		PhysicsObject: GetPhysicsObject(utils.MassAudi),
	}
}

// Move as MegaBike, called by server
func (audi *Audi) Move() {
	if audi.target == nil { // no target, audi will stop
		audi.velocity = 0.0
	} else {
		audi.velocity = 1.0 / audi.mass
		audi.orientation = computeOrientation(audi.coordinates, audi.target.GetPosition())
	}
	audi.coordinates = phy.GetNewPosition(audi.coordinates, audi.velocity, audi.orientation)
}

func (audi *Audi) UpdateGameState(state IGameState) {
	// search for target
	minDistance := math.Inf(1)
	audi.target = nil
	for _, bike := range state.GetMegaBikes() {
		if bike.GetVelocity() != 0.0 {
			continue
		}
		distance := computeDistance(audi.coordinates, bike.GetPosition())
		if distance < minDistance {
			minDistance = distance
			audi.target = bike
		}
	}
}

func computeOrientation(src utils.Coordinates, target utils.Coordinates) float64 {
	xDiff := target.X - src.X
	yDiff := target.Y - src.Y
	return math.Atan(yDiff/xDiff) / math.Pi
}

func computeDistance(src utils.Coordinates, target utils.Coordinates) float64 {
	return math.Pow(src.X-target.X, 2) + math.Pow(src.Y-target.Y, 2)
}
