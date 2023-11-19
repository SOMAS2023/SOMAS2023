package LootAlloc

import (
	"github.com/google/uuid"

	"SOMAS2023/internal/common/objects"
)

func calculateResourceAllocation(gameState objects.IGameState, self objects.IBaseBiker, method string) map[uuid.UUID]float32 {
	allocations := make(map[uuid.UUID]float32)

	//how to get id of my megabike?
	var bikeID uuid.UUID
	// bikeID = self.megaBikeId

	bike := gameState.GetMegaBikes()[bikeID]
	agentsOnBike := bike.GetAgents()

	for _, agent := range agentsOnBike {
		allocations[agent.GetID()] = generateAllocation(agent, self, method)
	}

	allocations = normaliseMap(allocations)

	return allocations
}

func generateAllocation(agent objects.IBaseBiker, self objects.IBaseBiker, method string) float32 {
	var value float32

	switch method {
	case "equal":
		value = 1
	case "greedy":
		if agent.GetID() == self.GetID() {
			value = 1
		} else {
			value = 0
		}
	default:
		//default to equal
		value = 1
	}

	//add more interesting allocation methods

	return value

}

func normaliseMap(m map[uuid.UUID]float32) map[uuid.UUID]float32 {
	sum := sumMap(m)

	for id, val := range m {
		m[id] = val / sum
	}

	return m
}

func sumMap(m map[uuid.UUID]float32) float32 {
	var sum float32 = 0
	for _, val := range m {
		sum += val
	}

	return sum
}
