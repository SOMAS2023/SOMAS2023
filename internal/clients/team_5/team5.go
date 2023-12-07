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
	state            int // 0 = normal, 1 = conservative
	OtherBikerForces []utils.Forces
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
	t5.updateState()
	t5.updateReputationOfAllAgents()
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
// func (t5 *team5Agent) VoteForKickout() map[uuid.UUID]int {
// 	voteResults := make(map[uuid.UUID]int)
// 	for _, agent := range t5.GetFellowBikers() {
// 		agentID := agent.GetID()
// 		if agentID != t5.GetID() {
// 			voteResults[agentID] = 0
// 		}
// 	}
// 	return voteResults
// }

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

func (bb *team5Agent) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages and which ones to send (return)
	println("get all")
	wantToSendMsg := true
	if wantToSendMsg {
		reputationMsg := bb.CreateReputationMessage()
		kickoutMsg := bb.CreatekickoutMessage()
		lootboxMsg := bb.CreateLootboxMessage()
		joiningMsg := bb.CreateJoiningMessage()
		governceMsg := bb.CreateGoverenceMessage()
		forcesMsg := bb.CreateForcesMessage()
		voteGoveranceMessage := bb.CreateVoteGovernanceMessage()
		voteLootboxDirectionMessage := bb.CreateVoteLootboxDirectionMessage()
		voteRulerMessage := bb.CreateVoteRulerMessage()
		voteKickoutMessage := bb.CreateVotekickoutMessage()
		return []messaging.IMessage[objects.IBaseBiker]{reputationMsg, kickoutMsg, lootboxMsg, joiningMsg, governceMsg, forcesMsg, voteGoveranceMessage, voteLootboxDirectionMessage, voteRulerMessage, voteKickoutMessage}
	}
	return []messaging.IMessage[objects.IBaseBiker]{}
}

func (bb *team5Agent) CreatekickoutMessage() objects.KickoutAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	println("create")
	return objects.KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		Kickout:     false,
	}
}

func (bb *team5Agent) CreateReputationMessage() objects.ReputationOfAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		Reputation:  1.0,
	}
}

func (bb *team5Agent) CreateJoiningMessage() objects.JoiningAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.JoiningAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		BikeId:      uuid.Nil,
	}
}
func (bb *team5Agent) CreateLootboxMessage() objects.LootboxMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.LootboxMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		LootboxId:   uuid.Nil,
	}
}

func (bb *team5Agent) CreateGoverenceMessage() objects.GovernanceMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		BikeId:       uuid.Nil,
		GovernanceId: 0,
	}
}

func (bb *team5Agent) CreateForcesMessage() objects.ForcesMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.ForcesMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		AgentForces: utils.Forces{
			Pedal: 0.0,
			Brake: 0.0,
			Turning: utils.TurningDecision{
				SteerBike:     false,
				SteeringForce: 0.0,
			},
		},
	}
}

func (bb *team5Agent) CreateVoteGovernanceMessage() objects.VoteGoveranceMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteGoveranceMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *team5Agent) CreateVoteLootboxDirectionMessage() objects.VoteLootboxDirectionMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteLootboxDirectionMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *team5Agent) CreateVoteRulerMessage() objects.VoteRulerMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteRulerMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *team5Agent) CreateVotekickoutMessage() objects.VoteKickoutMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(map[uuid.UUID]int),
	}
}

func (bb *team5Agent) HandleKickoutMessage(msg objects.KickoutAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
	println("handle")
	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// kickout := msg.kickout
}

func (bb *team5Agent) HandleReputationMessage(msg objects.ReputationOfAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// reputation := msg.Reputation
}

func (bb *team5Agent) HandleJoiningMessage(msg objects.JoiningAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// bikeId := msg.BikeId
}

func (bb *team5Agent) HandleLootboxMessage(msg objects.LootboxMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// lootboxId := msg.LootboxId
}

func (bb *team5Agent) HandleGovernanceMessage(msg objects.GovernanceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// bikeId := msg.BikeId
	// governanceId := msg.GovernanceId
}

func (bb *team5Agent) HandleForcesMessage(msg objects.ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// agentForces := msg.AgentForces

}

func (bb *team5Agent) HandleVoteGovernanceMessage(msg objects.VoteGoveranceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *team5Agent) HandleVoteLootboxDirectionMessage(msg objects.VoteLootboxDirectionMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *team5Agent) HandleVoteRulerMessage(msg objects.VoteRulerMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *team5Agent) HandleVoteKickoutMessage(msg objects.VoteKickoutMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}
