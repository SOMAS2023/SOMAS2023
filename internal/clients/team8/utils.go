package team_8

import (
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
