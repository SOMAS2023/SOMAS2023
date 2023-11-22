package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"github.com/google/uuid"
)

type Iteam5Agent interface {
	objects.BaseBiker
}

type team5Agent struct {
	objects.BaseBiker
	resourceAllocationMethod string
}

func NewTeam5Agent(totColours utils.Colour, bikeId uuid.UUID) *team5Agent {
	baseBiker := objects.GetBaseBiker(totColours, bikeId) // Use the constructor function
	// print
	fmt.Println("team5Agent: newTeam5Agent: baseBiker: ", baseBiker)
	return &team5Agent{
		BaseBiker:                *baseBiker,
		resourceAllocationMethod: "equal",
	}
}

func (t5 *team5Agent) GetBike() uuid.UUID {
	fmt.Println("team5Agent: GetBike: t5.BaseBiker.GetBike(): ", t5.BaseBiker.GetBike())
	return t5.BaseBiker.GetBike()
}

func (t5 *team5Agent) DecideAllocation() voting.IdVoteMap {
	fmt.Println("team5Agent: GetBike: t5.BaseBiker.DecideAllocation: ", t5.resourceAllocationMethod)
	return calculateResourceAllocation(t5.GetGameState(), t5)
}

func (t5 *team5Agent) ProposeDirection(pendingAgents map[uuid.UUID]float64) uuid.UUID {
    gameState := t5.GetGameState()
    agentID := t5.GetID()

    // Calculate the final preferences for all loot boxes
    finalPreferences := CalculateLootBoxPreferences(gameState, agentID, t5.cumulativePreferences)

    // Find the loot box with the highest preference
    var bestBox uuid.UUID
    maxPreference := -1.0
    for boxID, preference := range finalPreferences {
        if preference > maxPreference {
            bestBox = boxID
            maxPreference = preference
        }
    }

    return bestBox
}


