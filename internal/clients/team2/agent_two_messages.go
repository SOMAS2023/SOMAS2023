package team2

import (
	obj "SOMAS2023/internal/common/objects"
	"math"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

func (a *AgentTwo) CreateForcesMessage() obj.ForcesMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents

	return obj.ForcesMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](a, a.GetFellowBikers()),
		AgentId:     a.GetID(),
		AgentForces: a.forces,
	}
}

func (a *AgentTwo) CreateKickOffMessage() obj.KickOffAgentMessage {
	kickOff := false
	minAgentId := uuid.Nil
	minCapital := math.MaxFloat64
	for agentId, value := range a.SocialCapital {
		if value < minCapital {
			kickOff = true
			minCapital = value
			minAgentId = agentId
		}
	}

	return obj.KickOffAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](a, a.GetFellowBikers()),
		AgentId:     minAgentId,
		KickOff:     kickOff,
	}
}

func (a *AgentTwo) HandleKickOffMessage(msg obj.KickOffAgentMessage) {
	agentId := msg.AgentId
	if agentId != uuid.Nil {
		a.UpdateSocNetAgent(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
		a.updateInstitution(agentId, InstitutionEventWeight_KickedOut, InstitutionKickoffEventValue)
	}
}

func (a *AgentTwo) HandleForcesMessage(msg obj.ForcesMessage) {

	agentId := msg.AgentId
	agentForces := msg.AgentForces
	optimalLootbox := a.votedDirection
	optimalForces := a.GetVotedLootboxForces(optimalLootbox)

	EventValue := a.RuleAdhereanceValue(agentId, optimalForces, agentForces)

	a.UpdateSocNetAgent(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
	a.updateInstitution(agentId, InstitutionEventWeight_Adhereance, EventValue)

}

func (a *AgentTwo) HandleJoiningMessage(msg obj.JoiningAgentMessage) {

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// bikeId := msg.BikeId
	agentId := msg.AgentId
	a.UpdateSocNetAgent(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
	a.updateInstitution(agentId, InstitutionEventWeight_Accepted, InstitutionAcceptedEventValue)
}
