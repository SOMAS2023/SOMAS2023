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

func (t5 *team5Agent) ProposeDirection() uuid.UUID {
	return lootBoxPref(t5.GetGameState(), t5)
}

func (t5 *team5Agent) GetGameState() objects.IGameState {
	return t5.BaseBiker.GetGameState()
}

func (t5 *team5Agent) GetMegaBikeId() uuid.UUID {
	return t5.BaseBiker.GetMegaBikeId()
}

func (t5 *team5Agent) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	gameState := t5.GetGameState()
	finalPreferences := CalculateLootBoxPreferences(gameState, t5 /*t5.cumulativePreferences*/)

	finalVote := SortPreferences(finalPreferences)

	// fmt.Print("finalVote: ")
	// fmt.Print(finalVote)
	// fmt.Print("\n")

	return finalVote
}

func (t5 *team5Agent) GetMegaBike() []objects.IBaseBiker {
    return t5.BaseBiker.GetGameState().GetMegaBikes()[t5.BaseBiker.GetMegaBikeId()].GetAgents()
}
