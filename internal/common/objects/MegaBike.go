package objects

import (
	phy "SOMAS2023/internal/common/physics"
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type IMegaBike interface {
	IPhysicsObject
	AddAgent(biker IBaseBiker)
	RemoveAgent(bikerId uuid.UUID)
	GetAgents() []IBaseBiker
	UpdateMass()
	CalculateForce() float64
	CalculateOrientation() float64
}

// MegaBike will have the following forces
type MegaBike struct {
	*PhysicsObject
	agents []IBaseBiker
}

// GetMegaBike is a constructor for MegaBike that initializes it with a new UUID and default position.
func GetMegaBike() *MegaBike {
	return &MegaBike{
		PhysicsObject: GetPhysicsObject(utils.MassBike),
	}
}

// adds
func (mb *MegaBike) AddAgent(biker IBaseBiker) {
	mb.agents = append(mb.agents, biker)
}

// Remove agent from bike, given its ID
func (mb *MegaBike) RemoveAgent(bikerId uuid.UUID) {
	// Create a new slice to store the updated agents
	var updatedAgents []IBaseBiker

	// Iterate through the agents and copy them to the updatedAgents slice
	for _, agent := range mb.agents {
		if agent.GetID() != bikerId {
			updatedAgents = append(updatedAgents, agent)
		}
	}

	// Replace the mb.agents slice with the updatedAgents slice
	mb.agents = updatedAgents
}

func (mb *MegaBike) GetAgents() []IBaseBiker {
	return mb.agents
}

// Calculate the mass of the bike with all it's agents
func (mb *MegaBike) UpdateMass() {
	mass := utils.MassBike
	mass += float64(len(mb.agents))
	mb.mass = mass
}

// Calculates the total force based on the Biker's force
func (mb *MegaBike) CalculateForce() float64 {
	if len(mb.agents) == 0 {
		return 0.0
	}
	totalPedal := 0.0
	totalBrake := 0.0
	for _, agent := range mb.agents {
		force := agent.GetForces()

		if force.Pedal != 0 {
			totalPedal += float64(force.Pedal)
		} else {
			totalBrake += float64(force.Brake)
		}
	}
	F := (float64(totalPedal) - float64(totalBrake))
	return F
}

// Calculates the final orientation of the Megabike, between -1 and 1 (-180° to 180°), given the Biker's Turning forces
func (mb *MegaBike) CalculateOrientation() float64 {
	if len(mb.agents) == 0 {
		return mb.orientation
	}
	totalTurning := 0.0
	for _, agent := range mb.agents {
		totalTurning += float64(agent.GetForces().Turning)
	}
	averageTurning := totalTurning / float64(len(mb.agents))
	mb.orientation += (averageTurning)
	// ensure the orientation wraps around if it exceeds the range 1.0 or -1.0
	if mb.orientation > 1.0 {
		mb.orientation -= 2
	} else if mb.orientation < -1.0 {
		mb.orientation += 2
	}
	return mb.orientation
}

// Moves the MegaBike to its new position after the agents have applied thier force
func (mb *MegaBike) Move() {
	mb.acceleration = phy.CalcAcceleration(mb.CalculateForce(), mb.mass)
	mb.velocity = phy.CalcVelocity(mb.acceleration, mb.velocity)
	mb.orientation = mb.CalculateOrientation()
	mb.coordinates = phy.GetNewPosition(mb.coordinates, mb.velocity, mb.orientation)
}
