package baseextendedagent

/*
	BaseExtendedAgent is the implentation of IExtendedAgent. All agents must be composed of a BaseExtendedAgent
*/

import (
	baseagent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	// "SOMAS2023/internal/common/baseclient"
)

type BaseExtendedAgent struct {
	// implements the base agent class with the IExtendAgent interface
	*baseagent.BaseAgent[IExtendedAgent]

	// implements the IEnvironmentObject agent class with the IEnvironmentObject interface
	*EnvironmentObject[IEnvironementObject]
}

// SetColour sets the agent's colour.
func (ea *ExtendedAgent) SetColour(lootBoxColour Colour) {
	ea.colour = lootBoxColour
}

// GetColour returns the agent's colour.
func (ea *ExtendedAgent) GetColour() Colour {
	return ea.colour
}
