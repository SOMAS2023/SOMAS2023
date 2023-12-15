package objects

import (
	phy "SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

type IAwdi interface {
	IPhysicsObject
	UpdateGameState(state IGameState)
	GetTargetID() uuid.UUID
}

type Awdi struct {
	*PhysicsObject
	target    IMegaBike
	gameState IGameState
}

// GetAwdi is a constructor for Awdi that initializes it with a new UUID and default position.
func GetAwdi() *Awdi {
	return &Awdi{
		PhysicsObject: GetPhysicsObject(utils.MassAwdi),
	}
}

func GetIAwdi() IAwdi {
	return &Awdi{
		PhysicsObject: GetPhysicsObject(utils.MassAwdi),
	}
}

// Calculates and returns the desired force of the awdi based on the current gamestate
func (awdi *Awdi) UpdateForce() {
	// Compute the target Megabike, which will update awdi.target
	awdi.ComputeTarget()

	if awdi.target == nil { // no target, awdi will not apply a force and eventually come to a stop
		awdi.force = 0.0
	} else {
		awdi.force = utils.AwdiMaxForce // Otherwise apply max force to get to target MegaBike
	}
}

// Calculates and returns the desired orientation of the awdi based on the current gamestate
func (awdi *Awdi) UpdateOrientation() {
	// If no target, awdi will not change orientation
	// Otherwise, new orientation is calculated based on positioning of target
	if awdi.target != nil {
		awdi.orientation = phy.ComputeOrientation(awdi.coordinates, awdi.target.GetPosition())
	}
}

// Computes the target Megabike based on current gameState
func (awdi *Awdi) ComputeTarget() {
	// search for target
	minDistance := math.Inf(1)
	minVelocity := math.Inf(1)
	awdi.target = nil
	for _, bike := range awdi.gameState.GetMegaBikes() {
		if utils.AwdiOnlyTargetsStationaryMegaBike {
			if bike.GetVelocity() != 0.0 {
				continue
			}
		}

		if !utils.AwdiTargetsEmptyMegaBike {
			agentsOnBike := bike.GetAgents()
			if len(agentsOnBike) == 0 {
				continue
			}
		}

		// ignore faster bike
		if bike.GetVelocity() > minVelocity {
			continue
		}

		distance := phy.ComputeDistance(awdi.coordinates, bike.GetPosition())
		// minimize the velocity first
		if bike.GetVelocity() < minVelocity {
			awdi.target = bike
		} else if distance < minDistance { // if same velocity, then minimize distance
			awdi.target = bike
		} else {
			continue
		}
		minVelocity = awdi.target.GetVelocity()
		minDistance = distance
	}
}

// Updates gameState member variable
func (awdi *Awdi) UpdateGameState(state IGameState) {
	awdi.gameState = state
}

func (awdi *Awdi) GetTargetID() uuid.UUID {
	if awdi.target != nil {
		return awdi.target.GetID()
	} else {
		return uuid.UUID{}
	}
}
