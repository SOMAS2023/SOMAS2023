package objects

import (
	obj "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

// Create new
type ExtendedBaseBiker struct {
	*obj.BaseBiker       // BaseBiker inherits functions from BaseAgent such as GetID(), GetAllMessages() and UpdateAgentInternalState()
	OtherBiker           obj.IBaseBiker
	OtherBikerReputation float64
	OtherBikerForces     []utils.Forces
}

// Produce new IExtendedBaseBiker
func NewExtendedBaseBiker(agentId uuid.UUID) *ExtendedBaseBiker {
	return &ExtendedBaseBiker{
		BaseBiker: obj.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()),
	}
}

func (ebb *ExtendedBaseBiker) HandleReputationMessage(msg obj.ReputationOfAgentMessage) {
	agentId := msg.AgentId
	reputation := msg.Reputation

	if agentId == ebb.OtherBiker.GetID() {
		ebb.OtherBikerReputation += reputation
	}
}

func (ebb *ExtendedBaseBiker) HandleForcesMessage(msg obj.ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	agentId := msg.AgentId
	agentForces := msg.AgentForces

	if agentId == ebb.OtherBiker.GetID() {
		ebb.OtherBikerForces = append(ebb.OtherBikerForces, agentForces)
	}

}

// Send messages to server which will be sent to the reciepients of the messages
func (ebb *ExtendedBaseBiker) GetAllMessages([]obj.IBaseBiker) []messaging.IMessage[obj.IBaseBiker] {
	reputationMsg := ebb.CreateReputationMessage()
	forcesMsg := ebb.CreateForcesMessage()
	return []messaging.IMessage[obj.IBaseBiker]{reputationMsg, forcesMsg}
}

func (ebb *ExtendedBaseBiker) CreateReputationMessage() obj.ReputationOfAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	var updatedAgents []obj.IBaseBiker
	updatedAgents = append(updatedAgents, ebb.OtherBiker)
	return obj.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](ebb, updatedAgents),
		AgentId:     ebb.OtherBiker.GetID(),
		Reputation:  0.1,
	}
}

func (ebb *ExtendedBaseBiker) CreateForcesMessage() obj.ForcesMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	var updatedAgents []obj.IBaseBiker
	updatedAgents = append(updatedAgents, ebb.OtherBiker)
	return obj.ForcesMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](ebb, updatedAgents),
		AgentId:     ebb.OtherBiker.GetID(),
		AgentForces: ebb.GetForces(),
	}
}

func TestBaseBikerMessaging(t *testing.T) {
	biker1 := NewExtendedBaseBiker(uuid.New())
	biker2 := NewExtendedBaseBiker(uuid.New())
	var bikers []obj.IBaseBiker
	bikers = append(bikers, biker1)

	biker1.OtherBiker = biker2
	biker2.OtherBiker = biker1

	for i := 0; i < 5; i++ {
		msgs1 := biker1.GetAllMessages(bikers)
		msgs2 := biker2.GetAllMessages(bikers)
		msgs1[0].InvokeMessageHandler(biker1)
		msgs2[0].InvokeMessageHandler(biker2)
	}

	// Assert that biker1.OtherBikerReputation and biker2.OtherBikerReputation are both equal to 5
	if biker1.OtherBikerReputation != 0.5 || biker2.OtherBikerReputation != 0.5 {
		t.Errorf("Expected both biker1 and biker2 reputations to be 0.5, but got biker1=%.2f and biker2=%.2f", biker1.OtherBikerReputation, biker2.OtherBikerReputation)
	}
}

func TestBaseBikerForcesMessaging(t *testing.T) {
	biker1 := NewExtendedBaseBiker(uuid.New())
	biker2 := NewExtendedBaseBiker(uuid.New())
	var bikers []obj.IBaseBiker
	bikers = append(bikers, biker1)

	biker1.OtherBiker = biker2
	biker2.OtherBiker = biker1

	for i := 0; i < 5; i++ {
		msgs1 := biker1.GetAllMessages(bikers)
		msgs2 := biker2.GetAllMessages(bikers)
		// fmt.Println("iteration ", i)
		// fmt.Println("msg biker 1: ", msgs1)
		// fmt.Println("msg biker 2: ", msgs2)
		msgs1[0].InvokeMessageHandler(biker1)
		msgs2[0].InvokeMessageHandler(biker2)

		msgs1[1].InvokeMessageHandler(biker1)
		msgs2[1].InvokeMessageHandler(biker2)

	}

	// Assert that biker1.OtherBikerReputation and biker2.OtherBikerReputation are both equal to 5
	if len(biker1.OtherBikerForces) != 5 || len(biker2.OtherBikerForces) != 5 {
		t.Errorf("Expected both biker1 and biker2 to store 5 past forces, but got biker1=%d and biker2=%d", len(biker1.OtherBikerForces), len(biker2.OtherBikerForces))
	}
}
