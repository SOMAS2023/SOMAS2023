package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math/rand"

	"github.com/google/uuid"
)

type BaselineAgent struct {
	objects.BaseBiker
	currentBike *objects.MegaBike
}

// DecideAction only pedal
func (agent *BaselineAgent) DecideAction() objects.BikerAction {
	if not objects.leave {
		return objects.Pedal
	} else {
		return objects.ChangeBike
	}

}

func (agent *BaselineAgent) DecideForces() {
	energyLevel := agent.GetEnergyLevel()

	randomPedalForce := rand.Float64() * energyLevel
	forces := utils.Forces{
		Pedal:   randomPedalForce,
		Brake:   0.0,
		Turning: 0.0,
	}
	println("forces for each round", forces)
}

func (agent *BaselineAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	trust := make(map[uuid.UUID]int)
	for _, agent := range pendingAgents {
		//trust
		if trust[agent] >= 5 {
			decision[agent] = true
		}
	}
	return decision
}


