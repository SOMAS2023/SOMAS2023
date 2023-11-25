package team_8

import (
	"SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

type Agent8 struct {
	*objects.BaseBiker
	ID           uuid.UUID
	soughtColour utils.Colour
	energyLevel  float64
	gameState    IGameState
}

// type BaseBiker struct {
// 	// Fields from your existing structure
// 	ID               uuid.UUID
// 	soughtColour     string // Assuming utils.Colour is a string type
// 	onBike           bool
// 	energyLevel      float64
// 	points           int
// 	megaBikeId       uuid.UUID
// 	gameState        IGameState
// 	allocationParams ResourceAllocationParams
// }

type Colour string

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

func (bb *Agent8) GetEnergyLevel() float64 {
	return bb.energyLevel
}

// CalculateAverageEnergy calculates the average energy level for agents on a specific bike.
func (bb *Agent8) CalculateAverageEnergy(bikeID uuid.UUID) float64 {
	// Step 1: Get fellowBikers from the specified bike
	fellowBikers := bb.gameState.GetMegaBikes()[bikeID].GetAgents()

	// Step 2: Ensure there is at least one agent
	if len(fellowBikers) == 0 {
		return 0.0 // or handle this case according to your requirements
	}

	// Step 3: Calculate the sum of energy levels
	sum := 0.0
	for _, agent := range fellowBikers {
		sum += agent.energyLevel
	}

	// Step 4: Calculate the average
	average := sum / float64(len(fellowBikers))

	return average
}

func (bb *Agent8) GetColour() utils.Colour {
	return bb.soughtColour
}

// CountAgentsWithSameColour counts the number of agents with the same colour as the reference agent on a specific bike.
func (bb *Agent8) CountAgentsWithSameColour(bikeID uuid.UUID) int {
	// Step 1: Get reference colour from the BaseBiker
	referenceColour := bb.GetColour()

	// Step 2: Get fellowBikers from the specified bike
	fellowBikers := bb.gameState.GetMegaBikes()[bikeID].GetAgents()

	// Step 3: Ensure there is at least one agent
	if len(fellowBikers) == 0 {
		return 0 // or handle this case according to your requirements
	}

	// Step 4: Count agents with the same colour as the reference agent
	count := 0
	for _, agent := range fellowBikers {
		if agent.soughtColour == referenceColour {
			count++
		}
	}

	return count
}

func (bb *Agent8) CalculateBordaScore() uuid.UUID {
	// Get all the bikes from the game state
	megaBikes := bb.gameState.GetMegaBikes()

	// Initialize a map to store Borda scores for each bike
	bordaScores := make(map[uuid.UUID]float64)

	// Iterate through each bike
	for bikeID, megabike := range megaBikes {
		// Calculate the Borda score for the current bike
		bordaScore := bb.CalculateAverageEnergy(bikeID) +
			float64(bb.CountAgentsWithSameColour(bikeID)) +
			CalculateGiniIndexFromAB(float64(bb.CountAgentsWithSameColour(bikeID)), float64(len(megaBikes.GetAgents())))

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
