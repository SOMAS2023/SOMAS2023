package team3

import (
	obj "SOMAS2023/internal/common/objects"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
)

func (a *SmartAgent) CreateForcesMessage() obj.ForcesMessage {
	return obj.ForcesMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](a, a.GetFellowBikers()),
		AgentId:     a.GetID(),
		AgentForces: a.GetForces(),
	}
}

// This function updates all the messages for that agent i.e. both sending and receiving.
// And returns the new messages from other agents to your agent
func (a *SmartAgent) GetAllMessages([]obj.IBaseBiker) []messaging.IMessage[obj.IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages and which ones to send (return)
	wantToSendMsg := true
	if wantToSendMsg {
		// fmt.Printf("Agent %v is getting all messages\n", a.GetID())
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
