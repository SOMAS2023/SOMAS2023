package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

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
	return objects.KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		Kickout:     false,
	}
}

func (bb *Agent8) CreateReputationMessage() objects.ReputationOfAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		Reputation:  1.0,
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
		LootboxId:   uuid.Nil,
	}
}

func (bb *Agent8) CreateGoverenceMessage() objects.GovernanceMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	return objects.GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		BikeId:       uuid.Nil,
		GovernanceId: 0,
	}
}

func (bb *Agent8) CreateForcesMessage() objects.ForcesMessage {
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
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *Agent8) CreateVoteRulerMessage() objects.VoteRulerMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteRulerMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *Agent8) CreateVotekickoutMessage() objects.VoteKickoutMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(map[uuid.UUID]int),
	}
}

func (bb *Agent8) HandleKickoutMessage(msg objects.KickoutAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// kickout := msg.Kickout
}

func (bb *Agent8) HandleReputationMessage(msg objects.ReputationOfAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// reputation := msg.Reputation
}

func (bb *Agent8) HandleJoiningMessage(msg objects.JoiningAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// bikeId := msg.BikeId
}

func (bb *Agent8) HandleLootboxMessage(msg objects.LootboxMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// lootboxId := msg.LootboxId
}

func (bb *Agent8) HandleGovernanceMessage(msg objects.GovernanceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// bikeId := msg.BikeId
	// governanceId := msg.GovernanceId
}

func (bb *Agent8) HandleForcesMessage(msg objects.ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// agentForces := msg.AgentForces

}

func (bb *Agent8) HandleVoteGovernanceMessage(msg objects.VoteGoveranceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *Agent8) HandleVoteLootboxDirectionMessage(msg objects.VoteLootboxDirectionMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *Agent8) HandleVoteRulerMessage(msg objects.VoteRulerMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

func (bb *Agent8) HandleVoteKickoutMessage(msg objects.VoteKickoutMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap
}

//===============================================================================================================================================================
