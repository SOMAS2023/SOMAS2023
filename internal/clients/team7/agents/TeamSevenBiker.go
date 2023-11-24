package agents

import (
	"SOMAS2023/internal/clients/team7/frameworks"
	objects "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

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
	environmentHandler    *EnvironmentHandler
	personality           *Personality
}

// Produce new BaseTeamSevenBiker
func NewBaseTeamSevenBiker(agentId uuid.UUID) *BaseTeamSevenBiker {
	baseBiker := objects.GetBaseBiker(utils.GenerateRandomColour(), agentId)
	return &BaseTeamSevenBiker{
		BaseBiker:             baseBiker,
		navigationFramework:   frameworks.NewNavigationDecisionFramework(),
		bikeDecisionFramework: frameworks.NewBikeDecisionFramework(),
		opinionFramework:      frameworks.NewOpinionFramework(frameworks.OpinionFrameworkInputs{}),
		socialNetwork:         frameworks.NewSocialNetwork(),
		votingFramework:       frameworks.NewVotingFramework(),
		environmentHandler:    NewEnvironmentHandler(baseBiker.GetGameState(), baseBiker.GetMegaBikeId(), agentId),
		personality:           NewDefaultPersonality(),
	}
}

// Override base biker functions
func (biker *BaseTeamSevenBiker) DecideForce(direction uuid.UUID) {
	navInputs := frameworks.NavigationInputs{
		DesiredLootbox:  biker.environmentHandler.GetNearestLootBox().GetPosition(),
		CurrentLocation: biker.GetLocation(),
	}
	navOutput := biker.navigationFramework.GetDecision(navInputs)

	biker.SetForces(navOutput)
}

/*
// Ally will update this as soon as the infrastructure is merged!

// VOTING FUNCTIONS
func (biker *BaseTeamSevenBiker) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	voteInputs := frameworks.VoteInputs{
		DecisionType:   frameworks.VoteToAcceptNewAgent,
		Candidates:     pendingAgents,
		VoteParameters: frameworks.YesNo,
	}

	voteOutput := biker.votingFramework.GetDecision(voteInputs)

	return voteOutput
}
*/
