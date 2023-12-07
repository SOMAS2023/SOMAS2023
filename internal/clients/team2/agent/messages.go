package agent

import (
	obj "SOMAS2023/internal/common/objects"
	"fmt"

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
	fmt.Printf("[HandleKickOffMessage] Received message from Agent %v\n", msg.AgentId)

	agentId := msg.AgentId
	if agentId == uuid.Nil {
		return
	}

	a.Modules.SocialCapital.UpdateSocialNetwork(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
	a.Modules.SocialCapital.UpdateInstitution(agentId, InstitutionEventValue_Kickoff, InstitutionEventWeight_Kickoff)
}

func (a *AgentTwo) HandleForcesMessage(msg obj.ForcesMessage) {
	agentId := msg.AgentId
	if agentId == uuid.Nil {
		return
	}

	fmt.Printf("[HandleForcesMessage] Received message from Agent %v\n", agentId)

	agentPosition := a.GetLocation()
	optimalLootbox := a.Modules.Environment.GetNearestLootboxByColor(agentId, a.GetColour())
	lootboxPosition := a.Modules.Environment.GetLootboxPos(optimalLootbox)
	optimalForces := a.Modules.Utils.GetForcesToTarget(agentPosition, lootboxPosition)
	eventValue := a.Modules.Utils.ProjectForce(optimalForces, msg.AgentForces)

	fmt.Printf("Agent Social Network Before: %v\n", a.Modules.SocialCapital.SocialNetwork)
	a.Modules.SocialCapital.UpdateSocialNetwork(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
	a.Modules.SocialCapital.UpdateInstitution(agentId, InstitutionEventWeight_Adhereance, eventValue)
	fmt.Printf("Agent Social Network After: %v\n", a.Modules.SocialCapital.SocialNetwork)
}

func (a *AgentTwo) HandleJoiningMessage(msg obj.JoiningAgentMessage) {
	fmt.Printf("[HandleJoiningMessage] Received message from Agent %v\n", msg.AgentId)

	agentId := msg.AgentId
	if agentId == uuid.Nil {
		return
	}

	a.Modules.SocialCapital.UpdateSocialNetwork(agentId, SocialEventValue_AgentSentMsg, SocialEventWeight_AgentSentMsg)
	a.Modules.SocialCapital.UpdateInstitution(agentId, InstitutionEventValue_Accepted, InstitutionEventWeight_Accepted)
}

// This function updates all the messages for that agent i.e. both sending and receiving.
// And returns the new messages from other agents to your agent
func (a *AgentTwo) GetAllMessages([]obj.IBaseBiker) []messaging.IMessage[obj.IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages and which ones to send (return)
	wantToSendMsg := true
	if wantToSendMsg {
		fmt.Printf("Agent %v is getting all messages\n", a.GetID())
		reputationMsg := a.CreateReputationMessage()
		kickoutMsg := a.CreatekickoutMessage()
		lootboxMsg := a.CreateLootboxMessage()
		joiningMsg := a.CreateJoiningMessage()
		governceMsg := a.CreateGoverenceMessage()
		forcesMsg := a.CreateForcesMessage()
		voteGoveranceMessage := a.CreateVoteGovernanceMessage()
		voteLootboxDirectionMessage := a.CreateVoteLootboxDirectionMessage()
		voteRulerMessage := a.CreateVoteRulerMessage()
		voteKickoutMessage := a.CreateVotekickoutMessage()
		return []messaging.IMessage[obj.IBaseBiker]{reputationMsg, kickoutMsg, lootboxMsg, joiningMsg, governceMsg, forcesMsg, voteGoveranceMessage, voteLootboxDirectionMessage, voteRulerMessage, voteKickoutMessage}
	}
	return []messaging.IMessage[obj.IBaseBiker]{}
}
