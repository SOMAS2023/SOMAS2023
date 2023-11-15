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
	target    *MegaBike
	gameState IGameState
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
	// Compute the target Megabike, which will update audi.target
	audi.ComputeTarget()

	if audi.target == nil { // no target, audi will stop
		audi.velocity = 0.0
	} else {
		audi.velocity = 1.0 / audi.mass
		audi.orientation = phy.ComputeOrientation(audi.coordinates, audi.target.GetPosition())
	}
	audi.coordinates = phy.GetNewPosition(audi.coordinates, audi.velocity, audi.orientation)
}

// Computes the target Megabike based on current gameState
func (audi *Audi) ComputeTarget() {
	// search for target
	minDistance := math.Inf(1)
	audi.target = nil
	for _, bike := range audi.gameState.GetMegaBikes() {
		if bike.GetVelocity() != 0.0 {
			continue
		}
		distance := phy.ComputeDistance(audi.coordinates, bike.GetPosition())
		if distance < minDistance {
			minDistance = distance
			audi.target = bike
		}
	}
}

// Updates gameState member variable
func (audi *Audi) UpdateGameState(state IGameState) {
	audi.gameState = state
}
