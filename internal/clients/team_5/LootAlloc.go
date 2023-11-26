package team5Agent

import (
	"github.com/google/uuid"

	"SOMAS2023/internal/common/objects"
)

func calculateResourceAllocation(gameState objects.IGameState, b *team5Agent) map[uuid.UUID]float64 {
	allocations := make(map[uuid.UUID]float64)

	//how to get id of my megabike?
	var bikeID uuid.UUID
	bikeID = b.GetBike()

	bike := gameState.GetMegaBikes()[bikeID]
	agentsOnBike := bike.GetAgents()

	for _, agent := range agentsOnBike {
		allocations[agent.GetID()] = generateAllocation(agent, b)
	}

	allocations = normaliseMap(allocations)

	return allocations
}

// gets Bike ID from gamestate, to be removed after getter added to basebiker
func getBikeIdFromGameState(b *team5Agent, gameState objects.IGameState) uuid.UUID {
	bikes := gameState.GetMegaBikes()

	for id, bike := range bikes {
		for _, agent := range bike.GetAgents() {
			if agent.GetID() == b.GetID() {
				return id
			}
		}
	}

	return uuid.Nil
}

func generateAllocation(agent objects.IBaseBiker, b *team5Agent) float64 {
	var value float64

	switch b.resourceAllocationMethod {
	case "equal":
		value = 1
	case "greedy":
		if agent.GetID() == b.GetID() {
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
