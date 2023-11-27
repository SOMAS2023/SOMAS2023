package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
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
	lootBoxes := bb.GetGameState().GetLootBoxes()
	preferences := make(map[uuid.UUID]float64)
	softmaxPreferences := make(map[uuid.UUID]float64)

	// Calculate preferences
	for _, lootBox := range lootBoxes {
		distance := calculateDistance(bb.GetLocation(), lootBox.GetPosition())
		colorPreference := calculateColorPreference(bb.GetColour(), lootBox.GetColour())
		energyWeighting := calculateEnergyWeighting(bb.GetEnergyLevel())

		preferences[lootBox.GetID()] = colorPreference + distance*energyWeighting
	}

	// Apply softmax to convert preferences to a probability distribution
	softmaxPreferences = softmax(preferences)

	// Rank loot boxes based on preferences
	rankedLootBoxes := rankByPreference(softmaxPreferences)

	// Select the top choice(s) based on ranking
	selectedLootBox := selectTopChoices(rankedLootBoxes, bb.GetGameState().GetVotingListLength())

	// Consider social dilemma factors if applicable
	finalSelection := considerSocialDilemma(selectedLootBox, bb.GetGameState())

	// If the biker is a leader, adjust the final selection accordingly
	//if bb.isLeader() {
	//	finalSelection = leaderStrategyAdjustment(finalSelection, bb.gameState)
	//}

	return finalSelection
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

// softmax applies the softmax function to the preferences to get a probability distribution
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

// considerSocialDilemma adjusts the selection based on social factors
func considerSocialDilemma(selected uuid.UUID, gameState IGameState) uuid.UUID {
	// Placeholder: The actual implementation would depend on the social dilemma factors considered in the game
	return selected
}

// leaderStrategyAdjustment adjusts the selection if the agent is a leader
func leaderStrategyAdjustment(selected uuid.UUID, gameState IGameState) uuid.UUID {
	// Placeholder: The actual implementation would depend on the leader's strategy in the game
	return selected
}

// //  Nemo started

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *Agent8) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	votes := make(voting.LootboxVoteMap)

	distance_threshold := 30
	energy_threshold := 30
	currLocation := bb.GetLocation()
	var chosenBox uuid.UUID
	minDistance := math.MaxFloat64
	foundDesiredColor := false

	// Check for desired color within the distance range
	for _, proposal := range proposals {
		for _, lootBox := range bb.GetGameState().GetLootBoxes() {
			if lootBox.GetID() == proposal {
				// x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
				// // distance := calculateDistance(bb.GetLocation(), lootBox.GetPosition())
				// distance := math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
				distance := calculateDistance(bb.GetLocation(), lootBox.GetPosition())

				if distance <= float64(distance_threshold) && lootBox.GetColour() == bb.GetColour() {
					if distance < minDistance {
						minDistance = distance
						chosenBox = proposal
						foundDesiredColor = true
					}
				}
			}
		}
	}

	// If no desired color found within range or energy level is low, vote for the closest loot box
	if !foundDesiredColor || bb.GetEnergyLevel() < float64(energy_threshold) {
		for _, proposal := range proposals {
			for _, lootBox := range bb.GetGameState().GetLootBoxes() {
				if lootBox.GetID() == proposal {
					x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
					distance := math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))

					if distance < minDistance {
						minDistance = distance
						chosenBox = proposal
					}
				}
			}
		}
	}

	// Cast vote
	for _, proposal := range proposals {
		if proposal == chosenBox {
			votes[proposal] = 1.0 // Full vote for the chosen box
		} else {
			votes[proposal] = 0.0 // No vote for the other boxes
		}
	}

	return votes
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	//TODO： need to be implemented
	distribution := make(map[uuid.UUID]float64)
	return distribution
}
