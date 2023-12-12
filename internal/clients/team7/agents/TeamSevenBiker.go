package agents

import (
	"SOMAS2023/internal/clients/team7/frameworks"
	objects "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

	"math/rand"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

type ITeamSevenBiker interface {
	objects.IBaseBiker
}

type BaseTeamSevenBiker struct {
	*objects.BaseBiker    // BaseBiker inherits functions from BaseAgent such as GetID(), GetAllMessages() and UpdateAgentInternalState()
	navigationFramework   *frameworks.NavigationDecisionFramework
	bikeDecisionFramework *frameworks.BikeDecisionFramework
	opinionFramework      *frameworks.OpinionFramework
	socialNetwork         *frameworks.SocialNetwork
	environmentHandler    *frameworks.EnvironmentHandler
	personality           *frameworks.Personality

	// Memory
	memoryLength          int
	proposedDirections    []float64
	bikeProposedLootboxes []uuid.UUID
	locations             []utils.Coordinates
	myProposedLootboxes   []uuid.UUID
	distanceFromMyLootbox []float64
	time                  int

	votedForResources bool
	voteAllocationMap voting.IdVoteMap
	voteDirectionMap  voting.LootboxVoteMap
	voteKickingMap    map[uuid.UUID]int
	voteJoiningMap    map[uuid.UUID]bool

	reputationMessages           []objects.ReputationOfAgentMessage
	kickoutMessages              []objects.KickoutAgentMessage
	lootboxMessages              []objects.LootboxMessage
	joiningMessages              []objects.JoiningAgentMessage
	governanceMessages           []objects.GovernanceMessage
	forcesMessages               []objects.ForcesMessage
	voteGovernanceMessages       []objects.VoteGoveranceMessage
	voteLootboxDirectionMessages []objects.VoteLootboxDirectionMessage
	voteRulerMessages            []objects.VoteRulerMessage
	voteKickoutMessages          []objects.VoteKickoutMessage
	voteAllocationMessages       []objects.VoteAllocationMessage

	currentOpinionsOfAgents    map[uuid.UUID]float64
	currentOpinionsOfLootboxes map[uuid.UUID]float64
}

// Produce new BaseTeamSevenBiker
func NewBaseTeamSevenBiker(baseBiker *objects.BaseBiker) *BaseTeamSevenBiker {
	agentId := baseBiker.GetID()
	personality := frameworks.NewDefaultPersonality()
	return &BaseTeamSevenBiker{
		BaseBiker:             baseBiker,
		navigationFramework:   frameworks.NewNavigationDecisionFramework(),
		bikeDecisionFramework: frameworks.NewBikeDecisionFramework(),
		opinionFramework:      frameworks.NewOpinionFramework(frameworks.OpinionFrameworkInputs{}),
		socialNetwork:         frameworks.NewSocialNetwork(agentId, personality),
		personality:           personality,
		environmentHandler:    frameworks.NewEnvironmentHandler(baseBiker.GetGameState(), baseBiker.GetBike(), agentId),

		memoryLength:       10,
		time:               -1,
		proposedDirections: []float64{0, 0},

		reputationMessages:           make([]objects.ReputationOfAgentMessage, 0),
		kickoutMessages:              make([]objects.KickoutAgentMessage, 0),
		lootboxMessages:              make([]objects.LootboxMessage, 0),
		joiningMessages:              make([]objects.JoiningAgentMessage, 0),
		governanceMessages:           make([]objects.GovernanceMessage, 0),
		forcesMessages:               make([]objects.ForcesMessage, 0),
		voteGovernanceMessages:       make([]objects.VoteGoveranceMessage, 0),
		voteLootboxDirectionMessages: make([]objects.VoteLootboxDirectionMessage, 0),
		voteRulerMessages:            make([]objects.VoteRulerMessage, 0),
		voteKickoutMessages:          make([]objects.VoteKickoutMessage, 0),

		currentOpinionsOfAgents:    make(map[uuid.UUID]float64),
		currentOpinionsOfLootboxes: make(map[uuid.UUID]float64),
	}
}

func (biker *BaseTeamSevenBiker) UpdateGameState(gameState objects.IGameState) {
	biker.BaseBiker.UpdateGameState(gameState)
	biker.environmentHandler.UpdateGameState(gameState)
	biker.environmentHandler.UpdateCurrentBikeId(biker.GetBike())
}

// Override UpdateAgentInternalState
func (biker *BaseTeamSevenBiker) UpdateAgentInternalState() {
	biker.time++
	biker.environmentHandler.UpdateCurrentBikeId(biker.GetBike())

	fellowBikers := biker.environmentHandler.GetAgentsOnCurrentBike()

	// First formulate the data that we have access to directly
	agentColours := make(map[uuid.UUID]utils.Colour)
	agentEnergyLevels := make(map[uuid.UUID]float64)

	agentIds := make([]uuid.UUID, len(fellowBikers))
	for i, fellowBiker := range fellowBikers {
		agentId := fellowBiker.GetID()
		agentIds[i] = agentId
		agentColours[agentId] = fellowBiker.GetColour()
		agentEnergyLevels[agentId] = fellowBiker.GetEnergyLevel()
		// TODO: Implement once we can message biker to ask for allocation
		// if biker.votedForResources {
		// 	agentResourceVotes[agentId] = fellowBiker.DecideAllocation()
		// 	biker.votedForResources = false
		// }
	}

	// Formulate data based on messages and communication
	allAgentForceInformation := make(map[uuid.UUID](map[uuid.UUID]utils.Forces))
	for _, msg := range biker.forcesMessages {
		agentId := msg.AgentId
		allAgentForceInformation[agentId] = make(map[uuid.UUID]utils.Forces)
		allAgentForceInformation[agentId][msg.GetSender().GetID()] = msg.AgentForces
	}

	agentForces := make(map[uuid.UUID]utils.Forces)
	// Use the agent force from the agent with the highest trust level
	trustLevels := biker.socialNetwork.GetAverageTrustLevels()
	for agentId, agentForceInformation := range allAgentForceInformation {
		mostTrustedAgentId := uuid.Nil
		for senderId := range agentForceInformation {
			if mostTrustedAgentId == uuid.Nil || trustLevels[senderId] > trustLevels[mostTrustedAgentId] {
				mostTrustedAgentId = senderId
			}
		}
		agentForces[agentId] = agentForceInformation[mostTrustedAgentId]
	}

	// Get the agents' votes on allocation. At this point we are just trusting what they say to be true for now.
	agentResourceVotes := make(map[uuid.UUID]voting.IdVoteMap)
	for _, msg := range biker.voteAllocationMessages {
		agentResourceVotes[msg.GetSender().GetID()] = msg.VoteMap
	}

	socialNetworkInput := frameworks.SocialNetworkUpdateInput{
		AgentDecisions:     agentForces,
		AgentResourceVotes: agentResourceVotes,
		AgentEnergyLevels:  agentEnergyLevels,
		AgentColours:       agentColours,
		BikeTurnAngle:      biker.proposedDirections[len(biker.proposedDirections)-1],
	}

	biker.socialNetwork.UpdateSocialNetwork(agentIds, socialNetworkInput)

	// Next, update opinions
	// Update opinion on agents
	biker.updateOpinionsOnAgents(agentIds)
	// Update opinion on lootboxes
	biker.updateOpinionsOnLootboxes(agentIds)

	// Update memory
	if len(biker.locations) < biker.memoryLength {
		biker.locations = append(biker.locations, biker.GetLocation())
	} else {
		biker.locations = append(biker.locations[1:], biker.GetLocation())
	}

	// Clear messages which were used in this round
	biker.reputationMessages = make([]objects.ReputationOfAgentMessage, 0)
	biker.lootboxMessages = make([]objects.LootboxMessage, 0)
	biker.forcesMessages = make([]objects.ForcesMessage, 0)
	biker.voteAllocationMessages = make([]objects.VoteAllocationMessage, 0)
	// These were not used but clear them just for good practice
	biker.voteLootboxDirectionMessages = make([]objects.VoteLootboxDirectionMessage, 0)
	biker.voteKickoutMessages = make([]objects.VoteKickoutMessage, 0)
	biker.voteGovernanceMessages = make([]objects.VoteGoveranceMessage, 0)
	biker.voteRulerMessages = make([]objects.VoteRulerMessage, 0)
	biker.kickoutMessages = make([]objects.KickoutAgentMessage, 0)
	biker.joiningMessages = make([]objects.JoiningAgentMessage, 0)
	biker.governanceMessages = make([]objects.GovernanceMessage, 0)
}

func (biker *BaseTeamSevenBiker) get2DReputationMap() map[uuid.UUID](map[uuid.UUID]float64) {
	// Get reputation of each agent from each message
	// This is a map of agentId to a map of agentId to reputation
	// {
	// 	agentA: {
	// 		agentB: 1,
	// 		agentC: 0.5,
	// 	},
	// 	agentB: {
	// 		agentA: 0.5,
	// 		agentC: 0.2,
	// 	},
	// 	...
	// }
	//
	reputation2DMap := make(map[uuid.UUID](map[uuid.UUID]float64))
	for _, msg := range biker.reputationMessages {
		if _, ok := reputation2DMap[msg.AgentId]; !ok {
			reputation2DMap[msg.AgentId] = make(map[uuid.UUID]float64)
		}
		reputation2DMap[msg.AgentId][msg.GetSender().GetID()] = msg.Reputation
	}

	return reputation2DMap
}

func (biker *BaseTeamSevenBiker) updateOpinionsOnAgents(agentIds []uuid.UUID) {
	// Update opinions
	// Calculate overall opinion of each agent
	for _, agentId := range agentIds {
		_, hasOpinion := biker.currentOpinionsOfAgents[agentId]
		if !hasOpinion {
			biker.currentOpinionsOfAgents[agentId] = biker.socialNetwork.GetAverageTrustLevels()[agentId]
		}
	}

	// Opinion of agents
	reputation2DMap := biker.get2DReputationMap()
	for _, agentId := range agentIds {
		bikerOpinionsOfAgents, hasData := reputation2DMap[agentId]
		if hasData {
			opinionFrameworkInputs := frameworks.OpinionFrameworkInputs{
				AgentOpinion: bikerOpinionsOfAgents,
				Mindset:      biker.currentOpinionsOfAgents[agentId],
				OpinionType:  frameworks.AgentOpinions,
			}

			opinion := biker.opinionFramework.GetOpinion(opinionFrameworkInputs)
			biker.currentOpinionsOfAgents[agentId] = opinion
		}
	}
}

func (biker *BaseTeamSevenBiker) updateOpinionsOnLootboxes(agentIds []uuid.UUID) {
	// Opinion of lootboxes
	biker.currentOpinionsOfLootboxes = make(map[uuid.UUID]float64, 0)

	lootboxInterest := biker.getLootboxInterest()
	myProposedLootbox := biker.getDesiredLootboxId()
	if _, ok := lootboxInterest[myProposedLootbox]; !ok {
		lootboxInterest[myProposedLootbox] = make([]uuid.UUID, 0)
	}
	lootboxInterest[myProposedLootbox] = append(lootboxInterest[myProposedLootbox], biker.GetID())
	for lootboxId, agentIdsInterested := range lootboxInterest {
		opinionsOnLootbox := make(map[uuid.UUID]float64)
		for _, agentId := range agentIdsInterested {
			opinionsOnLootbox[agentId] = 1
		}
		for _, agentId := range agentIds {
			if _, ok := opinionsOnLootbox[agentId]; !ok {
				opinionsOnLootbox[agentId] = 0
			}
		}
		opinionFrameworkInputs := frameworks.OpinionFrameworkInputs{
			AgentOpinion: opinionsOnLootbox,
			Mindset:      biker.currentOpinionsOfLootboxes[lootboxId],
			OpinionType:  frameworks.LootboxOpinions,
		}

		opinion := biker.opinionFramework.GetOpinion(opinionFrameworkInputs)
		biker.currentOpinionsOfLootboxes[lootboxId] = opinion
	}
}

func (biker *BaseTeamSevenBiker) getLootboxInterest() map[uuid.UUID]([]uuid.UUID) {
	lootboxMap := make(map[uuid.UUID]([]uuid.UUID))
	// Get lootbox interest of each agent from each message
	for _, msg := range biker.lootboxMessages {
		sender := msg.GetSender().GetID()
		lootboxId := msg.LootboxId
		if _, ok := lootboxMap[lootboxId]; !ok {
			lootboxMap[lootboxId] = make([]uuid.UUID, 0)
		}
		lootboxMap[lootboxId] = append(lootboxMap[lootboxId], sender)
	}

	return lootboxMap
}

func (biker *BaseTeamSevenBiker) ProposeDirection() uuid.UUID {

	myProposedLootbox := biker.getDesiredLootboxId()

	// Update Memory
	if len(biker.myProposedLootboxes) < biker.memoryLength {
		biker.myProposedLootboxes = append(biker.myProposedLootboxes, myProposedLootbox)
	} else {
		biker.myProposedLootboxes = append(biker.myProposedLootboxes[1:], myProposedLootbox)
	}

	return myProposedLootbox
}

func (biker *BaseTeamSevenBiker) getDesiredLootboxId() uuid.UUID {
	myProposedLootboxObject := biker.environmentHandler.GetNearestLootBoxByColour(biker.GetColour())

	if biker.GetEnergyLevel() < 0.25 || myProposedLootboxObject == nil {
		return biker.environmentHandler.GetNearestLootBox().GetID()
	}

	return myProposedLootboxObject.GetID()
}

// TODO: Implement a strategy for choosing the final vote
func (biker *BaseTeamSevenBiker) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	myDesired := biker.getDesiredLootboxId()

	voteInputs := frameworks.VoteOnLootBoxesInput{
		LootBoxCandidates: proposals,
		MyPersonality:     biker.personality,
		MyDesired:         myDesired,
		MyOpinion:         biker.currentOpinionsOfLootboxes,
	}
	voteHandler := frameworks.NewVoteOnProposalsHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)
	return voteOutput
}

