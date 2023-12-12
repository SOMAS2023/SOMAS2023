package team6

import (
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

func (bb *Team6Biker) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	weights := make(map[uuid.UUID]float64)
	agents := bb.GetFellowBikers()
	for _, agent := range agents {
		// Based on our calculation of reputation, change the weights of voting -- to do
		weights[agent.GetID()] = 1.0
	}
	return weights
}
