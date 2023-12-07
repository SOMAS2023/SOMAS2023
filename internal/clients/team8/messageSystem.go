package team8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/voting"
	"math"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Message System <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// This function updates all the messages for that agent i.e. both sending and receiving.
// And returns the new messages from other agents to your agent
func (bb *Agent8) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages and which ones to send (return)
	wantToSendMsg := false
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

func (bb *Agent8) CreatekickoutMessage() objects.KickoutAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	bikeID := bb.GetBike()
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	kickAgent := uuid.Nil
	protectAgent := uuid.Nil
	initalreputation := 0.0
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if bb.QueryReputation(agentID) < initalreputation {
			kickAgent = agentID
			initalreputation = bb.QueryReputation(agentID)
		}
	}
	if kickAgent == uuid.Nil {
		for _, agent := range fellowBikers {
			agentID := agent.GetID()
			if bb.QueryReputation(agentID) > initalreputation {
				protectAgent = agentID
				initalreputation = bb.QueryReputation(agentID)
			}
		}
		return objects.KickoutAgentMessage{
			BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
			AgentId:     protectAgent,
			Kickout:     false,
		}
	}
	return objects.KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     kickAgent,
		Kickout:     true,
	}
}

func (bb *Agent8) CreateReputationMessage() objects.ReputationOfAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	bikeID := bb.GetBike()
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	bestAgent := uuid.Nil
	reputation := -1.0
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if bb.QueryReputation(agentID) > reputation {
			bestAgent = agentID
			reputation = bb.QueryReputation(agentID)
		}
	}
	return objects.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     bestAgent,
		Reputation:  reputation,
	}
}

func (bb *Agent8) CreateJoiningMessage() objects.JoiningAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.JoiningAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		BikeId:      uuid.Nil,
	}
}
func (bb *Agent8) CreateLootboxMessage() objects.LootboxMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.LootboxMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		LootboxId:   bb.ProposeDirection(),
	}
}

func (bb *Agent8) CreateGoverenceMessage() objects.GovernanceMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		BikeId:       bb.ChangeBike(),
		GovernanceId: int(bb.DecideGovernance()),
	}
}

func (bb *Agent8) CreateForcesMessage() objects.ForcesMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	bb.DecideForce(bb.ProposeDirection())
	return objects.ForcesMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     bb.GetID(),
		AgentForces: bb.GetForces(),
	}
}

func (bb *Agent8) CreateVoteGovernanceMessage() objects.VoteGoveranceMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteGoveranceMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *Agent8) CreateVoteLootboxDirectionMessage() objects.VoteLootboxDirectionMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteLootboxDirectionMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     bb.overallLootboxPreferences.GetVotes(),
	}
}

func (bb *Agent8) CreateVoteRulerMessage() objects.VoteRulerMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	reputationMap := bb.GetReputation()
	bikeID := bb.GetBike()
	voteRulerMap := make(voting.IdVoteMap)
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		voteRulerMap[agentID] = reputationMap[agentID]
		if agentID == bb.GetID() {
			voteRulerMap[agentID] = 1.0
		}
	}
	voteRulerMap = softmax(voteRulerMap)
	return objects.VoteRulerMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     voteRulerMap,
	}
}

func (bb *Agent8) CreateVotekickoutMessage() objects.VoteKickoutMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	voteResults := make(map[uuid.UUID]int)
	bikeID := bb.GetBike()
	kickAgentID := uuid.Nil
	lowestReputation := 0.0
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		voteResults[agentID] = 0
		if bb.QueryReputation(agentID) < lowestReputation {
			kickAgentID = agentID
			lowestReputation = bb.QueryReputation(agentID)
		}
	}
	voteResults[kickAgentID] = 1
	return objects.VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     voteResults,
	}
}

func (bb *Agent8) HandleKickoutMessage(msg objects.KickoutAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	agentId := msg.AgentId
	kickout := msg.Kickout
	senderReputation := bb.QueryReputation(sender.GetID())
	kickAgentReputation := bb.QueryReputation(agentId)
	if kickout {
		if senderReputation > 0.5 {
			if kickAgentReputation < 0.2 {
				bb.SetReputation(agentId, math.Max(kickAgentReputation-0.05, -1))
			}
		}
	} else {
		if senderReputation > 0.5 {
			if kickAgentReputation > 0.2 {
				bb.SetReputation(agentId, math.Min(kickAgentReputation+0.05, 1))
			}
		}
	}
}

func (bb *Agent8) HandleReputationMessage(msg objects.ReputationOfAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	senderReputation := bb.QueryReputation(sender.GetID())
	agentId := msg.AgentId
	reputation := msg.Reputation

	if senderReputation > 0.5 {
		bb.SetReputation(agentId, math.Min(1, math.Max(-1, bb.QueryReputation(agentId)+
			0.1*(reputation-bb.GetAverageReputation(sender))/bb.GetAverageReputation(sender))))
	}
}

func (bb *Agent8) HandleJoiningMessage(msg objects.JoiningAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
	//ignore
}

func (bb *Agent8) HandleLootboxMessage(msg objects.LootboxMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	lootboxId := msg.LootboxId
	senderReputation := bb.QueryReputation(sender.GetID())
	if senderReputation >= 0.0 {
		if lootboxId == bb.ProposeDirection() {
			bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+0.05, 1))
		} else {
			bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())-0.05, -1))
		}
	}
}

func (bb *Agent8) HandleGovernanceMessage(msg objects.GovernanceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	bikeId := msg.BikeId
	governanceId := msg.GovernanceId
	senderReputation := bb.QueryReputation(sender.GetID())
	if senderReputation >= 0.0 {
		if bikeId == bb.GetBike() {
			if governanceId == int(bb.DecideGovernance()) {
				bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+0.05, 1))
			} else {
				bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())+0.05, -1))
			}
		}
	}
}

func (bb *Agent8) HandleForcesMessage(msg objects.ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	agentId := msg.AgentId
	agentForces := msg.AgentForces
	senderReputation := bb.QueryReputation(sender.GetID())
	if senderReputation >= 0.0 {
		if agentForces.Turning == bb.GetForces().Turning {
			bb.SetReputation(agentId, math.Min(bb.QueryReputation(agentId)+0.05*agentForces.Pedal, 1))
		} else {
			bb.SetReputation(agentId, math.Max(bb.QueryReputation(agentId)-0.05*agentForces.Pedal, -1))
		}
	}
}

func (bb *Agent8) HandleVoteGovernanceMessage(msg objects.VoteGoveranceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
	//Todo
	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap

}

func (bb *Agent8) HandleVoteLootboxDirectionMessage(msg objects.VoteLootboxDirectionMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
	//Todo
	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *Agent8) HandleVoteRulerMessage(msg objects.VoteRulerMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
	sender := msg.BaseMessage.GetSender()
	voteMap := msg.VoteMap
	firstRuler := uuid.Nil
	bestScore := 0.0
	chooseMe := false
	for agent, score := range voteMap {
		if agent == bb.GetID() {
			chooseMe = true
		}
		if score > bestScore {
			bestScore = score
			firstRuler = agent
		}
	}

	if firstRuler == bb.GetID() {
		bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+0.1, 1))
	}

	if chooseMe {
		bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+0.05, 1))
	} else {
		bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())-0.1, -1))
	}
}

func (bb *Agent8) HandleVoteKickoutMessage(msg objects.VoteKickoutMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	voteMap := msg.VoteMap
	for agent, score := range voteMap {
		if agent == bb.GetID() && score == 1 {
			bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())-0.2, -1))
		}
	}
}

//===============================================================================================================================================================