// Override base biker functions
func (biker *BaseTeamSevenBiker) DecideForce(direction uuid.UUID) {
	proposedLootbox := biker.environmentHandler.GetLootboxById(direction)

	var proposedLocation utils.Coordinates
	if proposedLootbox != nil {
		proposedLocation = proposedLootbox.GetPosition()
	} else {
		proposedLocation = utils.Coordinates{X: 0, Y: 0}
	}

	navInputs := frameworks.NavigationInputs{
		IsDestination:          proposedLootbox != nil,
		Destination:            proposedLocation,
		CurrentLocation:        biker.GetLocation(),
		CurrentEnergy:          biker.GetEnergyLevel(),
		ConscientiousnessLevel: biker.personality.Conscientiousness,
	}

	proposedDirection := biker.navigationFramework.GetTurnAngle(navInputs)

	navOutput := biker.navigationFramework.GetDecision(navInputs)

	biker.SetForces(navOutput)

	// Update Memory
	if len(biker.bikeProposedLootboxes) < biker.memoryLength {
		biker.bikeProposedLootboxes = append(biker.bikeProposedLootboxes, direction)
	} else {
		biker.bikeProposedLootboxes = append(biker.bikeProposedLootboxes[1:], direction)
	}

	if len(biker.proposedDirections) < biker.memoryLength {
		biker.proposedDirections = append(biker.proposedDirections, proposedDirection)
	} else {
		biker.proposedDirections = append(biker.proposedDirections[1:], proposedDirection)
	}

	if len(biker.myProposedLootboxes) > 0 {
		distanceFromMyProposal := biker.environmentHandler.GetDistanceBetweenLootboxes(direction, biker.myProposedLootboxes[len(biker.myProposedLootboxes)-1])
		if len(biker.distanceFromMyLootbox) < biker.memoryLength {
			biker.distanceFromMyLootbox = append(biker.distanceFromMyLootbox, distanceFromMyProposal)
		} else {
			biker.distanceFromMyLootbox = append(biker.distanceFromMyLootbox[1:], distanceFromMyProposal)
		}
	}
}

