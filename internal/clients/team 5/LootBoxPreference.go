package preference

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

func ProposeDirection(gameState objects.IGameState, agentID uuid.UUID) map[uuid.UUID]float64 {
	preferenceMap := make(map[uuid.UUID]float64)

	// Get the megabike, agent and lootboxes
	bikeId, agent := findAgentandBike(gameState, agentID)
	bike := gameState.GetMegaBikes()[bikeId]
	lootBoxes := gameState.GetLootBoxes()

	// Weights for distance, energy and colour
	var wd, we, wc = 0.3, 0.3, 0.4

	// Calculate the preference for each lootbox
	for id, loot := range lootBoxes {
		distance := calculateDistance(bike.GetPosition(), bike.GetPosition())
		energy := energyPreference(agent.GetEnergyLevel(), loot.GetTotalResources(), gameState, agent.GetID(), bike)
		colour := colourMatch(agent.GetColour(), loot.GetColour())

		preference := wd/(1+distance) + we*energy + wc*colour
		preferenceMap[id] = preference
	}

	return preferenceMap
}

// Find the agent with the given ID and return the bike ID and agent
func findAgentandBike(gameState objects.IGameState, agentID uuid.UUID) (uuid.UUID, objects.IBaseBiker) {
	for id, bike := range gameState.GetMegaBikes() {
		for _, agent := range bike.GetAgents() {
			if agent.GetID() == agentID {
				return id, agent
			}
		}
	}
	return uuid.Nil, nil
}

// Calculate the Euclidean distance between two coordinates
func calculateDistance(a, b utils.Coordinates) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// Calculate the preference for a lootbox based on colour
func colourMatch(agentColour, lootColour utils.Colour) float64 {
	if agentColour == lootColour {
		return 1.0 // Colour match between agent and lootbox
	}
	return 0.2
}

// Calculate the preference for a lootbox based on agent energy
func energyPreference(agentEnergy, lootResources float64, gameState objects.IGameState, agentID uuid.UUID) float64 {
	averageEnergyOthers := calculateAverageEnergyOthers(gameState, agentID)                    // Average energy of other agents
	altruismFactor := averageEnergyPreference(agentEnergy, lootResources, averageEnergyOthers) // Altruism factor that takes into account other agents' energy

	return altruismFactor * lootResources * math.Pow(1-agentEnergy, 2) // Quadratic function for energy preference as to give a greater effect on urgency to replenish energy when energy gets lower
}

// Calculate the altruism factor based on the agent's energy level
func averageEnergyPreference(agentEnergy, lootResources float64, averageEnergyOthers float64) float64 {
	if agentEnergy < averageEnergyOthers {
		return 1.0 // Agent has less energy than average of other agents, so it is more urgent to replenish energy
	}
	return 0.2 // Agent has more energy than average of other agents, so it is less urgent to replenish energy
}

// Calculate the average energy of other agents
func calculateAverageEnergyOthers(gameState objects.IGameState, agentID uuid.UUID, megabike objects.IMegaBike) float64 {
	totalEnergy := 0.0

	agents := megabike.GetAgents()
	totAgents := len(agents) - 1

	if totAgents == 0 {
		return 0
	}

	// Calculate the total energy of other agents on the bike
	for _, agent := range agents {
		if agent.GetID() != agentID {
			totalEnergy += agent.GetEnergyLevel()
		}
	}

	return totalEnergy / float64(totAgents)
}
