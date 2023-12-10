package team1

import (
	obj "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"fmt"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

// -----------------MESSAGING FUNCTIONS------------------

// Handle a message received from anyone, ensuring they are trustworthy and come from the right place (e.g. our bike)
func (bb *Biker1) VerifySender(sender obj.IBaseBiker) bool {
	// check if sender is on our bike
	if sender.GetBike() == bb.GetBike() {
		// check if sender is trustworthy
		if bb.opinions[sender.GetID()].trust > trustThreshold && bb.opinions[sender.GetID()].opinion > 0.5 {
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
		if msg.AgentId != uuid.Nil {
			if msg.Kickout {
				if bb.opinions[msg.AgentId].opinion > 0.5 {
					penalty := 0.9
					bb.UpdateOpinion(sender.GetID(), penalty)
				} else {
					sameOpinionreward := 1.1
					bb.UpdateOpinion(sender.GetID(), sameOpinionreward)
				}
			} else {
				if bb.opinions[msg.AgentId].opinion > 0.5 {
					sameOpinionreward := 1.1
					bb.UpdateOpinion(sender.GetID(), sameOpinionreward)
				} else {
					penalty := 0.9
					bb.UpdateOpinion(sender.GetID(), penalty)
				}
			}
		}

	}

}

// Agent receives a reputation of another agent
func (bb *Biker1) HandleReputationMessage(msg obj.ReputationOfAgentMessage) {
	sender := msg.GetSender()
	verified := bb.VerifySender(sender)

	if verified {
		// TODO: SOME FORMULA TO UPDATE OPINION BASED ON REPUTATION given
		if msg.AgentId != uuid.Nil {
			// Retrieve the struct from the map
			opinion, ok := bb.opinions[msg.AgentId]
			if ok {
				// Update the field
				opinion.trust += msg.Reputation * reputationScaling
				bb.opinions[msg.AgentId] = opinion
			}
			currentReputation := bb.GetReputation()[msg.AgentId] + msg.Reputation
			bb.SetReputation(msg.AgentId, currentReputation*reputationScaling)
		}
	}
	// ask fellow bikers what their reputation of incoming biker is..
}

// Agent receives a message from another agent to join
func (bb *Biker1) HandleJoiningMessage(msg obj.JoiningAgentMessage) {
	sender := msg.GetSender()
	// different from Verify sender since they are not on our bike
	if bb.opinions[sender.GetID()].trust > trustThreshold && bb.opinions[sender.GetID()].opinion > 0.5 {
		// check if sender is on our bike
		if msg.AgentId != uuid.Nil {
			agentToJoin := bb.GetAgentFromId(msg.AgentId)
			if agentToJoin.GetColour() == bb.GetColour() {
				sameColourReward := 1.1
				bb.UpdateOpinion(sender.GetID(), sameColourReward)
			}
		}

	}

}

// Agent receives a message from another agent say what lootbox they want to go to
func (bb *Biker1) HandleLootboxMessage(msg obj.LootboxMessage) {
	sender := msg.GetSender()
	verified := bb.VerifySender(sender)
	if verified {
		if msg.LootboxId != uuid.Nil {
			if sender.GetColour() == bb.GetColour() {
				sameColourReward := 1.2
				bb.UpdateOpinion(sender.GetID(), sameColourReward)
			}
		}
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
			// Lower opinion if they are not pedalling
			if msg.AgentForces.Pedal == 0 {
				bb.UpdateOpinion(sender.GetID(), 0.9)
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
	opinions, ok := bb.opinions[bb.GetID()]
	reputation := opinions.opinion
	if !ok {
		reputation = 0.0
	} else {
		reputation = bb.opinions[bb.GetID()].opinion
	}

	return obj.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](bb, bb.GetTrustedRecepients()),
		AgentId:     uuid.Nil,
		Reputation:  reputation,
	}
}

func (bb *Biker1) CreateJoiningMessage() obj.JoiningAgentMessage {
	// Tell the truth (for now)
	// receipients = fellowBikers
	biketoJoin := bb.PickBestBike()
	fmt.Printf(biketoJoin.String())
	gs := bb.GetGameState()
	joiningBike := gs.GetMegaBikes()[biketoJoin]
	fmt.Printf("Joining bike: %v", joiningBike)
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
	sendKickMessage = false
	sendReputationMessage = false
	sendJoiningMessage = false
	sendLootboxMessage = false
	sendGovernanceMessage = false

	// TODO: add logic to decide which messages to send and when
	if bb.GetBike() == uuid.Nil && !bb.GetBikeStatus(){
		sendGovernanceMessage = true
		sendJoiningMessage = false
	} else if bb.GetBike() == uuid.Nil {
		fmt.Printf("Bike is nil\n")
		sendJoiningMessage = true
	} else {
		for _, agent := range bb.GetFellowBikers() {
			if bb.opinions[agent.GetID()].opinion < kickThreshold {
				sendKickMessage = true
			}

			if (bb.opinions[agent.GetID()].trust > trustThreshold) && (bb.opinions[agent.GetID()].opinion > 0.5) {
				sendGovernanceMessage = true
				sendLootboxMessage = true
				// Never send reputation message			}
			}
		}
	}

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
