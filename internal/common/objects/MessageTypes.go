package objects

import (
	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

type ReputationOfAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	agentId    uuid.UUID // agent who's reputation you are talking about
	reputation float64   // your agent's reputation expected from 0-1
}

type KickOffAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	agentId uuid.UUID // agent who you do/do not want to kick off
	kickOff bool      // true if you want to kick off, otherwise false
}

type JoiningAgentMessage struct {
	messaging.BaseMessage[IBaseBiker]
	agentId uuid.UUID // agent who wants to join this bike. DOESN’T MOVE YOU ONTO THAT BIKE, IS DECLARING INTENTION
	bikeId  uuid.UUID // the bike this agent wants to join
}

type LootboxMessage struct {
	messaging.BaseMessage[IBaseBiker]
	lootboxId uuid.UUID // the lootbox that agent wants
}

type GovernanceMessage struct { //"I would like to operate under this governance system" NOTE: NOT VOTING TO CHANGE GOVERNMENT
	messaging.BaseMessage[IBaseBiker]
	bikeId       uuid.UUID // the bike this agent wants to join
	governanceId int       // the governce type that this agent wants
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
