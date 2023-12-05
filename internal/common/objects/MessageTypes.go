package objects

import (
	"SOMAS2023/internal/common/utils"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

type ReputationOfAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId    uuid.UUID // agent who's reputation you are talking about
	Reputation float64   // your agent's reputation expected from 0-1
}

type KickOffAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId uuid.UUID // agent who you do/do not want to kick off
	KickOff bool      // true if you want to kick off, otherwise false
}

type JoiningAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId uuid.UUID // agent who wants to join this bike. DOESN’T MOVE YOU ONTO THAT BIKE, IS DECLARING INTENTION
	BikeId  uuid.UUID // the bike this agent wants to join
}

type LootboxMessage struct {
	messaging.BaseMessage[IBaseBiker]
	LootboxId uuid.UUID // the lootbox that agent wants
}

type GovernanceMessage struct { //"I would like to operate under this governance system" NOTE: NOT VOTING TO CHANGE GOVERNMENT
	messaging.BaseMessage[IBaseBiker]
	BikeId       uuid.UUID // the bike this agent wants to join
	GovernanceId int       // the governce type that this agent wants
}

type ForcesMessage struct {
	messaging.BaseMessage[IBaseBiker]
	AgentId     uuid.UUID    // the agent whose forces are shared
	AgentForces utils.Forces // the forces
}

func (msg ReputationOfAgentMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleReputationMessage(msg)
}

func (msg KickOffAgentMessage) InvokeMessageHandler(agent IBaseBiker) {
	agent.HandleKickOffMessage(msg)
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
