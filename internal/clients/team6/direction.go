package team6

import (
	voting "SOMAS2023/internal/common/voting"
	"math"

	"github.com/google/uuid"
)

// ============================================= Propose a direction =============================================
func (bb *Team6Biker) ProposeDirection() uuid.UUID {
	nearestSameColourBox := bb.nearestSameColourLoot() // Get the nearest lootbox of the same colour as the biker.
	nearestBox := bb.nearestLoot()                     // Get the nearest lootbox of any colour

	if bb.GetEnergyLevel() < energyThreshold {
		// When the Biker's current energy is below the energy threshold, go to the lootbox.
		return nearestBox
	} else {
		if nearestSameColourBox != uuid.Nil {
			// If the nearest lootbox of the same colour as the biker exists, go to this lootbox
			return nearestSameColourBox
		} else {
			// Go to the nearest lootbox
			return nearestBox
		}
	}
}

// This function is the same as the code in BaseBiker.go
func (bb *Team6Biker) nearestLoot() uuid.UUID {
	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestBox uuid.UUID
	var currDist float64
	for _, loot := range bb.GetGameState().GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return nearestBox
}

// Derive the UUID of the nearest lootbox of the same colour as the biker
func (bb *Team6Biker) nearestSameColourLoot() uuid.UUID {

	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestSameColourBox = bb.nearestLoot()
	var currDist float64
	bikerColour := bb.GetColour()
	for _, loot := range bb.GetGameState().GetLootBoxes() {
		lootColour := loot.GetColour() // Get the colour of the lootbox
		if lootColour == bikerColour {
			x, y := loot.GetPosition().X, loot.GetPosition().Y
			currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
			if currDist < shortestDist {
				nearestSameColourBox = loot.GetID()
				shortestDist = currDist
			}
		}
	}
	return nearestSameColourBox

	//return bb.nearestLoot()
}

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *Team6Biker) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	votes := make(voting.LootboxVoteMap)
	countProposals := make(map[uuid.UUID]int)
	maxCount := 0
	mostCommonProposal := []uuid.UUID{}

	for _, proposal := range proposals {
		countProposals[proposal]++ // Count the occurrance of each proposal
		votes[proposal] = 0.0      // Intialise the weights of each proposal
	}

	// Iterate through the list of proposals
	for proposal, count := range countProposals {
		if count > maxCount {
			// If the current proposal's count is greater than the previous maximum count,
			// update the most common proposal and the maximum count
			mostCommonProposal = []uuid.UUID{proposal}
			maxCount = count
		} else if count == maxCount {
			// If the current proposal's count is equal to the maximum count,
			// add the proposal to the list of most common proposals
			mostCommonProposal = append(mostCommonProposal, proposal)
		}
	}

	// Check if bb's proposal (proposals[bb.GetID()]) is in mostCommonProposal
	proposalID := proposals[bb.GetID()]
	if contains(mostCommonProposal, proposalID) {
		votes[proposalID] = 1.0
	} else {
		// If mostCommonProposal contains only one proposal, give it a weight of 0.5
		if len(mostCommonProposal) == 1 {
			votes[mostCommonProposal[0]] = 0.5
			votes[proposalID] = 0.5
		} else {
			// If mostCommonProposal contains more than one proposal, find the nearest proposal
			currLocation := bb.GetLocation()
			shortestDist := math.MaxFloat64
			var nearestProposal uuid.UUID
			var currDist float64
			for _, proposal := range mostCommonProposal {
				loot := bb.GetGameState().GetLootBoxes()[proposal]
				x, y := loot.GetPosition().X, loot.GetPosition().Y
				currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
				if currDist < shortestDist {
					nearestProposal = proposal
					shortestDist = currDist
				}
			}

			// Set the weight of the nearest proposal to 0.5
			votes[nearestProposal] = 0.5
			votes[proposalID] = 0.5
		}
	}

	return votes
}

func contains(slice []uuid.UUID, item uuid.UUID) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
