package preference

import (
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/objects"
	"math"

	"github.com/google/uuid"
)

func ProposeDirection(gameState objects.IGameState, agentID uuid.UUID) map[uuid.UUID]float64 {
	preferenceMap := make(map[uuid.UUID]float64)
	agent := gameState.GetMegaBikes()[agentID]
	lootBoxes := gameState.GetLootBoxes()

	var preference float64
	var wd = 0.3
	var we = 0.3
	var wc = 0.4

	for id, loot := range lootBoxes {
		distance := calculateDistance(agent.GetPosition(), loot.GetPosition())
		energy := energyPreference(agent.GetEnergyLevel(), loot.GetTotalResources(), gameState, agentID) // averageEnergyOthers
		colour := colourMatch(agent.GetColour(), loot.GetColour())

		preference = wd/(1+distance) + we*energy + wc*colour

		preferenceMap[id] = preference
	}

	return preferenceMap
}

// calculates distance between biker and lootbox

func calculateDistance(a, b utils.Coordinates) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// calculates preference based on color match
func colourMatch(agentColour, lootColour utils.Colour) float64 {
	if agentColour == lootColour {
		return 1.0
	}
	return 0.2
}

/*

func energyPreference(agentEnergy, lootResources float64) float64 {
    return lootResources * math.Pow(1-agentEnergy, 2) // quadratic
}

*/

func energyPreference(agentEnergy, lootResources float64, gameState objects.IGameState, agentID uuid.UUID) float64 {
	averageEnergyOthers := calculateAverageEnergyOthers(gameState, agentID)
	altruismFactor := averageEnergyPreference(agentEnergy, lootResources, averageEnergyOthers)

	return altruismFactor * lootResources * math.Pow(1-agentEnergy, 2) // quadratic
}

func averageEnergyPreference(agentEnergy, lootResources float64, averageEnergyOthers float64) float64 {

	if agentEnergy < averageEnergyOthers {
		return 1.0 // high preference
	}
	return 0.2
}

func calculateAverageEnergyOthers(gameState objects.IGameState, agentID uuid.UUID) float64 {
	totalEnergy := 0.0
	megabike := gameState.GetMegaBikes()[agentID]

	agents := megabike.GetAgents()
	totAgents := len(agents) - 1

	for id, agent := range agents {
		if id != agentID {
			totalEnergy += agent.GetEnergyLevel()
		}
	}

	if totAgents == 0 {
		return 0
	}

	return totalEnergy / float64(totAgents)
}
