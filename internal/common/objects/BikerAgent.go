package baseextendedagent

/*
	BikerAgent is the implentation of IBikerAgent. All extended biker agents must be composed of a BikerAgent object
*/

import (
	baseagent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
	// "SOMAS2023/internal/common/utils"
)

type BikerAgent struct {
	colour Colour
	forces Forces
	position [2]float64
}

// generates a new BikerAgent instance with a randomised ID
NewBikerAgent [T IBikerAgent [ T ]]() ∗ BikerAgent [ T ] {
	return &BaseAgent [ T] {
		id uuid.ID
	}
}

// SetColour sets the color of the BikerAgent.
func (ba *BikerAgent) SetColour(lootBoxColour Colour) {
	ba.colour = lootBoxColour
}

// GetColour returns the color of the BikerAgent.
func (ba *BikerAgent) GetColour() Colour {
	return ba.colour
}

// GetForces returns the current forces of the BikerAgent.
func (ba *BikerAgent) GetForces() Forces {
	return ba.forces
}

// UpdateForces updates the forces based on some decision-making process.
func (ba *BikerAgent) UpdateForces() {
	// Implement the decision-making process to update the forces.
	// This is an example and should be replaced with actual logic.
	ba.forces.PedalForce = 0.5    // example value
	ba.forces.BrakeForce = 0.0    // example value
	ba.forces.TurningForce = -0.1 // example value indicating a slight turn to the left
}

// GetID is a method from BaseAgent that needs to be implemented if not already present.
func (ba *BikerAgent) GetID() uuid.UUID {
	// Return the ID from the embedded BaseAgent struct.
	return ba.BaseAgent.GetID()
}

// UpdateAgentInternalState is another method from BaseAgent that should be implemented.
func (ba *BikerAgent) UpdateAgentInternalState() {
	// Implement the logic to update the agent's internal state.
}
