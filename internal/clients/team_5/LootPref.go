package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

func ProposeDirection(gameState objects.IGameState, b *team5Agent) uuid.UUID {
	preferenceMap := make(map[uuid.UUID]float64)

	// Get the megabike, agent and lootboxes
	bikeId := b.GetBike()
	bike := gameState.GetMegaBikes()[bikeId]
	lootBoxes := gameState.GetLootBoxes()

	// Weights for distance, energy and colour
	var wd, we, wc = 0.3, 0.3, 0.4

	averageEnergyOthers := calculateAverageEnergyOthers(gameState, b)                 // Average energy of other agents on the bike
	urgencyFactor := averageEnergyPreference(b.GetEnergyLevel(), averageEnergyOthers) // Urgency factor based on agent's energy level compared to his bike mates

	// Calculate the preference for each lootbox
	for id, loot := range lootBoxes {
		distance := calculateDistance(bike.GetPosition(), loot.GetPosition())
		energy := energyPreference(b.GetEnergyLevel(), loot.GetTotalResources())
		colour := colourMatch(b.GetColour(), loot.GetColour())

		weightedDistancePreference := urgencyFactor * wd / (0.01 * distance)
		weightedEnergyPreference := we * energy
		weightedColourPreference := wc * colour

		preference := weightedDistancePreference + weightedEnergyPreference + weightedColourPreference
		preferenceMap[id] = preference
	}

	// Find the lootbox with the highest preference
	var max = 0.0
	var prefLootId uuid.UUID

	for lootId, preference := range preferenceMap {
		if preference > max {
			prefLootId = lootId
			max = preference
		}
	}

	// Return the lootbox with the highest preference
	return prefLootId
}

// Calculate the Euclidean distance between two coordinates
func calculateDistance(a, b utils.Coordinates) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// Calculate the preference for a lootbox based on colour
func colourMatch(agentColour, lootColour utils.Colour) float64 {
	if agentColour == lootColour {
		return 0.1 // Colour match between agent and lootbox
	}
	return 0.0 // No colour match between agent and lootbox
}

// Calculate the preference for a lootbox based on agent energy
func energyPreference(agentEnergy, lootResources float64) float64 {
	return lootResources * math.Pow(1/(1+agentEnergy), 2) // Quadratic function for energy preference as to give a greater effect on urgency to replenish energy when energy gets lower
}

// Calculate the altruism factor based on the agent's energy level
func averageEnergyPreference(agentEnergy, averageEnergyOthers float64) float64 {
	if agentEnergy < averageEnergyOthers {
		return 1.5 // Agent has less energy than average of other agents, so it is more urgent to replenish energy
	}
	return 1.0 // Agent has more energy than average of other agents, so it is less urgent to replenish energy
}

// Calculate the average energy of other agents
func calculateAverageEnergyOthers(gameState objects.IGameState, b *team5Agent) float64 {
	totalEnergy := 0.0

	id := b.GetID()

	agents := b.GetFellowBikers()
	totAgents := len(agents) - 1

	if totAgents == 0 {
		return 0
	}

	// Calculate the total energy of other agents on the bike
	for _, agent := range agents {
		if agent.GetID() != id {
			totalEnergy += agent.GetEnergyLevel()
		}
	}

	return totalEnergy / float64(totAgents)
}
