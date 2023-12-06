package team2

import (
	"SOMAS2023/internal/clients/team2/agent"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
)

// this function is going to be called by the server to instantiate bikers in the MVP
func GetBiker(colour utils.Colour, id uuid.UUID) objects.IBaseBiker {
	return agent.NewBaseTeam2Biker(id, colour)
}
