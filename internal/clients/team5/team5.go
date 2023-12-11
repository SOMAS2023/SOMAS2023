package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

type Iteam5Agent interface {
	objects.BaseBiker
}

type team5Agent struct {
	*objects.BaseBiker
	resourceAllocMethod ResourceAllocationMethod
	state               int
	prevEnergy          map[uuid.UUID]float64
	roundCount          int
	otherBikerForces    map[uuid.UUID]utils.Forces
	otherBikerRep       map[uuid.UUID]float64
	finalPreferences    map[uuid.UUID]float64
}

type ResourceAllocationMethod int

const (
	Equal ResourceAllocationMethod = iota
	Greedy
	Needs
	Contributions
	Reputation
)

// Creates an instance of Team 5 Biker
func GetBiker(baseBiker *objects.BaseBiker) objects.IBaseBiker {
	baseBiker.GroupID = 5
	// fmt.Println("team5Agent: newTeam5Agent: baseBiker: ", baseBiker)
	return &team5Agent{
		BaseBiker:           baseBiker,
		resourceAllocMethod: Equal,
		state:               1, //observer state
		roundCount:          0,
		otherBikerForces:    make(map[uuid.UUID]utils.Forces),
		otherBikerRep:       make(map[uuid.UUID]float64),
		finalPreferences:    make(map[uuid.UUID]float64),
	}
}

func (t5 *team5Agent) UpdateAgentInternalState() {
	t5.updateState()
	t5.updateReputationOfAllAgents()
	t5.roundCount = (t5.roundCount + 1) % utils.RoundIterations
}

func (t5 *team5Agent) DecideGovernance() utils.Governance {
	return utils.Democracy
}

func (t5 *team5Agent) DecideAction() objects.BikerAction {
	return objects.Pedal
}

func (t5 *team5Agent) ChangeBike() uuid.UUID {
	//get reputation of all bikes
	bikeReps := t5.getReputationOfAllBikes()
	//get ID for maximum reputation bike if the bike is not full (<8 agents)
	maxRep := 0.0
	maxRepID := uuid.Nil
	for bikeID, rep := range bikeReps {
		//get length from GetAgents()
		numAgentsOnbike := len(t5.GetGameState().GetMegaBikes()[bikeID].GetAgents())
		if rep > maxRep && numAgentsOnbike < 8 {
			maxRep = rep
			maxRepID = bikeID
		}
	}
	return maxRepID
}

func (t5 *team5Agent) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	gameState := t5.GetGameState()
	finalPreferences := t5.CalculateLootBoxPreferences(gameState, proposals /*t5.cumulativePreferences*/)

	finalVote := SortPreferences(finalPreferences)

	return finalVote
}

func (t5 *team5Agent) DecideAllocation() voting.IdVoteMap {
	//fmt.Println("team5Agent: GetBike: t5.BaseBiker.DecideAllocation: ", t5.resourceAllocationMethod)
	method := t5.resourceAllocMethod
	return t5.calculateResourceAllocation(method)
}

func (t5 *team5Agent) VoteDictator() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	fellowBikers := t5.GetFellowBikers()
	var value float64 = 0
	for _, fellowBiker := range fellowBikers {
		value = t5.QueryReputation(fellowBiker.GetID())
		if fellowBiker.GetColour() == t5.GetColour() {
			value += 1
		}

		votes[fellowBiker.GetID()] = value
	}
	return votes
}

func (t5 *team5Agent) VoteLeader() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	fellowBikers := t5.GetFellowBikers()
	var value float64 = 0
	for _, fellowBiker := range fellowBikers {
		value = t5.QueryReputation(fellowBiker.GetID())
		if fellowBiker.GetColour() == t5.GetColour() {
			value += 1
		}

		votes[fellowBiker.GetID()] = value
	}
	return votes

}
