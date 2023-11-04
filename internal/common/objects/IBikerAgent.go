package ibikeragent

/*

The IBikerAgent is an extension to the BaseAgent class that will be implemented by BikerAgent

*/

import (
	baseagent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	// "SOMAS2023/internal/common/environment"
	// "SOMAS2023/internal/common/utils"
)

type IBikerAgent interface {
	// Inherits all functionality from BaseAgent i.e. GetID() and UpdateAgentInternalState()
	baseagent.IAgent[IBikerAgent]
	IEnvironementObject[IBikerAgent]

	// Sets colour of agent for loot boxes
	SetColour(lootBoxColour Colour)

	// Gets colour of agent
	GetColour() Colour

	// Gets the current force, which will be called by the server after the agents have called UpdateForces()
	// during the UpdateAgentInternalState() method call.
	// The forces are[pedal, brake, turning]
	// The pedal and brake forces is a float from 0.0 to 1.0, the turning force is a float from -1.0 (90° left) to 1.0 (90° right)
	GetForces() Forces

	// After making a decision, the agent must update their force each round.
	// The forces are forces["pedalForce"], forces["brakeForce"] and forces["turningForce"]
	// The pedal and brake forces is a float from 0.0 to 1.0, the turning force is a float from -1.0 (90° left) to 1.0 (90° right)
	UpdateForces()
}
