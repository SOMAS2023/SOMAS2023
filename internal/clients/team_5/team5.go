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

func (t5 *team5Agent) UpdateAgentInternalState() {
	t5.updateReputationOfAllAgents()
}

func (t5 *team5Agent) DecideAllocation() voting.IdVoteMap {
	//fmt.Println("team5Agent: GetBike: t5.BaseBiker.DecideAllocation: ", t5.resourceAllocationMethod)
	return calculateResourceAllocation(t5.GetGameState(), t5)
}

func (t5 *team5Agent) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	gameState := t5.GetGameState()
	finalPreferences := CalculateLootBoxPreferences(gameState, t5, proposals /*t5.cumulativePreferences*/)

	finalVote := SortPreferences(finalPreferences)

	// fmt.Print("finalVote: ")
	// fmt.Print(finalVote)
	// fmt.Print("\n")

	return finalVote
}

func (t5 *team5Agent) GetAgentsOnMegaBike() []objects.IBaseBiker {
	return t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetAgents()
}