func (biker *BaseTeamSevenBiker) DecideAction() objects.BikerAction {
	// Decide whether to pedal, brake or coast
	decisionInputs := frameworks.BikeDecisionInputs{
		CurrentLocation: biker.GetLocation(),
		DecisionType:    frameworks.StayOrLeaveBike,
		AvailableBikes:  biker.environmentHandler.GetBikeMap(),
	}

	bikeOutput := biker.bikeDecisionFramework.GetDecision(decisionInputs)

	if bikeOutput.LeaveBike {
		return objects.ChangeBike
	}

	return objects.Pedal
}

// Vote on whether to accept new agent onto bike.
func (biker *BaseTeamSevenBiker) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates: pendingAgents,
	}
	voteHandler := frameworks.NewVoteToAcceptAgentHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	biker.voteJoiningMap = voteOutput

	return voteOutput
}

// Vote on allocation of resources
func (biker *BaseTeamSevenBiker) DecideAllocation() voting.IdVoteMap {
	agentIds := biker.environmentHandler.GetAgentIdsOnCurrentBike()

	voteInputs := frameworks.VoteOnAllocationInput{
		AgentCandidates: agentIds,
		MyPersonality:   biker.personality,
		MyId:            biker.GetID(),
	}

	voteHandler := frameworks.NewVoteOnAllocationHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	biker.votedForResources = true
	biker.voteAllocationMap = voteOutput
	return voteOutput
}

