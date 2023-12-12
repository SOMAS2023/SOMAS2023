package agents

import (
	"SOMAS2023/internal/clients/team7/frameworks"
	objects "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

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
	voteKickoutMessgaes          []objects.VoteKickoutMessage
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
		voteKickoutMessgaes:          make([]objects.VoteKickoutMessage, 0),
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
	agentForces := make(map[uuid.UUID]utils.Forces)
	agentColours := make(map[uuid.UUID]utils.Colour)
	agentEnergyLevels := make(map[uuid.UUID]float64)
	agentResourceVotes := make(map[uuid.UUID]voting.IdVoteMap)

	agentIds := make([]uuid.UUID, len(fellowBikers))
	for i, fellowBiker := range fellowBikers {
		agentId := fellowBiker.GetID()
		agentIds[i] = agentId
		// agentForces[agentId] = fellowBiker.GetForces()
		agentColours[agentId] = fellowBiker.GetColour()
		agentEnergyLevels[agentId] = fellowBiker.GetEnergyLevel()
		// TODO: Implement once we can message biker to ask for allocation
		// if biker.votedForResources {
		// 	agentResourceVotes[agentId] = fellowBiker.DecideAllocation()
		// 	biker.votedForResources = false
		// }
	}

	socialNetworkInput := frameworks.SocialNetworkUpdateInput{
		AgentDecisions:     agentForces,
		AgentResourceVotes: agentResourceVotes,
		AgentEnergyLevels:  agentEnergyLevels,
		AgentColours:       agentColours,
		BikeTurnAngle:      biker.proposedDirections[len(biker.proposedDirections)-1],
	}

	biker.socialNetwork.UpdateSocialNetwork(agentIds, socialNetworkInput)

	// Update memory
	if len(biker.locations) < biker.memoryLength {
		biker.locations = append(biker.locations, biker.GetLocation())
	} else {
		biker.locations = append(biker.locations[1:], biker.GetLocation())
	}
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
	votes := make(voting.LootboxVoteMap)
	totOptions := len(proposals)
	normalDist := 1.0 / float64(totOptions)
	for _, proposal := range proposals {
		if val, ok := votes[proposal]; ok {
			votes[proposal] = val + normalDist
		} else {
			votes[proposal] = normalDist
		}
	}

	biker.voteDirectionMap = votes
	return votes
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
		IsDestination:   proposedLootbox != nil,
		Destination:     proposedLocation,
		CurrentLocation: biker.GetLocation(),
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
		MyId:            biker.GetID(),
	}

	voteHandler := frameworks.NewVoteOnAllocationHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	biker.votedForResources = true
	biker.voteAllocationMap = voteOutput
	return voteOutput
}

func (biker *BaseTeamSevenBiker) VoteForKickout() map[uuid.UUID]int {
	voteResults := make(map[uuid.UUID]int)

	fellowBikerIds := biker.environmentHandler.GetAgentIdsOnCurrentBike()

	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates:      fellowBikerIds,
		CurrentSocialNetwork: biker.socialNetwork.GetSocialNetwork(),
	}
	voteHandler := frameworks.NewVoteToKickAgentHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	for _, agent := range fellowBikerIds {
		if voteOutput[agent] {
			voteResults[agent] = 1
		} else {
			voteResults[agent] = 0
		}
	}

	biker.voteKickingMap = voteResults

	return voteResults
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
	// Currently this returns a default/meaningless message
	// For team's agent, add your own logic to communicate with other agents
	return objects.VoteKickoutMessage{
		BaseMessage: messaging.CreateMessage[objects.IBaseBiker](biker, biker.GetFellowBikers()),
		VoteMap:     biker.voteKickingMap,
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
	biker.voteKickoutMessgaes = append(biker.voteKickoutMessgaes, msg)
}
