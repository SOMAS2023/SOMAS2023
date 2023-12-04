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
	kickOut := make([]uuid.UUID, 0)
	fellowBikers := t5.GetFellowBikers()
	if len(fellowBikers) > 4 {
		for _, agent := range fellowBikers {
			id := agent.GetID()
			if t5.QueryReputation(id) < 0.3 {
				kickOut = append(kickOut, id)
			}
		}
	}
	return kickOut
}

func (t5 *team5Agent) DecideDictatorAllocation() voting.IdVoteMap {
	return t5.calculateResourceAllocation(Reputation)
}
