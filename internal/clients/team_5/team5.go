package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

type Iteam5Agent interface {
	objects.BaseBiker
}

type team5Agent struct {
	objects.BaseBiker
	resourceAllocMethod ResourceAllocationMethod
	//set state default to 0
	state int // 0 = normal, 1 = conservative
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
func NewTeam5Agent(totColours utils.Colour, bikeId uuid.UUID) *team5Agent {
	baseBiker := objects.GetBaseBiker(totColours, bikeId) // Use the constructor function
	baseBiker.GroupID = 5
	// print
	fmt.Println("team5Agent: newTeam5Agent: baseBiker: ", baseBiker)
	return &team5Agent{
		BaseBiker:           *baseBiker,
		resourceAllocMethod: Equal,
		state:               0,
	}
}

func (t5 *team5Agent) UpdateAgentInternalState() {
	t5.updateReputationOfAllAgents()
	t5.updateState()
}

// needs fixing always democracy
func (t5 *team5Agent) DecideGovernance() utils.Governance {
	return utils.Democracy
}

//Functions can be called in any scenario

// needs fixing never gets off bike
func (t5 *team5Agent) DecideAction() objects.BikerAction {
	return objects.Pedal
}

// needs fixing doesn't pick a bike to join
// Decides which bike to join based on reputation and space available
// Todo: create a formula that combines reputation, space available, people with same colour, governance system (rn only uses rep)
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

// needs fixing currently never votes off

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

func (t5 *team5Agent) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {

	wantToSendMsg := true
	if wantToSendMsg {
		forcesMsg := t5.CreateForcesMessage()
		return []messaging.IMessage[objects.IBaseBiker]{forcesMsg}
	}
	return []messaging.IMessage[objects.IBaseBiker]{}
}

func (t5 *team5Agent) CreateForcesMessage() objects.ForcesMessage {
	fmt.Println()
	return objects.ForcesMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](t5, t5.GetFellowBikers()),
		AgentId:     t5.GetID(),
		AgentForces: utils.Forces{
			Pedal: 1.0,
			Brake: 1.0,
			Turning: utils.TurningDecision{
				SteerBike:     false,
				SteeringForce: 1.0,
			},
		},
	}
}

func (t5 *team5Agent) HandleForcesMessage(msg objects.ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// agentForces := msg.AgentForces

}

func (bb *team5Agent) HandleKickoutMessage(msg objects.KickoutAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// kickout := msg.kickout
}
