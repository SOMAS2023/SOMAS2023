package team6

import (
	"SOMAS2023/internal/common/objects"
	//utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

func (bb *Team6Biker) CreatekickoutMessage() objects.KickoutAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	fellowbikers := bb.GetFellowBikers()
	var kickAgent = uuid.Nil
	var lowestRep = 10.0
	for _, agents := range fellowbikers {
		agentid := agents.GetID()
		if agentid == bb.GetID() {
			continue
		}
		if bb.QueryReputation(agentid) < reputationThreshold && bb.QueryReputation(agentid) < lowestRep {
			kickAgent = agentid
			lowestRep = bb.QueryReputation(agentid)
		}
	}

	if kickAgent == uuid.Nil {
		return objects.KickoutAgentMessage{
			BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
			AgentId:     uuid.Nil,
			Kickout:     false,
		}
	} else {
		return objects.KickoutAgentMessage{
			BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
			AgentId:     kickAgent,
			Kickout:     true,
		}
	}
}

// Haven't changed yet
func (bb *Team6Biker) CreateReputationMessage() objects.ReputationOfAgentMessage {
	return objects.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     uuid.Nil,
		Reputation:  1.0,
	}
}

func (bb *Team6Biker) CreateJoiningMessage() objects.JoiningAgentMessage {
	// Send messages to the all the bikers on the bike that the agent wants to join
	bikeToJoin := bb.ChangeBike()
	return objects.JoiningAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetGameState().GetMegaBikes()[bikeToJoin].GetAgents()),
		AgentId:     bb.GetID(),
		BikeId:      bikeToJoin,
	}
}

func (bb *Team6Biker) CreateLootboxMessage() objects.LootboxMessage {
	// Propose a direction and tell all the fellow bikers
	return objects.LootboxMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		LootboxId:   bb.ProposeDirection(),
	}
}

func (bb *Team6Biker) CreateGoverenceMessage() objects.GovernanceMessage {
	// Sent messages of goverence to all fellow bikers
	governance := int(bb.DecideGovernance())
	return objects.GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		BikeId:       bb.GetBike(),
		GovernanceId: governance,
	}
}

func (bb *Team6Biker) CreateForcesMessage() objects.ForcesMessage {
	return objects.ForcesMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		AgentId:     bb.GetID(),
		AgentForces: bb.BaseBiker.GetForces(),
	}
}

func (bb *Team6Biker) CreateVoteGovernanceMessage() objects.VoteGoveranceMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteGoveranceMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     make(voting.IdVoteMap),
	}
}

func (bb *Team6Biker) CreateVoteLootboxDirectionMessage() objects.VoteLootboxDirectionMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteLootboxDirectionMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     bb.ProposedVote(),
	}
}

func (bb *Team6Biker) CreateVoteRulerMessage() objects.VoteRulerMessage {
	RulerMap := make(voting.IdVoteMap)
	for _, agent := range bb.GetFellowBikers() {
		RulerMap[agent.GetID()] = bb.GetReputation()[agent.GetID()]
		if agent.GetID() == bb.GetID() {
			RulerMap[agent.GetID()] = 1.0
		}
	}
	return objects.VoteRulerMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     RulerMap,
	}
}

func (bb *Team6Biker) CreateVotekickoutMessage() objects.VoteKickoutMessage {
	return objects.VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](bb, bb.GetFellowBikers()),
		VoteMap:     bb.VoteForKickout(),
	}
}
