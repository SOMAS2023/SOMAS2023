package team5Agent

import (
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

//Only called when agent is democratic leader

func (t5 *team5Agent) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	weights := make(map[uuid.UUID]float64)
	agents := t5.GetFellowBikers()
	for _, agent := range agents {
		weights[agent.GetID()] = 1.0
	}
	return weights
}

//Only called when agent is dictator

func (t5 *team5Agent) DictateDirection() uuid.UUID {
	nearest := t5.ProposeDirection()
	return nearest
}

// needs fixing never kicks out
func (t5 *team5Agent) DecideKickOut() []uuid.UUID {
	return (make([]uuid.UUID, 0))
}

// needs fixing always allocates evenly
func (t5 *team5Agent) DecideDictatorAllocation() voting.IdVoteMap {
	bikeID := t5.GetBike()
	fellowBikers := t5.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	equalDist := 1.0 / float64(len(fellowBikers))
	for _, agent := range fellowBikers {
		distribution[agent.GetID()] = equalDist
	}
	return distribution
}
