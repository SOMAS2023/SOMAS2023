package team7

import (
	"SOMAS2023/internal/clients/team7/agents"
	"SOMAS2023/internal/common/objects"
)

func GetTeamSevenBiker(baseBiker *objects.BaseBiker) objects.IBaseBiker {
	return agents.NewBaseTeamSevenBiker(baseBiker)
}
