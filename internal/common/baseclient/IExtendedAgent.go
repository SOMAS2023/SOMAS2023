package iextendedagent

/*

The IExtendedAgent is an extension to the BaseAgent class that will be implemented by BaseExtendedAgent

*/

import (
	baseagent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	// "SOMAS2023/internal/common/baseclient"
)

type IExtendedAgent interface {
	// Inherits all functionality from BaseAgent i.e. GetID() and UpdateAgentInternalState()
	baseagent.IAgent[IExtendedAgent]

	// Sets colour of agent for loot boxes
	SetColour(lootBoxColour Colour)

	// Gets colour of agent
	GetColour() Colour
}
