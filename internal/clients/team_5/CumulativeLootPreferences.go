package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"
	"sort"

	"github.com/google/uuid"
)

// calculates final preferences for loot boxes based on various factors
func CalculateLootBoxPreferences(gameState objects.IGameState, b *team5Agent) map[uuid.UUID]float64 {
	finalPreferences := make(map[uuid.UUID]float64)

	// retrieve agent and loot boxes from game state
	lootBoxes := gameState.GetLootBoxes()
	position := gameState.GetMegaBikes()[b.GetBike()].GetPosition()

	for id, lootBox := range lootBoxes {
		distance := calculateDistanceToLootbox(position, lootBox.GetPosition())
		colorPreference := calculateColorPreference(b.GetColour(), lootBox.GetColour())
		energyPreference := calculateEnergyPreference(b.GetEnergyLevel(), lootBox.GetTotalResources())
		// cumulativePreference := cumulativePreferences[id]

		// combine preferences (weights: 0.4 for distance, 0.3 for color, 0.2 for energy, 0.1 for cumulative)
		// ensure that if cant get first preference, get second preference and so on

		finalPreferences[id] = 0.4*(1/distance) + 0.3*colorPreference + 0.2*energyPreference /*+ 0.1*cumulativePreference*/
	}

	return finalPreferences
}

func SortPreferences(prefs map[uuid.UUID]float64) map[uuid.UUID]float64 {
	finalVotes := make(map[uuid.UUID]float64)
	ids := make([]uuid.UUID, 0, len(prefs))

	for id := range prefs {
		ids = append(ids, id)
	}

	sort.SliceStable(ids, func(i, j int) bool {
		return prefs[ids[i]] < prefs[ids[j]]
	})

	for idx, id := range ids {
		finalVotes[id] = float64(idx + 1)
	}

	return finalVotes

}

// ensure that if cant get first preference, get second preference and so on

// calculates the Euclidean distance between two points
func calculateDistanceToLootbox(a, b utils.Coordinates) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// calculates preference based on color match
func calculateColorPreference(agentColor, lootBoxColor utils.Colour) float64 {
	if agentColor == lootBoxColor {
		return 1.0
	}
	return 0.0
}

// calculates preference based on energy level and loot resources
// example: higher preference if energy level is low and loot is high
func calculateEnergyPreference(agentEnergy, lootResources float64) float64 {
	return lootResources * (1 - agentEnergy)
}
