package team_8

import (
	"SOMAS2023/internal/common/utils"
	"math"
	"sort"

	"github.com/google/uuid"
)

// CalculateGiniIndexFromAB calculates the Gini index using the given values of A and B.
func CalculateGiniIndexFromAB(A, B float64) float64 {
	// Ensure that the denominator is not zero to avoid division by zero
	if A+B == 0 {
		return 0.0 // or handle this case according to your requirements
	}

	// Calculate the Gini index
	giniIndex := A / (A + B)

	return giniIndex
}

// sort loot boxes by their scores
func sortLootBoxesByScore(combinedScores map[uuid.UUID]float64) []uuid.UUID {
	// Create a slice of boxes to sort
	var boxes []uuid.UUID
	for boxID := range combinedScores {
		boxes = append(boxes, boxID)
	}

	// Sort the slice based on scores
	sort.Slice(boxes, func(i, j int) bool {
		return combinedScores[boxes[i]] > combinedScores[boxes[j]]
	})

	return boxes
}

func softmax(preferences map[uuid.UUID]float64) map[uuid.UUID]float64 {
	sum := 0.0
	for _, pref := range preferences {
		sum += math.Exp(pref)
	}

	softmaxPreferences := make(map[uuid.UUID]float64)
	for id, pref := range preferences {
		softmaxPreferences[id] = math.Exp(pref) / sum
	}

	return softmaxPreferences
}

// calculateDistance computes the Euclidean distance between two points
func calculateDistance(a, b utils.Coordinates) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// calculateColorPreference returns 1 if the colors match, 0 otherwise
func calculateColorPreference(agentColor, boxColor utils.Colour) float64 {
	if agentColor == boxColor {
		return 1
	}
	return 0
}

// calculateEnergyWeighting adjusts the distance preference based on the agent's energy level
func calculateEnergyWeighting(energyLevel float64) float64 {
	// Assuming the energy level is between 0 and 1, inverse it to give higher weight to closer loot boxes when energy is low
	return 1 - energyLevel
}

// rankByPreference sorts the loot boxes by preference
func rankByPreference(preferences map[uuid.UUID]float64) []uuid.UUID {
	type kv struct {
		ID         uuid.UUID
		Preference float64
	}

	var sorted []kv
	for id, pref := range preferences {
		sorted = append(sorted, kv{id, pref})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Preference > sorted[j].Preference
	})

	var rankedIDs []uuid.UUID
	for _, kv := range sorted {
		rankedIDs = append(rankedIDs, kv.ID)
	}

	return rankedIDs
}

// selectTopChoices selects the top choices based on the ranking
func selectTopChoices(rankedIDs []uuid.UUID, numChoices int) uuid.UUID {
	if len(rankedIDs) == 0 {
		return uuid.Nil // No loot boxes available
	}
	if numChoices > len(rankedIDs) {
		numChoices = len(rankedIDs)
	}
	// For this example, just select the top choice
	return rankedIDs[0]
}
