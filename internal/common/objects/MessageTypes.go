package objects

import (
	"SOMAS2023/internal/common/utils"

	voting "SOMAS2023/internal/common/voting"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

// "I have the following reputation of Agent X"
type ReputationOfAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId    uuid.UUID // agent who's reputation you are talking about
	Reputation float64   // your agent's reputation expected from 0-1
}

// "I want to kick off this agent"
type KickoutAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId uuid.UUID // agent who you do/do not want to kick off
	Kickout bool      // true if you want to kick off, otherwise false
}

// "I want to move to this bike"
type JoiningAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId uuid.UUID // agent who wants to join this bike. DOESN’T MOVE YOU ONTO THAT BIKE, IS DECLARING INTENTION
	BikeId  uuid.UUID // the bike this agent wants to join
}

// "I want to go to this lootbox next iteration"
type LootboxMessage struct {
	messaging.BaseMessage[IBaseBiker]
	LootboxId uuid.UUID // the lootbox that agent wants
}

// "I would like to operate under this governance system" NOTE: NOT VOTING TO CHANGE GOVERNMENT
type GovernanceMessage struct {
	messaging.BaseMessage[IBaseBiker]
	BikeId       uuid.UUID // the bike this agent wants to join
	GovernanceId int       // the governce type that this agent wants
}

// "I applied the following force in this iteration" or "I know Agent X applied the following force in this iteration"
type ForcesMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId     uuid.UUID    // the agent whose forces are shared
	AgentForces utils.Forces // the forces
}

// "I voted for this governance in this iteration"
type VoteGoveranceMessage struct {
	messaging.BaseMessage[IBaseBiker]
	VoteMap voting.IdVoteMap // the vote map that you voted for (if you are telling the truth)
}

// "I voted for this lootbox direction in this iteration"
type VoteLootboxDirectionMessage struct {
	messaging.BaseMessage[IBaseBiker]
	VoteMap voting.IdVoteMap // the vote map that you voted for (if you are telling the truth)
}

// "I voted for this ruler in this iteration"
type VoteRulerMessage struct {
	messaging.BaseMessage[IBaseBiker]
	VoteMap voting.IdVoteMap // the vote map that you voted for (if you are telling the truth)
}

// "I voted for kicking out this biker in this iteration"
type VoteKickoutMessage struct {
	messaging.BaseMessage[IBaseBiker]
	VoteMap map[uuid.UUID]int // the vote map that you voted for (if you are telling the truth)
}

func (msg ReputationOfAgentMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleReputationMessage(msg)
}

func (msg KickoutAgentMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleKickoutMessage(msg)
}

func (msg JoiningAgentMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleJoiningMessage(msg)
}

func (msg LootboxMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleLootboxMessage(msg)
}

func (msg GovernanceMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleGovernanceMessage(msg)
}

func (msg ForcesMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleForcesMessage(msg)
}

func (msg VoteGoveranceMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleVoteGovernanceMessage(msg)
}

func (msg VoteLootboxDirectionMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleVoteLootboxDirectionMessage(msg)
}

func (msg VoteRulerMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleVoteRulerMessage(msg)
}

func (msg VoteKickoutMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleVoteKickoutMessage(msg)
}
