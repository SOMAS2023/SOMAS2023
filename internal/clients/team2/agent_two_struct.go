// AccessBaseBiker is an example of how to access BaseBiker fields and methods
// func (a *AgentTwo) AccessBaseBiker() {
//     // Accessing a field of BaseBiker
//     a.BaseBiker.SomeField = "some value"

//     // Calling a method of BaseBiker
//     a.BaseBiker.SomeMethod()
// }

// TODO: Reputation evaluation

package team2

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

func NewBaseTeam2Biker(agentId uuid.UUID) *AgentTwo {
	color := utils.GenerateRandomColour()
	baseBiker := objects.GetBaseBiker(color, agentId)
	return &AgentTwo{
		BaseBiker:          baseBiker,
		SocialCapital:      make(map[uuid.UUID]float64),
		Reputation:         make(map[uuid.UUID]float64),
		Institution:        make(map[uuid.UUID]float64),
		Network:            make(map[uuid.UUID]float64),
		GameIterations:     0,
		forgivenessCounter: 0,
		gameState:          nil,
		megaBikeId:         uuid.UUID{},
		bikeCounter:        make(map[uuid.UUID]int32),
		actions:            make([]Action, 0),
		soughtColour:       color,
		onBike:             false,
		energyLevel:        1.0,
		points:             0,
		forces:             utils.Forces{},
		allocationParams:   objects.ResourceAllocationParams{},
	}
}
