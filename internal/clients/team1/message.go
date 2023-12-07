package team1

import (
	obj "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

// -----------------MESSAGING FUNCTIONS------------------

// Handle a message received from anyone, ensuring they are trustworthy and come from the right place (e.g. our bike)
func (bb *Biker1) VerifySender(sender obj.IBaseBiker) bool {
	// check if sender is on our bike
	if sender.GetBike() == bb.GetBike() {
		// check if sender is trustworthy
		if bb.opinions[sender.GetID()].trust > trustThreshold {
			return true
		}
	}
	return false
}

// Agent receives a who to kick off message
func (bb *Biker1) HandleKickOffMessage(msg obj.KickoutAgentMessage) {
	sender := msg.GetSender()
	verified := bb.VerifySender(sender)
	if verified {
		// slightly penalise view of person who sent message
		penalty := 0.9
		bb.UpdateOpinion(sender.GetID(), penalty)
	}

}

// Agent receives a reputation of another agent
func (bb *Biker1) HandleReputationMessage(msg obj.ReputationOfAgentMessage) {
	sender := msg.GetSender()
	verified := bb.VerifySender(sender)
	if verified {
		// TODO: SOME FORMULA TO UPDATE OPINION BASED ON REPUTATION given
	}
}

// Agent receives a message from another agent to join
func (bb *Biker1) HandleJoiningMessage(msg obj.JoiningAgentMessage) {
	sender := msg.GetSender()
	// check if sender is trustworthy
	if bb.opinions[sender.GetID()].trust > trustThreshold {
		// TODO: some update on opinon maybe???
	}

}

// Agent receives a message from another agent say what lootbox they want to go to
func (bb *Biker1) HandleLootboxMessage(msg obj.LootboxMessage) {
	sender := msg.GetSender()
	verified := bb.VerifySender(sender)
	if verified {
		// TODO: some update on lootbox decision maybe??
	}
}

// Agent receives a message from another agent saying what Governance they want
func (bb *Biker1) HandleGovernanceMessage(msg obj.GovernanceMessage) {
	sender := msg.GetSender()
	verified := bb.VerifySender(sender)
	if verified {
		// TODO: some update on governance decision maybe??
	}
}

// HandleForcesMessage
func (bb *Biker1) HandleForcesMessage(msg obj.ForcesMessage) {
	sender := msg.GetSender()
	verified := bb.VerifySender(sender)
	if verified {
		//if we are dictator and the pedal force is 0, or they are braking, or they are turning differently, add them to the kick list
		bikeID := bb.GetBike()
		gs := bb.GetGameState()
		bike := gs.GetMegaBikes()[bikeID]
		if bb.GetID() == bike.GetRuler() && bike.GetGovernance() == utils.Dictatorship {
			if msg.AgentForces.Brake > 0 || msg.AgentForces.Turning.SteerBike {
				//set our opinion of them to 0, should be kicked in next loop
				bb.UpdateOpinion(sender.GetID(), 0)
			}
			if msg.AgentForces.Pedal == 0 {
				bb.UpdateOpinion(sender.GetID(), bb.opinions[sender.GetID()].opinion*0.9)
			}
		}
		return
	}
}

func (bb *Biker1) GetTrustedRecepients() []obj.IBaseBiker {
	fellowBikers := bb.GetFellowBikers()
	var trustedRecepients []obj.IBaseBiker
	for _, agent := range fellowBikers {
		if bb.opinions[agent.GetID()].trust > trustThreshold {
			trustedRecepients = append(trustedRecepients, agent)
		}
	}
	return trustedRecepients
}

// CREATING MESSAGES
func (bb *Biker1) CreateKickOffMessage() obj.KickoutAgentMessage {
	// Receipients = fellowBikers
	agentToKick := bb.lowestOpinionKick()
	var kickDecision bool
	// send kick off message if we have a low opinion of someone
	if agentToKick != uuid.Nil {
		kickDecision = true
	} else {
		kickDecision = false
	}

	return obj.KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](bb, bb.GetTrustedRecepients()),
		AgentId:     agentToKick,
		Kickout:     kickDecision,
	}
}

func (bb *Biker1) CreateReputationMessage() obj.ReputationOfAgentMessage {
	// Tell the truth (for now)
	// TODO: receipients = fellowBikers that we trust?
	return obj.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](bb, bb.GetTrustedRecepients()),
		AgentId:     uuid.Nil,
		Reputation:  1.0,
	}
}

func (bb *Biker1) CreateJoiningMessage() obj.JoiningAgentMessage {
	// Tell the truth (for now)
	// receipients = fellowBikers
	biketoJoin := bb.ChangeBike()
	gs := bb.GetGameState()
	joiningBike := gs.GetMegaBikes()[biketoJoin]
	return obj.JoiningAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](bb, joiningBike.GetAgents()),
		AgentId:     bb.GetID(),
		BikeId:      biketoJoin,
	}
}
func (bb *Biker1) CreateLootboxMessage() obj.LootboxMessage {
	// Tell the truth (for now)
	// receipients = fellowBikers
	chosenLootbox := bb.ProposeDirection()
	return obj.LootboxMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](bb, bb.GetTrustedRecepients()),
		LootboxId:   chosenLootbox,
	}
}

func (bb *Biker1) CreateGoverenceMessage() obj.GovernanceMessage {
	// Tell the truth (using same logic as deciding governance for voting) (for now)
	// receipients = fellowBikers
	chosenGovernance := bb.DecideGovernance()
	// convert to int for now
	// send governance message to all agents (as not on a  bike yet)
	// todo: improve it so that it only sends it to trusted agents (among all the other agents in the game)
	agentMap := bb.GetGameState().GetAgents()
	allAgents := make([]obj.IBaseBiker, len(agentMap))
	i := 0
	for _, agent := range agentMap {
		allAgents[i] = agent
		i++
	}
	return obj.GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[obj.IBaseBiker](bb, allAgents),
		BikeId:       bb.GetBike(),
		GovernanceId: int(chosenGovernance),
	}
}

// Agent sending messages to other agents
func (bb *Biker1) GetAllMessages([]obj.IBaseBiker) []messaging.IMessage[obj.IBaseBiker] {
	var sendKickMessage, sendReputationMessage, sendJoiningMessage, sendLootboxMessage, sendGovernanceMessage bool

	// TODO: add logic to decide which messages to send and when

	var messageList []messaging.IMessage[obj.IBaseBiker]
	if sendKickMessage {
		messageList = append(messageList, bb.CreateKickOffMessage())
	}
	if sendReputationMessage {
		messageList = append(messageList, bb.CreateReputationMessage())
	}
	if sendJoiningMessage {
		messageList = append(messageList, bb.CreateJoiningMessage())
	}
	if sendLootboxMessage {
		messageList = append(messageList, bb.CreateLootboxMessage())

	}
	if sendGovernanceMessage {
		messageList = append(messageList, bb.CreateGoverenceMessage())

	}
	return messageList
}

// -----------------END MESSAGING FUNCTIONS------------------
