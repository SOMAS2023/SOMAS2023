package objects

import (
	phy "SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
	"math"
)

type IAudi interface {
	IPhysicsObject
	UpdateGameState(state IGameState)
	GetTargetID() uuid.UUID
}

type Audi struct {
	*PhysicsObject
	target    IMegaBike
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

// Calculates and returns the desired force of the audi based on the current gamestate
func (audi *Audi) UpdateForce() {
	// Compute the target Megabike, which will update audi.target
	audi.ComputeTarget()

	if audi.target == nil { // no target, audi will not apply a force and eventually come to a stop
		audi.force = 0.0
	} else {
		audi.force = utils.AudiMaxForce // Otherwise apply max force to get to target MegaBike
	}
}

// Calculates and returns the desired orientation of the audi based on the current gamestate
func (audi *Audi) UpdateOrientation() {
	// If no target, audi will not change orientation
	// Otherwise, new orientation is calculated based on positioning of target
	if audi.target != nil {
		audi.orientation = phy.ComputeOrientation(audi.coordinates, audi.target.GetPosition())
	}
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
		if !utils.AudiTargetEmpty {
			agentsOnBike := bike.GetAgents()
			if agentsOnBike == nil || len(agentsOnBike) == 0 {
				continue
			}
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

func (audi *Audi) GetTargetID() uuid.UUID {
	if audi.target != nil {
		return audi.target.GetID()
	} else {
		return uuid.UUID{}
	}
}
