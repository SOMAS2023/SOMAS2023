package iextendedagent

/*

The IExtendedAgent is an extension to the BaseAgent class.

Every BaseExtendedAgent will have to inherit IExtendedAgent to exist in the environment. All other objects such as Audi and LootBox must inherit IObject

*/

import (
	BaseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	// "SOMAS2023/internal/common/baseclient"
)

type IExtendedAgent interface {
	// Inherits all functionality from BaseAgent i.e. GetID() and UpdateAgentInternalState()
	BaseAgent.IAgent[IExtendedAgent]

	// Returns if agent is alive or dead
	IsAlive() bool

	// Current Coordinates of Agent
	GetCoordinates() [2]float64

	// sets colour of agent for loot boxes
	SetColour(lootBoxColour Colour)

	// gets colour of agent
	GetColour() Colour

	// This method overrides the BaseAgent's UpdateAgentInternalState.
	// After making a decision, the agent must return a force each round.
	// The forces are forces["pedalForce"], forces["brakeForce"] and forces["turningForce"]
	// The pedal and brake forces is a float from 0.0 to 1.0, the turning force is a float from -1.0 (90° left) to 1.0 (90° right)
	UpdateAgentInternalState() [3]float64
}
