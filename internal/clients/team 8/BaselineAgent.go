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

// in the MVP the biker's action defaults to pedaling (as it won't be able to change bikes)
// in future implementations this function will be overridden by the agent's specific strategy
// which will be used to determine whether to pedalor try to change bike
func (bb *Agent8) DecideAction() objects.BikerAction {
	//TODO： need to be implemented
	return 0
}

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	//TODO： need to be implemented
}

// AgentRanking holds an agent's UUID and their energy level
type AgentRanking struct {
	ID          uuid.UUID
	energyLevel float64
}

func GetAgentByUUID(id uuid.UUID) *Agent8 {
	// Implement the logic to retrieve Agent8 instance
	// This is just a placeholder
	return &Agent8{}
}

// GetEnergyLevel returns the energy level of the agent.
func (bb *Agent8) GetEnergyLevel() float64 {
	return bb.energyLevel
}

func (bb *Agent8) GetColor() string {
	return bb.color
}

type AgentBordaRank struct {
	ID         uuid.UUID
	BordaPoint float64
}

func RankAgentsWithEnergy(pendingAgents []uuid.UUID, weight float64) []AgentBordaRank {
	type agentWithEnergy struct {
		ID          uuid.UUID
		EnergyLevel float64
	}

	var agentsWithEnergy []agentWithEnergy
	for _, agentID := range pendingAgents {
		agent := GetAgentByUUID(agentID)
		agentsWithEnergy = append(agentsWithEnergy, agentWithEnergy{ID: agentID, EnergyLevel: agent.energyLevel})
	}

	// Sorting agents by energy level in descending order
	sort.Slice(agentsWithEnergy, func(i, j int) bool {
		return agentsWithEnergy[i].EnergyLevel > agentsWithEnergy[j].EnergyLevel
	})

	var AgentsScore []AgentBordaRank
	totalAgents := len(agentsWithEnergy)
	for i, agent := range agentsWithEnergy {
		bordaPoint := float64(totalAgents-i) * weight // Calculate Borda point
		AgentsScore = append(AgentsScore, AgentBordaRank{ID: agent.ID, BordaPoint: bordaPoint})
	}

	return AgentsScore
}

func SumBordaScores(canons ...[]AgentBordaRank) map[uuid.UUID]float64 {
	summedScores := make(map[uuid.UUID]float64)

	// Iterate over each canon
	for _, canon := range canons {
		// Iterate over each agent's Borda score in the canon
		for _, agentRank := range canon {
			// Sum the scores for each agent
			summedScores[agentRank.ID] += agentRank.BordaPoint
		}
	}

	return summedScores
}

func AssignPointsBasedOnColor(agentsMap map[uuid.UUID]*Agent8, ourAgentColor string, weighting float64) []AgentBordaRank {
	var ranks []AgentBordaRank
	totalAgents := float64(len(agentsMap))
	for id, agent := range agentsMap {
		points := 0.0
		if agent.color == ourAgentColor {
			points = totalAgents * weighting / 2 // Assign points if the color matches
		}
		ranks = append(ranks, AgentBordaRank{ID: id, BordaPoint: points})
	}
	return ranks
}

func (bb *Agent8) CalculateSingleAgentScore(agentID uuid.UUID, threshold float64) map[uuid.UUID]bool {
	agent := GetAgentByUUID(agentID)

	colorValue := 0.0
	energyScore := 0.0
	if agent.color == bb.color {
		colorValue = 1.0
	}
	if agent.energyLevel > bb.energyLevel {
		energyScore = 1.0
	}

	totalScore := energyScore + colorValue

	decisions := make(map[uuid.UUID]bool)
	decisions[agentID] = totalScore >= threshold

	return decisions
}

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *Agent8) DecideJoining(pendingAgents []uuid.UUID,

// n int, // number of agents we want to allow to join the bike
) map[uuid.UUID]bool {

	decisions := make(map[uuid.UUID]bool)

	if len(pendingAgents) == 0 {
		agentID := uuid.UUID{}
		decisions := bb.CalculateSingleAgentScore(agentID, 1)

		return decisions

	} else {
		agentsMap := make(map[uuid.UUID]*Agent8)

		// Functions for various canons
		// Weighting for each canon should be confirmed afterwards
		AgentsScore_energy := RankAgentsWithEnergy(pendingAgents, 0.3)
		ownAgentColor := bb.color
		AgentScore_color := AssignPointsBasedOnColor(agentsMap, ownAgentColor, 0.7)
		// AgentsScore_reputation := ()

		// sum weighted scores from different canons
		summedScores := SumBordaScores(AgentsScore_energy, AgentScore_color)

		// rank agents and select the best agent
		var summedRankings []AgentBordaRank
		for id, score := range summedScores {
			summedRankings = append(summedRankings, AgentBordaRank{ID: id, BordaPoint: score})
		}

		sort.Slice(summedRankings, func(i, j int) bool {
			return summedRankings[i].BordaPoint > summedRankings[j].BordaPoint
		})

		for i, agentRank := range summedRankings {
			decisions[agentRank.ID] = i < 1 // this should be the n from input afterwards
		}
	}
	return decisions
}

// decide which bike to go to
// for now it just returns a random uuid
func (bb *Agent8) ChangeBike() uuid.UUID {
	//TODO： need to be implemented
	return uuid.New()
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
	if bb.GetEnergyLevel() < 30 || !bb.hasDesiredColorInRange(proposals, 30) {
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
