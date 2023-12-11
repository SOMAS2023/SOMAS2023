package team4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/voting"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

// handle kick out message
const ifKickOutMsgValue = 0.2

// handle reputation
const reputationThreshold = 0.65
const ifReputationIsLow = 0.1

// handle loot box message
const likeLootBox = 0.2
const dislikeLootBox = 0.05

// handle governance
const sameColorDemorcracy = 0.1
const sameColorLeader = 0.05
const sameColorDictator = 0.05

const differentColorDemorcracy = 0.1
const differntColorLeader = 0.05
const differentColorDictatorship = 0.05

// hanle force message
const ifForcesTooLow = 0.4

// voteRulerMsg
const ifRulerSameColor = 0.1
const ifRulerDifferentColor = 0.1
const ifWeAreRuler = 0.2

// handle kickout message
const ifKickOutUs = 0.5

func (agent *BaselineAgent) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages and which ones to send (return)
	wantToSendMsg := false
	if wantToSendMsg {
		reputationMsg := agent.CreateReputationMessage()
		kickoutMsg := agent.CreatekickoutMessage()
		lootboxMsg := agent.CreateLootboxMessage()
		joiningMsg := agent.CreateJoiningMessage()
		governceMsg := agent.CreateGoverenceMessage()
		forcesMsg := agent.CreateForcesMessage()
		voteGoveranceMessage := agent.CreateVoteGovernanceMessage()
		voteLootboxDirectionMessage := agent.CreateVoteLootboxDirectionMessage()
		voteRulerMessage := agent.CreateVoteRulerMessage()
		voteKickoutMessage := agent.CreateVotekickoutMessage()
		return []messaging.IMessage[objects.IBaseBiker]{reputationMsg, kickoutMsg, lootboxMsg, joiningMsg, governceMsg, forcesMsg, voteGoveranceMessage, voteLootboxDirectionMessage, voteRulerMessage, voteKickoutMessage}
	}
	return []messaging.IMessage[objects.IBaseBiker]{}
}

func (agent *BaselineAgent) CreatekickoutMessage() objects.KickoutAgentMessage {
	// if reputation and honesty is both below average then kickout, otherwise we're friendly
	bikeID := agent.GetBike()
	fellowBikers := agent.GetGameState().GetMegaBikes()[bikeID].GetAgents()

	totalReputation, totalHonesty := 0.0, 0.0
	for _, biker := range fellowBikers {
		totalReputation += agent.QueryReputation(biker.GetID())
		totalHonesty += agent.QueryHonesty(biker.GetID())
	}
	avgReputation := totalReputation / float64(len(fellowBikers))
	avgHonesty := totalHonesty / float64(len(fellowBikers))

	var kickAgentID uuid.UUID
	for _, biker := range fellowBikers {
		agentID := biker.GetID()
		if agentID == agent.GetID() {
			continue
		}
		if agent.QueryReputation(agentID) < avgReputation && agent.QueryHonesty(agentID) < avgHonesty {
			kickAgentID = agentID
			break
		}
	}

	if kickAgentID != uuid.Nil {
		return objects.KickoutAgentMessage{
			BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, fellowBikers),
			AgentId:     kickAgentID,
			Kickout:     true,
		}
	} else {
		return objects.KickoutAgentMessage{
			BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, fellowBikers),
			AgentId:     uuid.Nil,
			Kickout:     false,
		}
	}
}

func (agent *BaselineAgent) CreateReputationMessage() objects.ReputationOfAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents
	bikeID := agent.GetBike()
	fellowBikers := agent.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	bestAgent := uuid.Nil
	reputation := -1.0
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if agent.QueryReputation(agentID) > reputation {
			bestAgent = agentID
			reputation = agent.QueryReputation(agentID)
		}
	}
	return objects.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, agent.GetFellowBikers()),
		AgentId:     bestAgent,
		Reputation:  reputation,
	}
}

func (agent *BaselineAgent) CreateJoiningMessage() objects.JoiningAgentMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents

	// don't know how to write this, i want to pass
	return objects.JoiningAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, agent.GetFellowBikers()),
		AgentId:     uuid.Nil,
		BikeId:      uuid.Nil,
	}
}

