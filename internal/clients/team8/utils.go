package team_8

import (
	"SOMAS2023/internal/common/objects"
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

// the function is used to map uuid of agents to real baseAgent object
func (bb *Agent8) UuidToAgentMap(pendingAgents []uuid.UUID) map[uuid.UUID]objects.IBaseBiker {
	agentMap := make(map[uuid.UUID]objects.IBaseBiker)
	megaBikes := bb.GetGameState().GetMegaBikes()

	for _, megaBike := range megaBikes {
		for _, agent := range megaBike.GetAgents() {
			for _, uuid := range pendingAgents {
				if agent.GetID() == uuid {
					agentMap[uuid] = agent
				}
			}
		}
	}

	return agentMap
}

// CalculateAverageEnergy calculates the average energy level for agents on a specific bike.
func (bb *Agent8) CalculateAverageEnergy(bikeID uuid.UUID) float64 {
	// Step 1: Get fellowBikers from the specified bike
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()

	// Step 2: Ensure there is at least one agent
	if len(fellowBikers) == 0 {
		return 0.0 // or handle this case according to your requirements
	}

	// Step 3: Calculate the sum of energy levels
	sum := 0.0
	for _, agent := range fellowBikers {
		sum += agent.GetEnergyLevel()
	}

	// Step 4: Calculate the average
	average := sum / float64(len(fellowBikers))

	return average
}

// CountAgentsWithSameColour counts the number of agents with the same colour as the reference agent on a specific bike.
func (bb *Agent8) CountAgentsWithSameColour(bikeID uuid.UUID) int {
	// Step 1: Get reference colour from the BaseBiker
	referenceColour := bb.GetColour()

	// Step 2: Get fellowBikers from the specified bike
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()

	// Step 3: Ensure there is at least one agent
	if len(fellowBikers) == 0 {
		return 0 // or handle this case according to your requirements
	}

	// Step 4: Count agents with the same colour as the reference agent
	count := 0
	for _, agent := range fellowBikers {
		if agent.GetColour() == referenceColour {
			count++
		}
	}

	return count
}

func (bb *Agent8) calculateValueJudgement(utilityLevels []float64, agentGoals []int, targetGoal int, turns []bool) float64 {

	averageUtility := bb.calculateAverageUtility(utilityLevels)
	percentageSameGoal := bb.calculatePercentageSameGoal(agentGoals, targetGoal)
	probabilitySatisfiedLoops := bb.calculateProbabilitySatisfiedLoops(turns)

	// Calculate the average score
	averageScore := (averageUtility + percentageSameGoal + probabilitySatisfiedLoops) / 3
	return averageScore
}

func (bb *Agent8) calculateAverageOfCostAndPercentage(decisions []bool, energyLevels []float64, threshold float64) float64 {

	costPercentage := bb.calculateCostInCollectiveImprovement(decisions)
	percentageLowEnergy := bb.calculatePercentageLowEnergyAgents(energyLevels, threshold)

	// Calculate the average
	averageResult := (costPercentage + percentageLowEnergy) / 2
	return averageResult
}

// calculateAverageUtilityPercentage calculates the average of utility levels and returns the percentage
// the utilitylevels need additional parameters to calculate
func (bb *Agent8) calculateAverageUtility(utilityLevels []float64) float64 {
	var sum float64
	for _, value := range utilityLevels {
		sum += value
	}
	average_utility := sum / float64(len(utilityLevels))
	return average_utility
}

// calculatePercentageSameGoal calculates the percentage of agents with the same goal
func (bb *Agent8) calculatePercentageSameGoal(agentGoals []int, targetGoal int) float64 {
	var countSameGoal int
	for _, goal := range agentGoals {
		if goal == targetGoal {
			countSameGoal++
		}
	}
	totalAgents := len(agentGoals)
	if totalAgents == 0 {
		return 0.0
	}
	percentage := (float64(countSameGoal) / float64(totalAgents))
	return percentage
}

// calculateProbabilitySatisfiedLoops calculates the probability of having 'true' in the array
func (bb *Agent8) calculateProbabilitySatisfiedLoops(turns []bool) float64 {
	var countTrue int
	for _, result := range turns {
		if result {
			countTrue++
		}
	}
	totalTurns := len(turns)
	if totalTurns == 0 {
		return 0.0
	}
	probability := float64(countTrue) / float64(totalTurns)
	return probability
}

func (bb *Agent8) calculateCostInCollectiveImprovement(decisions []bool) float64 {
	var countFalse int
	for _, decision := range decisions {
		if !decision {
			countFalse++
		}
	}

	totalDecisions := len(decisions)
	if totalDecisions == 0 {
		return 0.0
	}

	percentage := (float64(countFalse) / float64(totalDecisions))
	return percentage
}

// PercentageLowEnergyAgents can reflect relect the optimism regarding of the current megabike
func (bb *Agent8) calculatePercentageLowEnergyAgents(energyLevels []float64, threshold float64) float64 {
	var countLowEnergy int

	// Count the number of agents with energy levels below the threshold
	for _, energyLevel := range energyLevels {
		if energyLevel < threshold {
			countLowEnergy++
		}
	}

	// Calculate the percentage
	totalAgents := len(energyLevels)
	if totalAgents == 0 {
		return 0.0
	}

	percentage := (float64(countLowEnergy) / float64(totalAgents))
	return percentage
}

// calculate preference score(preference voting)
func (bb *Agent8) calculatePreferenceScores(proposals map[uuid.UUID]uuid.UUID) map[uuid.UUID]float64 {
	scores := make(map[uuid.UUID]float64)
	currLocation := bb.GetLocation()
	var distances []float64
	distanceToBox := make(map[float64]uuid.UUID)

	// Calculate distances to each loot box and sort them
	for _, proposal := range proposals {
		for _, lootBox := range bb.GetGameState().GetLootBoxes() {
			if lootBox.GetID() == proposal {
				x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
				distance := math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
				distances = append(distances, distance)
				distanceToBox[distance] = proposal
			}
		}
	}
	sort.Float64s(distances)

	// Scoring mechanism
	if bb.GetEnergyLevel() < 0.5 || !bb.hasDesiredColorInRange(proposals, 30) {
		// Score based on distance
		score := float64(len(distances))
		for _, distance := range distances {
			scores[distanceToBox[distance]] = score
			score--
		}
	} else {
		// Score based on color and distance
		score := float64(len(distances))
		for _, distance := range distances {
			lootBoxID := distanceToBox[distance]
			lootBoxColor := bb.GetGameState().GetLootBoxes()[lootBoxID].GetColour()
			if lootBoxColor == bb.GetColour() {
				scores[lootBoxID] = score
				score--
			}
		}
		for _, distance := range distances {
			lootBoxID := distanceToBox[distance]
			if scores[lootBoxID] == 0 { // Only score loot boxes that haven't been scored yet
				scores[lootBoxID] = score
				score--
			}
		}
	}

	return scores
}

// Check color in range:
func (bb *Agent8) hasDesiredColorInRange(proposals map[uuid.UUID]uuid.UUID, rangeThreshold float64) bool {
	currLocation := bb.GetLocation()
	for _, proposal := range proposals {
		for _, lootBox := range bb.GetGameState().GetLootBoxes() {
			if lootBox.GetID() == proposal {
				x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
				distance := math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
				if distance <= rangeThreshold && lootBox.GetColour() == bb.GetColour() {
					return true
				}
			}
		}
	}
	return false
}
