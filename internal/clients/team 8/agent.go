package main

import (
	"math/rand"
)

// BikerAction represents the possible actions a biker can take
type BikerAction int

const (
	Pedal BikerAction = iota
	ChangeBike
	// Add more actions as needed
)

// BaseBiker represents the base structure for a biker agent
type BaseBiker struct {
}

// calculateAverageUtilityPercentage calculates the average of utility levels and returns the percentage
func calculateAverageUtilityPercentage(utilityLevels []float64) float64 {
	var sum float64
	for _, value := range utilityLevels {
		sum += value
	}
	average := sum / float64(len(utilityLevels))
	percentage := (average / 100) * 100
	return percentage
}

// calculatePercentageSameGoal calculates the percentage of agents with the same goal
func calculatePercentageSameGoal(agentGoals []int, targetGoal int) float64 {
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
	percentage := (float64(countSameGoal) / float64(totalAgents)) * 100.0
	return percentage
}

// calculateProbabilitySatisfiedLoops calculates the probability of having 'true' in the array
func calculateProbabilitySatisfiedLoops(turns []bool) float64 {
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

func calculateValueJudgement(utilityLevels []float64, agentGoals []int, targetGoal int, turns []bool) float64 {
	/* Example usage
	utilityLevels := []float64{80.0, 90.0, 75.0, 85.0}
	agentGoals := []int{1, 2, 1, 1, 2, 2, 1, 1, 2, 2}
	targetGoal := 1
	turns := []bool{true, false, true, true, false, true, true, false, true}*/
	averageUtilityPercentage := calculateAverageUtilityPercentage(utilityLevels)
	percentageSameGoal := calculatePercentageSameGoal(agentGoals, targetGoal)
	probabilitySatisfiedLoops := calculateProbabilitySatisfiedLoops(turns)

	// Calculate the average score
	averageScore := (averageUtilityPercentage + percentageSameGoal + probabilitySatisfiedLoops) / 3
	return averageScore
}

func calculateCostInCollectiveImprovement(decisions []bool) float64 {
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

	percentage := (float64(countFalse) / float64(totalDecisions)) * 100.0
	return percentage
}

func calculatePercentageLowEnergyAgents(energyLevels []int, threshold int) float64 {
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

	percentage := (float64(countLowEnergy) / float64(totalAgents)) * 100.0
	return percentage
}

func calculateAverageOfCostAndPercentage(decisions []bool, energyLevels []int, threshold int) float64 {

	/*Example usage
	decisions := []bool{true, false, false, true, true, false, false, true}
	energyLevels := []int{80, 45, 60, 30, 70, 40, 55, 75, 90}
	threshold := 50

	// Calculate the average of values returned by calculateCostInCollectiveImprovement and calculatePercentageLowEnergyAgents
	averageResult := calculateAverageOfCostAndPercentage(decisions, energyLevels, threshold)*/
	costPercentage := calculateCostInCollectiveImprovement(decisions)
	percentageLowEnergy := calculatePercentageLowEnergyAgents(energyLevels, threshold)

	// Calculate the average
	averageResult := (costPercentage + percentageLowEnergy) / 2
	return averageResult
}

/*
	func calculateFairnessIndex(bb *BaseBiker) float64 {
		// Implement the logic to calculate fairness index
		return 0.0
	}
*/

func (bb *BaseBiker) DecideAction() BikerAction {
	// Example usage
	utilityLevels := []float64{80.0, 90.0, 75.0, 85.0, 45.0, 35.0, 60.0, 70.0, 65.0}
	agentGoals := []int{1, 2, 1, 1, 2, 2, 1, 1, 2}
	targetGoal := 1
	turns := []bool{true, false, true, true, false, true, true, false, true}
	decisions := []bool{true, false, false, true, true, false, false, true, false}
	energyLevels := []int{80, 45, 60, 30, 70, 40, 55, 75, 90}
	energyThreshold := 50

	// Find quantified ‘Value-judgement’
	valueJudgement := calculateValueJudgement(utilityLevels, agentGoals, targetGoal, turns)

	// Scale the ‘Cost in the collective improvement’
	AverageOfCostAndPercentage := calculateAverageOfCostAndPercentage(decisions, energyLevels, energyThreshold)

	// Find the overall ‘changeBike’ coefficient
	changeBikeCoefficient := 0.6*valueJudgement - 0.4*AverageOfCostAndPercentage

	// Make a decision based on the calculated coefficients
	if rand.Float64() < changeBikeCoefficient {
		return ChangeBike
	}

	// Default action
	return Pedal
}

func main() {
	// test
}
