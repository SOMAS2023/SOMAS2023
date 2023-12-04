package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

func (t5 *team5Agent) calculateResourceAllocation(gameState objects.IGameState) map[uuid.UUID]float64 {
	allocations := make(map[uuid.UUID]float64)

	//how to get id of my megabike?

	agentsOnBike := t5.GetFellowBikers()

	for _, agent := range agentsOnBike {
		allocations[agent.GetID()] = t5.generateAllocation(agent)
	}

	allocations = normaliseMap(allocations)

	return allocations
}

func (t5 *team5Agent) generateAllocation(agent objects.IBaseBiker) float64 {
	var value float64

	switch t5.resourceAllocationMethod {
	case "equal":
		value = 1
	case "greedy":
		if agent.GetID() == t5.GetID() {
			value = 1
		} else {
			value = 0
		}
	case "needs":
		value = 1 - agent.GetEnergyLevel()
	case "contributions":
		value = agent.GetForces().Pedal * utils.MovingDepletion
	// case "rep":
	// 	value = b.GetAgentReputation(agent.GetID())
	default:
		//default to equal
		value = 1
	}

	//add more interesting allocation methods

	return value

}

func normaliseMap(m map[uuid.UUID]float64) map[uuid.UUID]float64 {
	sum := sumMap(m)

	for id, val := range m {
		m[id] = val / sum
	}

	return m
}

func sumMap(m map[uuid.UUID]float64) float64 {
	var sum float64 = 0
	for _, val := range m {
		sum += val
	}

	return sum
}
