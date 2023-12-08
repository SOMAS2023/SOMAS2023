package team8

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

func (bb *Agent8) GetAverageReputation(agent objects.IBaseBiker) float64 {
	averageReputation := 0.0
	agentNum := 0
	for _, reputation := range agent.GetReputation() {
		averageReputation += reputation
		if reputation != 0 {
			agentNum++
		}
	}
	return averageReputation / float64(agentNum)
}
