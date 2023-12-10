package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	"math"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
)

// func (t5 *team5Agent) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
// 	var messages []messaging.IMessage[objects.IBaseBiker]

// 	// send message to all other agents on the bike containing our reputation value on them, expecting them to send back their reputation value on us
// 	for _, agent := range t5.GetFellowBikers() {
// 		if agent.GetID() != t5.GetID() && t5.QueryReputation(agent.GetID()) >= 0.6 {
// 			repMsg := t5.CreateReputationMessage(agent)
// 			messages = append(messages, repMsg)
// 		}
// 	}

// 	// send message to all other agents on the bike containing our forces, expecting them to send back their forces
// 	forcesMsg := t5.CreateForcesMessage()

// 	messages = append(messages, forcesMsg)

// 	return messages
// }

//-------------------- Create Messages ---------------------------------------------------------------------

func (t5 *team5Agent) CreateForcesMessage() objects.ForcesMessage {
	// Send our own forces
	return objects.ForcesMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](t5, t5.GetFellowBikers()),
		AgentId:     t5.GetID(),
		AgentForces: t5.GetForces(),
	}
}

func (t5 *team5Agent) CreateReputationMessage(agent objects.IBaseBiker) objects.ReputationOfAgentMessage {
	// Praise the agent that has a high reputation according to us
	return objects.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](t5, t5.GetFellowBikers()),
		AgentId:     agent.GetID(),
		Reputation:  t5.QueryReputation(agent.GetID()),
	}
}

//-------------------- Handle Messages ---------------------------------------------------------------------

func (t5 *team5Agent) HandleReputationMessage(msg objects.ReputationOfAgentMessage) {
	senderID := msg.BaseMessage.GetSender().GetID()
	agentId := msg.AgentId
	reputation := msg.Reputation

	// If the agent ID is the agent itself, store the reputation value
	if agentId == t5.GetID() {
		t5.otherBikerRep[senderID] = reputation
	}
}

func (t5 *team5Agent) HandleForcesMessage(msg objects.ForcesMessage) {
	senderID := msg.BaseMessage.GetSender().GetID()
	agentId := msg.AgentId
	forces := msg.AgentForces

	c := 0.05

	// If the sender gives his own forces and is on our bike, store the forces
	for _, agent := range t5.GetFellowBikers() {
		if senderID == agentId && agent.GetID() == senderID {
			t5.otherBikerForces[senderID] = forces

			// If the sender is slacking, lower the reputation value
			if forces.Pedal < 0.3 {
				t5.SetReputation(senderID, math.Max(t5.QueryReputation(senderID)-c, -1))
			}
		}
	}
}

func (t5 *team5Agent) HandleLootboxMessage(msg objects.LootboxMessage) {
	senderID := msg.BaseMessage.GetSender().GetID()
	preferredLootbox := msg.LootboxId

	ourPreference := t5.ProposeDirection()

	c := 0.05

	if preferredLootbox == ourPreference {
		t5.SetReputation(senderID, math.Min(t5.QueryReputation(senderID)+c, 1))
	} else {
		t5.SetReputation(senderID, math.Max(t5.QueryReputation(senderID)-c, -1))
	}
}