func (agent *BaselineAgent) CreateLootboxMessage() objects.LootboxMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents

	// no need to change (? maybe)
	return objects.LootboxMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, agent.GetFellowBikers()),
		LootboxId:   agent.ProposeDirection(),
	}
}

func (agent *BaselineAgent) CreateGoverenceMessage() objects.GovernanceMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents

	// true message of decide governance
	return objects.GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[objects.IBaseBiker](agent, agent.GetFellowBikers()),
		BikeId:       agent.ChangeBike(),
		GovernanceId: int(agent.DecideGovernance()),
	}
}

func (agent *BaselineAgent) CreateForcesMessage() objects.ForcesMessage {
	// Currently this returns a default message which sends to all bikers on the biker agent's bike
	// For team's agent, add your own logic to communicate with other agents

	// true value
	agent.DecideForce(agent.ProposeDirection())
	return objects.ForcesMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, agent.GetFellowBikers()),
		AgentId:     agent.GetID(),
		AgentForces: agent.GetForces(),
	}
}

func (agent *BaselineAgent) CreateVoteLootboxDirectionMessage() objects.VoteLootboxDirectionMessage {
	lootBoxes := agent.GetGameState().GetLootBoxes() // Assuming this returns all loot boxes in the game.
	voteMap := make(voting.IdVoteMap)

	distanceMap := make(map[uuid.UUID]float64)
	for _, lootBox := range lootBoxes {
		distance := physics.ComputeDistance(agent.GetLocation(), lootBox.GetPosition())
		if distance <= 20 {
			distanceMap[lootBox.GetID()] = distance
		}
	}

	var totalDistance float64
	for _, distance := range distanceMap {
		totalDistance += distance
	}
	//averageDistance := totalDistance / float64(len(distanceMap))

	// normalize?
	// for lootBoxID, distance := range distanceMap {
	// 	rank := 1 - math.Abs(distance-averageDistance)/averageDistance
	// 	voteMap[lootBoxID] = int(rank * 100)
	// }

	return objects.VoteLootboxDirectionMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, agent.GetFellowBikers()),
		VoteMap:     voteMap,
	}
}

func (agent *BaselineAgent) CreateVoteRulerMessage() objects.VoteRulerMessage {
	reputationMap := agent.GetReputation()
	honestyMap := agent.GetHonestyMatrix()
	honestyThreshold := 0.8

	var rulerCandidate uuid.UUID
	highestCombinedScore := 0.0

	// iterate through all agents to find the best candidate.
	for _, fellowAgent := range agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetAgents() {
		agentID := fellowAgent.GetID()
		if agentID == agent.GetID() {
			continue
		}

		agentReputation := reputationMap[agentID]
		agentHonesty := honestyMap[agentID]

		if agentReputation >= reputationThreshold && agentHonesty >= honestyThreshold {
			combinedScore := agentReputation + agentHonesty
			if combinedScore > highestCombinedScore {
				highestCombinedScore = combinedScore
				rulerCandidate = agentID
			}
		}
	}

	if rulerCandidate == uuid.Nil || rulerCandidate == agent.GetID() {
		rulerCandidate = agent.GetID()
	}

	voteRulerMap := make(voting.IdVoteMap)
	voteRulerMap[rulerCandidate] = 1

	return objects.VoteRulerMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, agent.GetFellowBikers()),
		VoteMap:     voteRulerMap,
	}
}

func (agent *BaselineAgent) CreateVotekickoutMessage() objects.VoteKickoutMessage {
	reputationMap := agent.GetReputation()
	honestyMap := agent.GetHonestyMatrix()

	var lowestScoreAgentID uuid.UUID
	voteResults := make(map[uuid.UUID]int)
	lowestReputation := 1.0
	lowestHonesty := 1.0

	fellowBikers := agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetAgents()

	for _, fellowAgent := range fellowBikers {
		agentID := fellowAgent.GetID()
		if agentID == agent.GetID() {
			continue
		}

		agentReputation := reputationMap[agentID]
		agentHonesty := honestyMap[agentID]

		if agentReputation <= lowestReputation && agentHonesty <= lowestHonesty {
			lowestReputation = agentReputation
			lowestHonesty = agentHonesty
			lowestScoreAgentID = agentID
		}
	}

	voteResults[lowestScoreAgentID] = 1
	return objects.VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](agent, fellowBikers),
		VoteMap:     voteResults,
	}
}

