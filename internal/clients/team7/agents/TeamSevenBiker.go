package agents

import (
	"SOMAS2023/internal/clients/team7/frameworks"
	objects "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

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
}

// Produce new BaseTeamSevenBiker
func NewBaseTeamSevenBiker(agentId uuid.UUID) *BaseTeamSevenBiker {
	baseBiker := objects.GetBaseBiker(utils.GenerateRandomColour(), agentId)
	personality := frameworks.NewDefaultPersonality()
	return &BaseTeamSevenBiker{
		BaseBiker:             baseBiker,
		navigationFramework:   frameworks.NewNavigationDecisionFramework(),
		bikeDecisionFramework: frameworks.NewBikeDecisionFramework(),
		opinionFramework:      frameworks.NewOpinionFramework(frameworks.OpinionFrameworkInputs{}),
		socialNetwork:         frameworks.NewSocialNetwork(personality),
		personality:           personality,
		environmentHandler:    frameworks.NewEnvironmentHandler(baseBiker.GetGameState(), baseBiker.GetBike(), agentId),
		memoryLength:          10,
		time:                  -1,
		proposedDirections:    []float64{0, 0},
	}
}

func (biker *BaseTeamSevenBiker) UpdateGameState(gameState objects.IGameState) {
	biker.BaseBiker.UpdateGameState(gameState)
	biker.environmentHandler.UpdateGameState(gameState)
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
	for _, fellowBiker := range fellowBikers {
		agentId := fellowBiker.GetID()
		agentIds = append(agentIds, agentId)
		agentForces[agentId] = fellowBiker.GetForces()
		agentColours[agentId] = fellowBiker.GetColour()
		agentEnergyLevels[agentId] = fellowBiker.GetEnergyLevel()
		if biker.votedForResources {
			agentResourceVotes[agentId] = fellowBiker.DecideAllocation()
			biker.votedForResources = false
		}
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
	if biker.GetEnergyLevel() < 0.25 {
		return biker.environmentHandler.GetNearestLootBox().GetID()
	}

	myProposedLootboxObject := biker.environmentHandler.GetNearestLootBoxByColour(biker.GetColour())
	var myProposedLootbox uuid.UUID
	if myProposedLootboxObject == nil {
		myProposedLootbox = biker.environmentHandler.GetNearestLootBox().GetID()
	} else {
		myProposedLootbox = myProposedLootboxObject.GetID()
	}

	// Update Memory
	if len(biker.myProposedLootboxes) < biker.memoryLength {
		biker.myProposedLootboxes = append(biker.myProposedLootboxes, myProposedLootbox)
	} else {
		biker.myProposedLootboxes = append(biker.myProposedLootboxes[1:], myProposedLootbox)
	}

	return myProposedLootbox
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

	distanceFromMyProposal := biker.environmentHandler.GetDistanceBetweenLootboxes(direction, biker.myProposedLootboxes[len(biker.myProposedLootboxes)-1])
	if len(biker.myProposedLootboxes) < biker.memoryLength {
		biker.distanceFromMyLootbox = append(biker.distanceFromMyLootbox, distanceFromMyProposal)
	} else {
		biker.distanceFromMyLootbox = append(biker.distanceFromMyLootbox[1:], distanceFromMyProposal)
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

	return voteOutput
}

// Vote on allocation of resources
func (biker *BaseTeamSevenBiker) DecideAllocation() voting.IdVoteMap {
	agentIds := biker.environmentHandler.GetAgentsOnCurrentBikeId()

	voteInputs := frameworks.VoteOnAllocationInput{
		AgentCandidates: agentIds,
		MyId:            biker.GetID(),
	}

	voteHandler := frameworks.NewVoteOnAllocationHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	biker.votedForResources = true
	return voteOutput
}

// Vote on whether to kick agent off bike
func (biker *BaseTeamSevenBiker) DecideKicking(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates: pendingAgents,
	}
	voteHandler := frameworks.NewVoteToKickAgentHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	return voteOutput
}

// Vote on Leader
// TODO: Uncomment when infrastructure have merged the new voting methods.
/*
func (biker *BaseTeamSevenBiker) VoteLeader() voting.IdVoteMap {
	agentIds := biker.environmentHandler.GetAgentsOnCurrentBikeId()

	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates: agentIds,
	}
	voteHandler := frameworks.NewVoteOnLeaderHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	return voteOutput
}

// Vote on Dictator
func (biker *BaseTeamSevenBiker) VoteDictator() voting.IdVoteMap {
	agentIds := biker.environmentHandler.GetAgentsOnCurrentBikeId()

	voteInputs := frameworks.VoteOnAgentsInput{
		AgentCandidates: agentIds,
	}
	voteHandler := frameworks.NewVoteOnDictatorHandler()
	voteOutput := voteHandler.GetDecision(voteInputs)

	return voteOutput
}
*/
