package team3

import (
	obj "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

// tell fellows who is the 'worst' agent on bike, and if he is still tolerable
func (agent *SmartAgent) CreatekickoutMessage() obj.KickoutAgentMessage {
	threshold := 0.0
	agentWithLowestRep := 0
	fellows := agent.GetFellowBikers()
	if len(fellows) < 2 {
		return obj.KickoutAgentMessage{
			BaseMessage: messaging.CreateMessage[obj.IBaseBiker](agent, agent.GetFellowBikers()),
			AgentId:     uuid.Nil,
			Kickout:     false,
		}
	}
	scores := make([]float64, len(fellows))
	for idx, onBikeAgent := range fellows {
		rep := agent.reputationMap[onBikeAgent.GetID()]
		// Cognitive dimension: is same belief?
		// Contribution and Achievement
		// Forgiveness: forgive agents pedal harder recently
		// Potential
		scores[idx] = rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution + rep.energyRemain + utils.Epsilon
		threshold += scores[idx]

		if scores[idx] < scores[agentWithLowestRep] {
			agentWithLowestRep = idx
		}
	}
	threshold = threshold / float64(len(fellows)) / 2.0

	return obj.KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](agent, agent.GetFellowBikers()),
		AgentId:     fellows[agentWithLowestRep].GetID(),
		Kickout:     scores[agentWithLowestRep] < threshold,
	}
}

func (agent *SmartAgent) CreateLootboxMessage() obj.LootboxMessage {
	return obj.LootboxMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](agent, agent.GetFellowBikers()),
		LootboxId:   agent.targetLootBox,
	}
}

// tell fellows how much I pedalled this turn. Tell lies if I find there is unfairness in current bike
// How difference the lies with truth depends on how unsatisfied I felt
func (agent *SmartAgent) CreateForcesMessage() obj.ForcesMessage {
	return obj.ForcesMessage{
		BaseMessage: messaging.CreateMessage[obj.IBaseBiker](agent, agent.GetFellowBikers()),
		AgentId:     agent.GetID(),
		AgentForces: utils.Forces{
			Pedal:   agent.GetForces().Pedal / agent.satisfactionOfRecentAllocation,
			Brake:   agent.GetForces().Brake,
			Turning: agent.GetForces().Turning,
		},
	}
}

func (agent *SmartAgent) HandleKickoutMessage(msg obj.KickoutAgentMessage) {
	// discard if msg is not about kickoff someone
	if !msg.Kickout {
		return
	}

	threshold := 0.0
	fellows := agent.GetFellowBikers()
	scores := make(map[uuid.UUID]float64)
	for _, fellow := range fellows {
		rep := agent.reputationMap[fellow.GetID()]
		// Cognitive dimension: is same belief?
		// Contribution and Achievement
		// Forgiveness: forgive agents pedal harder recently
		// Potential
		scores[fellow.GetID()] = rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution + rep.energyRemain + utils.Epsilon
		threshold += scores[fellow.GetID()]
	}
	threshold = threshold / float64(len(fellows)) / 2.0

	// if not same opinion, discard msg
	score, exist := scores[msg.AgentId]
	if exist && score > threshold {
		return
	}

	sender := msg.BaseMessage.GetSender()
	rep, ok := agent.reputationMap[sender.GetID()]
	if ok {
		rep.findSameOpinion()
	}
}

func (agent *SmartAgent) HandleLootboxMessage(msg obj.LootboxMessage) {
	lootboxId := msg.LootboxId
	// if not same target, discard message
	if lootboxId != agent.targetLootBox {
		return
	}
	// if same target, increase the judgment of this agent
	sender := msg.BaseMessage.GetSender()
	rep, ok := agent.reputationMap[sender.GetID()]
	if ok {
		rep.findSameOpinion()
	}
}

// This function updates all the messages for that agent i.e. both sending and receiving.
// And returns the new messages from other agents to your agent
func (a *SmartAgent) GetAllMessages([]obj.IBaseBiker) []messaging.IMessage[obj.IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages and which ones to send (return)
	wantToSendMsg := true
	if wantToSendMsg {
		//fmt.Printf("Agent %v is getting all messages\n", a.GetID())
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
