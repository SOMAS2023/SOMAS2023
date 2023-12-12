package team6

import (
	"SOMAS2023/internal/common/objects"
	"fmt"
	"math"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
)

const TrustThreshold = 0.3

func (bb *Team6Biker) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
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

func (bb *Team6Biker) HandleVoteRulerMessage(msg objects.VoteRulerMessage) {

}

func (bb *Team6Biker) HandleKickoutMessage(msg objects.KickoutAgentMessage) {
	sender := msg.BaseMessage.GetSender()
	agentid := msg.AgentId
	kickout := msg.Kickout
	if agentid == bb.GetID() && kickout == true {
		bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())-0.3, 0))
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("---------HandleKickoutMessage", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}

func (bb *Team6Biker) HandleJoiningMessage(msg objects.JoiningAgentMessage) {
	sender := msg.BaseMessage.GetSender()
	senderReputation := bb.QueryReputation(sender.GetID())
	fmt.Print("HandleJoiningMessage part is fine")
	if senderReputation < TrustThreshold {
		bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())-0.05, 1))
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("---------HandleJoiningMessage", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}

func (bb *Team6Biker) HandleLootboxMessage(msg objects.LootboxMessage) {
	sender := msg.BaseMessage.GetSender()
	senderColour := sender.GetColour()
	if senderColour == bb.GetColour() {
		bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+0.1, 1))
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("---------HandleLootboxMessage", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}

func (bb *Team6Biker) HandleGovernanceMessage(msg objects.GovernanceMessage) {
	sender := msg.BaseMessage.GetSender()
	senderGovernance := msg.GovernanceId
	if senderGovernance == int(bb.DecideGovernance()) {
		bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+0.1, 1))
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("-------HandleGovernanceMessage", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}

func (bb *Team6Biker) HandleForcesMessage(msg objects.ForcesMessage) {
	sender := msg.BaseMessage.GetSender()
	senderForces := msg.AgentForces
	if senderForces.Brake > 0 {
		bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())-senderForces.Brake/2, 0))
	}
	if senderForces.Pedal > 0 {
		bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+senderForces.Pedal/2, 1))
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("-----------HandleForcesMessage", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}

func (bb *Team6Biker) HandleVoteGovernanceMessage(msg objects.VoteGoveranceMessage) {

	sender := msg.BaseMessage.GetSender()
	voteMap := msg.VoteMap

	//type IdVoteMap map[uuid.UUID]float64

	for agent, score := range voteMap {
		if score == 1 {
			fmt.Print(agent, sender)
		}
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("-----------HandleVoteGovernanceMessage", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}

func (bb *Team6Biker) HandleVoteLootboxDirectionMessage(msg objects.VoteLootboxDirectionMessage) {

	sender := msg.BaseMessage.GetSender()
	//VoteMap voting.IdVoteMap // the vote map that you voted for (if you are telling the truth)
	voteMap := msg.VoteMap
	soughtloot := bb.ProposeDirection()
	for lootbox, score := range voteMap {
		if score == 1 && soughtloot == lootbox {
			bb.SetReputation(sender.GetID(), math.Min(bb.QueryReputation(sender.GetID())+0.35, 1))
		}
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("-----------HandleVoteLootboxDirectionMessageFromAgent", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}

func (bb *Team6Biker) HandleVoteKickoutMessage(msg objects.VoteKickoutMessage) {
	sender := msg.BaseMessage.GetSender()
	voteMap := msg.VoteMap
	for agent, score := range voteMap {
		if agent == bb.GetID() && score == 1 {
			bb.SetReputation(sender.GetID(), math.Max(bb.QueryReputation(sender.GetID())-0.3, 0))
		}
	}
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
	fmt.Println("-----------HandleVoteKickoutMessage", sender.GetID())
	fmt.Println("---------------------------------------")
	fmt.Println("---------------------------------------")
}