func (agent *BaselineAgent) HandleKickoutMessage(msg objects.KickoutAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender().GetID()
	agentId := msg.AgentId
	kickout := msg.Kickout
	if agentId == agent.GetID() && kickout == true {
		agent.DecreaseHonesty(sender, ifKickOutMsgValue)
	}
	//fmt.Println("message kickout")
}

func (agent *BaselineAgent) HandleReputationMessage(msg objects.ReputationOfAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender().GetID()
	agentId := msg.AgentId
	reputation := msg.Reputation

	if agentId == agent.GetID() && reputation < reputationThreshold {
		agent.DecreaseHonesty(sender, ifReputationIsLow)
	}
	//fmt.Println("message reputaion")

}

func (agent *BaselineAgent) HandleJoiningMessage(msg objects.JoiningAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
}

func (agent *BaselineAgent) HandleLootboxMessage(msg objects.LootboxMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	senderID := sender.GetID()
	lootboxId := msg.LootboxId

	if sender.GetColour() == agent.GetColour() && lootboxId != agent.targetLoot {
		agent.DecreaseHonesty(senderID, dislikeLootBox)
	}
	if lootboxId == agent.targetLoot {
		agent.IncreaseHonesty(senderID, likeLootBox)
	}
	//fmt.Println("message lootbox")
}

func (agent *BaselineAgent) HandleGovernanceMessage(msg objects.GovernanceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	senderID := sender.GetID()
	//bikeId := msg.BikeId
	governanceId := msg.GovernanceId
	if agent.GetColour() == sender.GetColour() {
		if governanceId == 0 {
			agent.IncreaseHonesty(senderID, sameColorDictator)
		}
		if governanceId == 1 {
			agent.IncreaseHonesty(senderID, sameColorLeader)
		}
		if governanceId == 2 {
			agent.DecreaseHonesty(senderID, sameColorDictator)
		}
	} else {
		if governanceId == 0 {
			agent.IncreaseHonesty(senderID, differentColorDemorcracy)
		}
		if governanceId == 1 {
			agent.DecreaseHonesty(senderID, differntColorLeader)
		}
		if governanceId == 2 {
			agent.DecreaseHonesty(senderID, differentColorDictatorship)
		}
	}
	//fmt.Println("message governance")
}

func (agent *BaselineAgent) HandleForcesMessage(msg objects.ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	agentId := msg.AgentId
	agentForces := msg.AgentForces

	if sender.GetColour() == agent.lootBoxColour && agentForces.Pedal < 0.4 {
		agent.DecreaseHonesty(agentId, ifForcesTooLow)
		//fmt.Println("message forces")
	}

}

func (agent *BaselineAgent) HandleVoteGovernanceMessage(msg objects.VoteGoveranceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap

	//don't think this is useful

}

func (agent *BaselineAgent) HandleVoteLootboxDirectionMessage(msg objects.VoteLootboxDirectionMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.
	// sender := msg.BaseMessage.GetSender()
	// voteMap := msg.VoteMap

	//useful?
}

func (agent *BaselineAgent) HandleVoteRulerMessage(msg objects.VoteRulerMessage) {

	senderID := msg.BaseMessage.GetSender().GetID()
	senderColor := msg.BaseMessage.GetSender().GetColour()

	var rulerAgentID uuid.UUID
	highestScore := -1.0

	for agentID, score := range msg.VoteMap {
		// find the agent with the highest vote score
		if score > highestScore {
			highestScore = score
			rulerAgentID = agentID
		}
	}
	//fmt.Println("message vote ruler")

	if rulerAgentID != agent.GetID() {
		if senderColor == agent.GetColour() {
			agent.IncreaseHonesty(senderID, ifRulerSameColor)
		} else {
			agent.DecreaseHonesty(senderID, ifRulerDifferentColor)
		}
	} else {
		agent.IncreaseHonesty(senderID, ifWeAreRuler)
	}
}

func (agent *BaselineAgent) HandleVoteKickoutMessage(msg objects.VoteKickoutMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	sender := msg.BaseMessage.GetSender()
	senderID := sender.GetID()
	voteMap := msg.VoteMap

	if voteMap[agent.GetID()] > 0 {
		agent.DecreaseHonesty(senderID, ifKickOutUs)
		//fmt.Println("message kickout")
	}
}
