package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/voting"

	"math/rand"

	"github.com/google/uuid"
)

type Agent8 struct {
	*objects.BaseBiker
}

// in the MVP the biker's action defaults to pedaling (as it won't be able to change bikes)
// in future implementations this function will be overridden by the agent's specific strategy
// which will be used to determine whether to pedalor try to change bike
/*func (bb *Agent8) DecideAction() objects.BikerAction {
	//see below
	return 0
}*/

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	//TODO： need to be implemented
}

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *Agent8) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	//TODO： need to be implemented
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

// decide which bike to go to
// for now it just returns a random uuid
func (bb *Agent8) ChangeBike() uuid.UUID {
	//TODO： need to be implemented
	return uuid.New()
}

// default implementation returns the id of the nearest lootbox
func (bb *Agent8) ProposeDirection() uuid.UUID {
	//TODO： need to be implemented
	return uuid.New()
}

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *Agent8) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	//TODO： need to be implemented
	votes := make(map[uuid.UUID]float64)
	return votes
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	//TODO： need to be implemented
	distribution := make(map[uuid.UUID]float64)
	return distribution
}

// DecideAction helper functions
type BikerAction int

const (
	Pedal BikerAction = iota
	ChangeBike
	// Add more actions as needed
)

// calculateAverageUtilityPercentage calculates the average of utility levels and returns the percentage
// the utilitylevels need additional parameters to calculate
func calculateAverageUtility(utilityLevels []float64) float64 {
	var sum float64
	for _, value := range utilityLevels {
		sum += value
	}
	average_utility := sum / float64(len(utilityLevels))
	return average_utility
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
	percentage := (float64(countSameGoal) / float64(totalAgents))
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
	averageUtility := calculateAverageUtility(utilityLevels)
	percentageSameGoal := calculatePercentageSameGoal(agentGoals, targetGoal)
	probabilitySatisfiedLoops := calculateProbabilitySatisfiedLoops(turns)

	// Calculate the average score
	averageScore := (averageUtility + percentageSameGoal + probabilitySatisfiedLoops) / 3
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

	percentage := (float64(countFalse) / float64(totalDecisions))
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

	percentage := (float64(countLowEnergy) / float64(totalAgents))
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

func (bb *Agent8) DecideAction() BikerAction {
	// Example usage: assume the game has run 9 iteration
	// requires a set of instances of Agent8
	utilityLevels := []float64{80.0, 90.0, 75.0, 85.0, 45.0, 35.0, 60.0, 70.0, 65.0}
	agentGoals := []int{1, 2, 1, 1, 2, 2, 1, 1, 2}
	targetGoal := 1
	turns := []bool{true, false, true, true, false, true, true, false, true}
	decisions := []bool{true, false, false, true, true, false, false, true, false}
	energyLevels := []int{80, 45, 60, 30, 70, 40, 55, 75, 90}
	EnergyThreshold := 50

	// Find quantified ‘Value-judgement’
	valueJudgement := calculateValueJudgement(utilityLevels, agentGoals, targetGoal, turns)

	// Scale the ‘Cost in the collective improvement’
	AverageOfCost := calculateAverageOfCostAndPercentage(decisions, energyLevels, EnergyThreshold)

	// Find the overall ‘changeBike’ coefficient
	changeBikeCoefficient := 0.6*valueJudgement - 0.4*AverageOfCost

	// Make a decision based on the calculated coefficients
	// rand.Float64() to be deicided
	if rand.Float64() > changeBikeCoefficient {
		return ChangeBike
	}

	// Default action
	return Pedal
}
