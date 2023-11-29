package team_8

import (
	"SOMAS2023/internal/common/objects"
	"math/rand"

	"SOMAS2023/internal/common/voting"
	"math"

	"sort"

	"github.com/google/uuid"
)

type Agent8 struct {
	*objects.BaseBiker
	energyLevel float64
	color       string
}

type Colour string

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	//TODO： need to be implemented
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

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *Agent8) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	var threshold float64 = 0.2
	decision := make(map[uuid.UUID]bool)
	agentMap := bb.UuidToAgentMap(pendingAgents)

	for uuid, agent := range agentMap {
		var score float64
		if agent.GetColour() == bb.GetColour() {
			score = (agent.GetEnergyLevel() - bb.CalculateAverageEnergy(bb.GetBike())) / bb.CalculateAverageEnergy(bb.GetBike())
		} else {
			score = 0.5 * (agent.GetEnergyLevel() - bb.CalculateAverageEnergy(bb.GetBike())) / bb.CalculateAverageEnergy(bb.GetBike())
		}
		if score >= threshold {
			decision[uuid] = true
		} else {
			decision[uuid] = false
		}

	}

	return decision
}

// decide which bike to go to
// for now it just returns a random uuid

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

func (bb *Agent8) ChangeBike() uuid.UUID {
	// Get all the bikes from the game state
	megaBikes := bb.GetGameState().GetMegaBikes()

	// Initialize a map to store Borda scores for each bike
	bordaScores := make(map[uuid.UUID]float64)

	// Iterate through each bike
	for bikeID, megabike := range megaBikes {
		// Calculate the Borda score for the current bike
		bordaScore := bb.CalculateAverageEnergy(bikeID) +
			float64(bb.CountAgentsWithSameColour(bikeID)) +
			CalculateGiniIndexFromAB(float64(bb.CountAgentsWithSameColour(bikeID)), float64(len(megabike.GetAgents())))

		// Store the Borda score in the map
		bordaScores[bikeID] = bordaScore
	}

	// Find the bike with the highest Borda score
	var highestBordaScore float64
	var winningBikeID uuid.UUID
	for bikeID, score := range bordaScores {
		if score > highestBordaScore {
			highestBordaScore = score
			winningBikeID = bikeID
		}
	}

	return winningBikeID
}

// default implementation returns the id of the nearest lootbox
// Alex
func (bb *Agent8) ProposeDirection() uuid.UUID {
	return uuid.New()
}

// DecideAction helper functions

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

func (bb *Agent8) calculateValueJudgement(utilityLevels []float64, agentGoals []int, targetGoal int, turns []bool) float64 {
	/* Example usage
	utilityLevels := []float64{80.0, 90.0, 75.0, 85.0}
	agentGoals := []int{1, 2, 1, 1, 2, 2, 1, 1, 2, 2}
	targetGoal := 1
	turns := []bool{true, false, true, true, false, true, true, false, true}*/
	averageUtility := bb.calculateAverageUtility(utilityLevels)
	percentageSameGoal := bb.calculatePercentageSameGoal(agentGoals, targetGoal)
	probabilitySatisfiedLoops := bb.calculateProbabilitySatisfiedLoops(turns)

	// Calculate the average score
	averageScore := (averageUtility + percentageSameGoal + probabilitySatisfiedLoops) / 3
	return averageScore
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

func (bb *Agent8) calculatePercentageLowEnergyAgents(energyLevels []int, threshold int) float64 {
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

func (bb *Agent8) calculateAverageOfCostAndPercentage(decisions []bool, energyLevels []int, threshold int) float64 {

	/*Example usage
	decisions := []bool{true, false, false, true, true, false, false, true}
	energyLevels := []int{80, 45, 60, 30, 70, 40, 55, 75, 90}
	threshold := 50

	// Calculate the average of values returned by calculateCostInCollectiveImprovement and calculatePercentageLowEnergyAgents
	averageResult := calculateAverageOfCostAndPercentage(decisions, energyLevels, threshold)*/
	costPercentage := bb.calculateCostInCollectiveImprovement(decisions)
	percentageLowEnergy := bb.calculatePercentageLowEnergyAgents(energyLevels, threshold)

	// Calculate the average
	averageResult := (costPercentage + percentageLowEnergy) / 2
	return averageResult
}

func (bb *Agent8) DecideAction() objects.BikerAction {
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
	valueJudgement := bb.calculateValueJudgement(utilityLevels, agentGoals, targetGoal, turns)

	// Scale the ‘Cost in the collective improvement’
	AverageOfCost := bb.calculateAverageOfCostAndPercentage(decisions, energyLevels, EnergyThreshold)

	// Find the overall ‘changeBike’ coefficient
	changeBikeCoefficient := 0.6*valueJudgement - 0.4*AverageOfCost

	// Make a decision based on the calculated coefficients
	// rand.Float64() to be deicided
	if rand.Float64() > changeBikeCoefficient {
		return objects.ChangeBike
	}

	// Default action
	return objects.Pedal
}

// //  Nemo started

// calculate preference score(preference voting)
func (bb *Agent8) calculatePreferenceScores(proposals []uuid.UUID) map[uuid.UUID]float64 {
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
	if bb.GetEnergyLevel() < 0.3 || !bb.hasDesiredColorInRange(proposals, 30) {
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
func (bb *Agent8) hasDesiredColorInRange(proposals []uuid.UUID, rangeThreshold float64) bool {
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

// Multi-voting system
func (bb *Agent8) FinalDirectionVote(proposals []uuid.UUID, overallScores voting.LootboxVoteMap) voting.LootboxVoteMap {
	// Calculate the biker's individual preference scores
	preferenceScores := bb.calculatePreferenceScores(proposals)

	combinedScores := make(map[uuid.UUID]float64)
	for _, proposal := range proposals {
		combinedScore := preferenceScores[proposal] + overallScores[proposal]
		combinedScores[proposal] = combinedScore
	}
	softmaxScores := softmax(combinedScores)

	return softmaxScores
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	//TODO： need to be implemented
	distribution := make(map[uuid.UUID]float64)
	return distribution

}
