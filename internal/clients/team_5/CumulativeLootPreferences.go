package team5Agent

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math"
	"sort"

	"github.com/google/uuid"
)

// calculates final preferences for loot boxes based on various factors
func (t5 *team5Agent) CalculateLootBoxPreferences(gameState objects.IGameState, proposals map[uuid.UUID]uuid.UUID) map[uuid.UUID]float64 {
	finalPreferences := make(map[uuid.UUID]float64)

	var lootBox objects.ILootBox

	// retrieve agent and loot boxes from game state
	position := gameState.GetMegaBikes()[t5.GetBike()].GetPosition()
	audiPos := t5.GetGameState().GetAudi().GetPosition()

	for _, lootBoxID := range proposals {
		lootBox = gameState.GetLootBoxes()[lootBoxID]
		distanceFromBike := calculateDistanceToObject(position, lootBox.GetPosition())
		colorPreference := calculateColorPreference(t5.GetColour(), lootBox.GetColour())
		energyPreference := calculateEnergyPreference(t5.GetEnergyLevel(), lootBox.GetTotalResources())

		distanceFromAudi := calculateDistanceToObject(audiPos, lootBox.GetPosition())
		var audiModifier float64 = 0

		if distanceFromAudi < (2 * utils.CollisionThreshold) {
			audiModifier = -0.5
		}

		// cumulativePreference := cumulativePreferences[id]

		// combine preferences (weights: 0.4 for distance, 0.3 for color, 0.2 for energy, 0.1 for cumulative)
		// ensure that if cant get first preference, get second preference and so on

		finalPreferences[lootBoxID] = 0.4*(1/distanceFromBike) + 0.3*colorPreference + 0.2*energyPreference + audiModifier
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
func calculateDistanceToObject(a, t5 utils.Coordinates) float64 {
	return math.Sqrt(math.Pow(a.X-t5.X, 2) + math.Pow(a.Y-t5.Y, 2))
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
