package objects

import (
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

// MegaBike will have the following forces
type MegaBike struct {
	id           uuid.UUID
	coordinates  utils.Coordinates
	agentForces  []utils.Forces
	Mass         float64
	acceleration float64
	velocity     float64
	orientation  float64
}

// GetMegaBike is a constructor for MegaBike that initializes it with a new UUID and default position.
func GetMegaBike() *MegaBike {
	return &MegaBike{
		id:          uuid.New(),                        // Generate a new unique identifier
		coordinates: utils.GenerateRandomCoordinates(), // Initialize to randomized position
	}
}

// returns the unique ID of the object
func (mb *MegaBike) GetID() uuid.UUID {
	return mb.id
}

// returns the current coordinates of the object
func (mb *MegaBike) GetPosition() utils.Coordinates {
	return mb.coordinates
}

// AddAgentForce should be called by the server once per agent, once all Bikers have called DecideForce
func (mb *MegaBike) AddAgentForce(agentForces utils.Forces) {
	mb.agentForces = append(mb.agentForces, agentForces)
}

// Calculates the total force based on the Biker's force
func (mb *MegaBike) CalculateForce() float64 {
	forceMap := 4.0
	if len(mb.agentForces) == 0 {
		return 0.0
	}
	totalPedal := 0.0
	totalBrake := 0.0
	totalMass := utils.MassBike
	for _, agent := range mb.agentForces {
		totalMass += utils.MassBiker
		if agent.Pedal != 0 {
			totalPedal += float64(agent.Pedal)
		} else {
			totalBrake += float64(agent.Brake)
		}
	}
	F := forceMap * (float64(totalPedal) - float64(totalBrake))
	return F
}

// Calculates the final orientation of the Megabike, between -1 and 1 (-180° to 180°), given the Biker's Turning forces
func (mb *MegaBike) CalculateOrientation() float64 {
	if len(mb.agentForces) == 0 {
		return mb.orientation
	}
	totalTurning := 0.0
	for _, agentForce := range mb.agentForces {
		totalTurning += float64(agentForce.Turning)
	}
	averageTurning := totalTurning / float64(len(mb.agentForces))
	mb.orientation += (averageTurning)
	// ensure the orientation wraps around if it exceeds the range 1.0 or -1.0
	if mb.orientation > 1.0 {
		mb.orientation -= 2
	} else if mb.orientation < -1.0 {
		mb.orientation += 2
	}

	return mb.orientation
}
