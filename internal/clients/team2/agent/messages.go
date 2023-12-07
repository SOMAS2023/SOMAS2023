package agent

import (
	obj "SOMAS2023/internal/common/objects"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

func (a *AgentTwo) CreateForcesMessage() obj.ForcesMessage {
	return obj.ForcesMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](a, a.GetFellowBikers()),
		AgentId:     a.GetID(),
		AgentForces: a.BaseBiker.GetForces(),
	}
}

func (a *AgentTwo) CreateKickOffMessage() obj.KickoutAgentMessage {
	agentId, _ := a.Modules.SocialCapital.GetMinimumSocialCapital()
	kickOff := false
	if agentId != a.GetID() {
		kickOff = true
	}

	return obj.KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](a, a.GetFellowBikers()),
		AgentId:     agentId,
		Kickout:     kickOff,
	}
}

func (a *AgentTwo) HandleKickOffMessage(msg obj.KickoutAgentMessage) {
	agentId := msg.AgentId

	if agentId != uuid.Nil {
		a.Modules.SocialCapital.UpdateSocialNetwork(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
		a.Modules.SocialCapital.UpdateInstitution(agentId, InstitutionEventValue_Kickoff, InstitutionEventWeight_Kickoff)
	}
}

func (a *AgentTwo) HandleForcesMessage(msg obj.ForcesMessage) {
	agentId := msg.AgentId
	agentPosition := a.GetLocation()
	optimalLootbox := a.Modules.VotedDirection
	lootboxPosition := a.Modules.Environment.GetLootboxPos(optimalLootbox)
	optimalForces := a.Modules.Utils.GetForcesToTarget(agentPosition, lootboxPosition)
	eventValue := a.Modules.Utils.ProjectForce(optimalForces, msg.AgentForces)

	a.Modules.SocialCapital.UpdateSocialNetwork(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
	a.Modules.SocialCapital.UpdateInstitution(agentId, InstitutionEventWeight_Adhereance, eventValue)
}

func (a *AgentTwo) HandleJoiningMessage(msg obj.JoiningAgentMessage) {
	agentId := msg.AgentId

	a.Modules.SocialCapital.UpdateSocialNetwork(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
	a.Modules.SocialCapital.UpdateInstitution(agentId, InstitutionEventValue_Accepted, InstitutionEventWeight_Accepted)
}
