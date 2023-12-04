package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math/rand"

	"SOMAS2023/internal/common/voting"
	"math"

	"sort"

	"github.com/google/uuid"
)

type GP struct {
	EnergyThreshold              float64
	DistanceThresholdForVoting   float64
	ThresholdForJoiningDecision  float64
	ThresholdForChangingMegabike float64
}

var GlobalParameters GP = GP{EnergyThreshold: 0.5, DistanceThresholdForVoting: 30, ThresholdForJoiningDecision: 0.2, ThresholdForChangingMegabike: 0.3}

type Agent8 struct {
	*objects.BaseBiker
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DecideGovernance <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// base biker defaults to democracy
func (bb *Agent8) DecideGovernance() voting.GovernanceVote {
	// TODO: implement this function
	governanceRanking := make(voting.GovernanceVote)
	governanceRanking[utils.Democracy] = 1.0
	governanceRanking[utils.Dictatorship] = 0.0
	governanceRanking[utils.Leadership] = 0.0
	return governanceRanking
}

// defaults to voting for first agent in the list
func (bb *Agent8) VoteDictator() voting.IdVoteMap {
	// TODO: implement this function
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()
	for i, fellowBiker := range fellowBikers {
		if i == 0 {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	return votes
}

// defaults to voting for first agent in the list
func (bb *Agent8) VoteLeader() voting.IdVoteMap {
	// TODO: implement this function
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()
	for i, fellowBiker := range fellowBikers {
		if i == 0 {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	return votes
}

//===============================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 1 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) VoteForKickout() map[uuid.UUID]int {
	// TODO: implement this function
	voteResults := make(map[uuid.UUID]int)
	bikeID := bb.GetBike()

	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if agentID != bb.GetID() {
			// random votes to other agents
			voteResults[agentID] = rand.Intn(2) // randomly assigns 0 or 1 vote
		}
	}

	return voteResults
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
	threshold := GlobalParameters.ThresholdForJoiningDecision
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

//===============================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 2 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) DecideAction() objects.BikerAction {

	var energyLevels []float64
	var target_goal int
	energy_threshold := GlobalParameters.EnergyThreshold
	changingbike_threshold := GlobalParameters.ThresholdForChangingMegabike

	// utility should be calculated by the fomula outlined on page7 of Lec6
	utilityLevels := []float64{80.0, 90.0, 75.0, 85.0, 45.0, 35.0, 60.0, 70.0, 65.0}
	turns := []bool{true, false, true, true, false, true, true, false, true}
	decisions := []bool{true, false, false, true, true, false, false, true, false}

	// get the energy level of all agents in the megabike
	fellowBikers := bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetAgents()

	for _, agent := range fellowBikers {
		if !(agent.GetID() == bb.GetID()) {
			energy_level := agent.GetEnergyLevel()
			energyLevels = append(energyLevels, energy_level)
		}
	}

	goalPreferenceList := make([]int, len(energyLevels))
	if bb.GetEnergyLevel() >= energy_threshold {
		target_goal = 1
	} else {
		target_goal = 0
	}

	// Convert energyLevels to 0 or 1 based on the threshold
	for i, energy := range energyLevels {
		if energy >= energy_threshold {
			goalPreferenceList[i] = 1
		} else {
			goalPreferenceList[i] = 0
		}
	}

	// Find quantified ‘Value-judgement’
	valueJudgement := bb.calculateValueJudgement(utilityLevels, goalPreferenceList, target_goal, turns)

	// Scale the ‘Cost in the collective improvement’
	AverageOfCost := bb.calculateAverageOfCostAndPercentage(decisions, energyLevels, energy_threshold)

	// Find the overall ‘changeBike’ coefficient
	changeBikeCoefficient := 0.6*valueJudgement - 0.4*AverageOfCost

	// Make a decision based on the calculated coefficients
	if changingbike_threshold > changeBikeCoefficient {
		return objects.ChangeBike
	}

	// Default action
	return objects.Pedal
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

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 3 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) ProposeDirection() uuid.UUID {
	lootBoxes := bb.GetGameState().GetLootBoxes()
	preferences := make(map[uuid.UUID]float64)
	softmaxPreferences := make(map[uuid.UUID]float64)

	// Calculate preferences
	for _, lootBox := range lootBoxes {
		distance := calculateDistance(bb.GetLocation(), lootBox.GetPosition())
		colorPreference := calculateColorPreference(bb.GetColour(), lootBox.GetColour())
		energyWeighting := calculateEnergyWeighting(bb.GetEnergyLevel())

		preferences[lootBox.GetID()] = colorPreference + (GlobalParameters.DistanceThresholdForVoting-distance)*energyWeighting
	}

	// Apply softmax to convert preferences to a probability distribution
	softmaxPreferences = softmax(preferences)

	// Rank loot boxes based on preferences
	rankedLootBoxes := rankByPreference(softmaxPreferences)

	return rankedLootBoxes[0]
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 4 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) DictateDirection() uuid.UUID {
	// TODO: implement this function
	return uuid.Nil
}

func (bb *Agent8) LeadDirection() uuid.UUID {
	// TODO: implement this function
	return uuid.Nil
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

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 5 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	// TODO: implement this function
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 6 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	// TODO: implement this function
	distribution := make(map[uuid.UUID]float64)
	return distribution

}

// =========================================================================================================================================================
