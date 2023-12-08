package team2

import (
	"SOMAS2023/internal/clients/team2/agent"
	"SOMAS2023/internal/common/objects"
)

// this function is going to be called by the server to instantiate bikers in the MVP
func GetBiker(baseBiker *objects.BaseBiker) objects.IBaseBiker {
	return agent.NewBaseTeam2Biker(baseBiker)
}
