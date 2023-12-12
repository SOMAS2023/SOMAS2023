package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

func (t5 *team5Agent) calculateResourceAllocation(method ResourceAllocationMethod) map[uuid.UUID]float64 {
	allocations := make(map[uuid.UUID]float64)

	//how to get id of my megabike?

	agentsOnBike := t5.GetFellowBikers()

	for _, agent := range agentsOnBike {
		allocations[agent.GetID()] = t5.generateAllocation(agent, method)
	}

	allocations = normaliseMap(allocations)

	return allocations
}

func (t5 *team5Agent) generateAllocation(agent objects.IBaseBiker, method ResourceAllocationMethod) float64 {
	var value float64

	switch method {
	case Equal:
		value = 1
	case Greedy:
		if agent.GetID() == t5.GetID() {
			value = 1
		} else {
			value = 0
		}
	case Needs:
		value = 1 - agent.GetEnergyLevel()
	case Contributions:
		value = agent.GetForces().Pedal * utils.MovingDepletion
	case Reputation:
		value = t5.QueryReputation(agent.GetID())
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