// Vote on kicking agent off bike.
func (biker *BaseTeamSevenBiker) VoteForKickout() map[uuid.UUID]int {

	fellowBikerIds := biker.environmentHandler.GetAgentIdsOnCurrentBike()

	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates:      fellowBikerIds,
		CurrentSocialNetwork: biker.socialNetwork.GetSocialNetwork(),
	}
	voteHandler := frameworks.NewVoteToKickAgentHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	biker.voteKickingMap = voteOutput

	return voteOutput
}

// Vote on Leader
func (biker *BaseTeamSevenBiker) VoteLeader() voting.IdVoteMap {
	agentIds := biker.environmentHandler.GetAgentIdsOnCurrentBike()

	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates:      agentIds,
		CurrentSocialNetwork: biker.socialNetwork.GetSocialNetwork(),
	}
	voteHandler := frameworks.NewVoteOnLeaderHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	return voteOutput
}

// Vote on Dictator
func (biker *BaseTeamSevenBiker) VoteDictator() voting.IdVoteMap {
	agentIds := biker.environmentHandler.GetAgentIdsOnCurrentBike()

	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates:      agentIds,
		CurrentSocialNetwork: biker.socialNetwork.GetSocialNetwork(),
	}
	voteHandler := frameworks.NewVoteOnDictatorHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	return voteOutput
}

