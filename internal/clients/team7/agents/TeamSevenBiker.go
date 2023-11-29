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
	votingFramework       *frameworks.VotingFramework
	environmentHandler    *frameworks.EnvironmentHandler
	personality           *frameworks.Personality

	previousProposedTurningDirection float64
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
		votingFramework:       frameworks.NewVotingFramework(),
		personality:           personality,
		environmentHandler:    frameworks.NewEnvironmentHandler(baseBiker.GetGameState(), baseBiker.GetMegaBikeId(), agentId),
	}
}

// Override base biker functions
func (biker *BaseTeamSevenBiker) DecideForce(direction uuid.UUID) {
	// Store previous proposed direction for next round's decisions

	proposedLootbox := biker.environmentHandler.GetLootboxById(direction)

	navInputs := frameworks.NavigationInputs{
		Destination:     proposedLootbox.GetPosition(),
		CurrentLocation: biker.GetLocation(),
	}

	biker.previousProposedTurningDirection = biker.navigationFramework.GetTurnAngle(navInputs)

	navOutput := biker.navigationFramework.GetDecision(navInputs)

	biker.SetForces(navOutput)
}

// Override UpdateAgentInternalState
func (biker *BaseTeamSevenBiker) UpdateAgentInternalState() {
	biker.BaseBiker.UpdateAgentInternalState()

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
		agentResourceVotes[agentId] = fellowBiker.DecideAllocation()
	}

	socialNetworkInput := frameworks.SocialNetworkUpdateInput{
		AgentDecisions:     agentForces,
		AgentResourceVotes: agentResourceVotes,
		AgentEnergyLevels:  agentEnergyLevels,
		AgentColours:       agentColours,
		BikeTurnAngle:      biker.previousProposedTurningDirection,
	}

	biker.socialNetwork.UpdateSocialNetwork(agentIds, socialNetworkInput)
}
