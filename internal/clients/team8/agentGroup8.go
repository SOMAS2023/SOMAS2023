package team_8

import (
	"SOMAS2023/internal/common/objects"

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

func UuidToAgentMap(pendingAgents []uuid.UUID, megaBikes map[uuid.UUID]objects.IMegaBike) map[uuid.UUID]objects.IBaseBiker {
	agentMap := make(map[uuid.UUID]objects.IBaseBiker)

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
	agentMap := UuidToAgentMap(pendingAgents, bb.GetGameState().GetMegaBikes())

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

// //  Nemo started

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
// NB: One vote system

// func (bb *Agent8) FinalDirectionVote(proposals []uuid.UUID, currentVotes voting.LootboxVoteMap) voting.LootboxVoteMap {
// 	votes := make(voting.LootboxVoteMap)

// 	distance_threshold := 30
// 	energy_threshold := 30
// 	currLocation := bb.GetLocation()
// 	var chosenBox uuid.UUID
// 	minDistance := math.MaxFloat64
// 	foundDesiredColor := false

// 	// Check for desired color within the distance range
// 	for _, proposal := range proposals {
// 		for _, lootBox := range bb.GetGameState().GetLootBoxes() {
// 			if lootBox.GetID() == proposal {
// 				// x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
// 				// // distance := calculateDistance(bb.GetLocation(), lootBox.GetPosition())
// 				// distance := math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
// 				distance := calculateDistance(bb.GetLocation(), lootBox.GetPosition())

// 				if distance <= float64(distance_threshold) && lootBox.GetColour() == bb.GetColour() {
// 					if distance < minDistance {
// 						minDistance = distance
// 						chosenBox = proposal
// 						foundDesiredColor = true
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// If no desired color found within range or energy level is low, vote for the closest loot box
// 	if !foundDesiredColor || bb.GetEnergyLevel() < float64(energy_threshold) {
// 		for _, proposal := range proposals {
// 			for _, lootBox := range bb.GetGameState().GetLootBoxes() {
// 				if lootBox.GetID() == proposal {
// 					x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
// 					distance := math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))

// 					if distance < minDistance {
// 						minDistance = distance
// 						chosenBox = proposal
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// // Cast vote
// 	// for _, proposal := range proposals {
// 	// 	if proposal == chosenBox {
// 	// 		votes[proposal] = 1.0 // Full vote for the chosen box
// 	// 	} else {
// 	// 		votes[proposal] = 0.0 // No vote for the other boxes
// 	// 	}
// 	// }

// 	//NB!!!
// 	// these above can be possibly replaced by previous storage of preference ***
// 	totalVoters := len(currentVotes)                              // Assuming you have a way to get the total number of voters
// 	leadingOption, secondOption := getTopTwoOptions(currentVotes) // Implement this method to find top two voted options

// 	// Decide on strategic voting
// 	if chosenBox != leadingOption && chosenBox != secondOption && currentVotes[leadingOption]-currentVotes[secondOption] > float64(totalVoters)/4 {
// 		// Vote for the leading option if the preferred option is not leading and the lead is significant
// 		votes[leadingOption] = 1.0
// 	} else if chosenBox == secondOption && currentVotes[leadingOption]-currentVotes[secondOption] < float64(totalVoters)/3 {
// 		votes[secondOption] = 1.0
// 	} else {
// 		// Otherwise, vote for the preferred option
// 		votes[chosenBox] = 1.0
// 	}

// 	return votes
// }

// for 1 voting system
// func getTopTwoOptions(currentVotes voting.LootboxVoteMap) (uuid.UUID, uuid.UUID) {
// 	var maxVoteID, secondMaxVoteID uuid.UUID
// 	maxVote, secondMaxVote := -1.0, -1.0

// 	for id, votes := range currentVotes {
// 		if votes > maxVote {
// 			// Update second max
// 			secondMaxVoteID = maxVoteID
// 			secondMaxVote = maxVote
// 			// Update max
// 			maxVoteID = id
// 			maxVote = votes
// 		} else if votes > secondMaxVote {
// 			secondMaxVoteID = id
// 			secondMaxVote = votes
// 		}
// 	}

// 	return maxVoteID, secondMaxVoteID
// }

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
	votes := make(voting.LootboxVoteMap)

	// Calculate the biker's individual preference scores
	preferenceScores := bb.calculatePreferenceScores(proposals)

	// Combine individual preferences with overall scores
	combinedScores := make(map[uuid.UUID]float64)
	for _, proposal := range proposals {
		combinedScore := preferenceScores[proposal] + overallScores[proposal]
		combinedScores[proposal] = combinedScore
	}

	// Sort the loot boxes based on the combined score
	sortedBoxes := sortLootBoxesByScore(combinedScores)

	// Assign scores based on ranking (6 for the top, then 5, 4, etc.)
	score := len(proposals)
	for _, boxID := range sortedBoxes {
		votes[boxID] = float64(score)
		score--
	}

	return votes
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

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	//TODO： need to be implemented
	distribution := make(map[uuid.UUID]float64)
	return distribution

}