// Vote on governance
func (biker *BaseTeamSevenBiker) DecideGovernance() utils.Governance {
	voteHandler := frameworks.NewVoteOnGovernanceHandler()
	voteOutput := voteHandler.GetDecision()

	return voteOutput
}

func (biker *BaseTeamSevenBiker) GetReputation() map[uuid.UUID]float64 {
	return biker.socialNetwork.GetCurrentTrustLevels()
}

func (biker *BaseTeamSevenBiker) QueryReputation(agentId uuid.UUID) float64 {
	trustLevels := biker.socialNetwork.GetCurrentTrustLevels()
	return trustLevels[agentId]
}

// This function updates all the messages for that agent i.e. both sending and receiving.
// And returns the new messages from other agents to your agent
func (biker *BaseTeamSevenBiker) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
	messages := make([]messaging.IMessage[objects.IBaseBiker], 0)

	// Get all the trust levels of the agents on the bike
	trustLevels := biker.socialNetwork.GetCurrentTrustLevels()
	for agentId, trustLevel := range trustLevels {
		reputationMessage := biker.CreateReputationMessage(agentId, trustLevel)
		messages = append(messages, reputationMessage)

		kickoutMessage := biker.CreatekickoutMessage(agentId, false)
		if trustLevel < 0.2 {
			kickoutMessage = biker.CreatekickoutMessage(agentId, true)
		}
		messages = append(messages, kickoutMessage)
	}

	voteKickoutMessage := biker.CreateVotekickoutMessage()
	messages = append(messages, voteKickoutMessage)

	voteDirectionMessage := biker.CreateVoteLootboxDirectionMessage()
	messages = append(messages, voteDirectionMessage)

	forcesMessage := biker.CreateForcesMessage()
	messages = append(messages, forcesMessage)

	voteAllocationMessage := biker.CreateVoteAllocationMessage()
	messages = append(messages, voteAllocationMessage)

	return messages
}

