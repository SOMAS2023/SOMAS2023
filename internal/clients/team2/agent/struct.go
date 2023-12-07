package agent

import (
	"SOMAS2023/internal/clients/team2/modules"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type IBaseBiker interface {
	objects.IBaseBiker
}

type AgentModules struct {
	Environment    *modules.EnvironmentModule
	SocialCapital  *modules.SocialCapital
	Decision       *modules.DecisionModule
	Utils          *modules.UtilsModule
	VotedDirection uuid.UUID
}

type AgentTwo struct {
	*objects.BaseBiker // Embedding the BaseBiker
	Modules            AgentModules
}

func NewBaseTeam2Biker(agentId uuid.UUID, colour utils.Colour) *AgentTwo {
	baseBiker := objects.GetBaseBiker(colour, agentId)
	baseBiker.GroupID = 2
	return &AgentTwo{
		BaseBiker: baseBiker,
		Modules: AgentModules{
			Environment:    modules.GetEnvironmentModule(baseBiker.GetID(), baseBiker.GetGameState(), baseBiker.GetBike()),
			SocialCapital:  modules.NewSocialCapital(),
			Decision:       modules.NewDecisionModule(),
			Utils:          modules.NewUtilsModule(),
			VotedDirection: uuid.UUID{},
		},
	}
}