func (biker *BaseTeamSevenBiker) CreatekickoutMessage(agentId uuid.UUID, kickout bool) objects.KickoutAgentMessage {
	return objects.KickoutAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		AgentId:     agentId,
		Kickout:     kickout,
	}
}

func (biker *BaseTeamSevenBiker) CreateReputationMessage(agentId uuid.UUID, reputation float64) objects.ReputationOfAgentMessage {
	return objects.ReputationOfAgentMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		AgentId:     agentId,
		Reputation:  reputation,
	}
}

func (biker *BaseTeamSevenBiker) CreateLootboxMessage() objects.LootboxMessage {
	return objects.LootboxMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		LootboxId:   biker.getDesiredLootboxId(),
	}
}

func (biker *BaseTeamSevenBiker) CreateGoverenceMessage() objects.GovernanceMessage {
	return objects.GovernanceMessage{
		BaseMessage:  messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		BikeId:       biker.GetBike(),
		GovernanceId: 0, // Always propose democracy (for now)
	}
}

func (biker *BaseTeamSevenBiker) CreateForcesMessage() objects.ForcesMessage {
	return objects.ForcesMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		AgentId:     biker.GetID(),
		AgentForces: biker.GetForces(),
	}
}

func (biker *BaseTeamSevenBiker) CreateVoteLootboxDirectionMessage() objects.VoteLootboxDirectionMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteLootboxDirectionMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		VoteMap:     biker.voteDirectionMap.GetVotes(),
	}
}

func (biker *BaseTeamSevenBiker) CreateVotekickoutMessage() objects.VoteKickoutMessage {
	// Low agreeableness => Uncooperative => More likely to lie about voting to kick off agent.
	// High agreeableness => Cooperative => Less likely to lie about voting to kick off agent.
	voteKickingMapMessage := biker.voteKickingMap
	randNum := rand.Float64()
	if biker.personality.Agreeableness < randNum {
		for agentId, vote := range biker.voteKickingMap {
			if vote == 1 {
				voteKickingMapMessage[agentId] = 0
			}
		}
	}

	return objects.VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		VoteMap:     voteKickingMapMessage,
	}
}

func (biker *BaseTeamSevenBiker) CreateVoteAllocationMessage() objects.VoteAllocationMessage {
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteAllocationMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		VoteMap:     biker.voteAllocationMap,
	}
}

func (biker *BaseTeamSevenBiker) HandleKickoutMessage(msg objects.KickoutAgentMessage) {
	biker.kickoutMessages = append(biker.kickoutMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleReputationMessage(msg objects.ReputationOfAgentMessage) {
	biker.reputationMessages = append(biker.reputationMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleJoiningMessage(msg objects.JoiningAgentMessage) {
	biker.joiningMessages = append(biker.joiningMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleLootboxMessage(msg objects.LootboxMessage) {
	biker.lootboxMessages = append(biker.lootboxMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleGovernanceMessage(msg objects.GovernanceMessage) {
	biker.governanceMessages = append(biker.governanceMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleForcesMessage(msg objects.ForcesMessage) {
	biker.forcesMessages = append(biker.forcesMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleVoteGovernanceMessage(msg objects.VoteGoveranceMessage) {
	biker.voteGovernanceMessages = append(biker.voteGovernanceMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleVoteLootboxDirectionMessage(msg objects.VoteLootboxDirectionMessage) {
	biker.voteLootboxDirectionMessages = append(biker.voteLootboxDirectionMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleVoteRulerMessage(msg objects.VoteRulerMessage) {
	biker.voteRulerMessages = append(biker.voteRulerMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleVoteKickoutMessage(msg objects.VoteKickoutMessage) {
	biker.voteKickoutMessages = append(biker.voteKickoutMessages, msg)
}

func (biker *BaseTeamSevenBiker) HandleVoteAllocationMessage(msg objects.VoteAllocationMessage) {
	biker.voteAllocationMessages = append(biker.voteAllocationMessages, msg)
}
