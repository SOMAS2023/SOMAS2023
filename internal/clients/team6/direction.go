package team6

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"math"
	"sort"

	"github.com/google/uuid"
)

// ============================================= Propose a direction =============================================
func (bb *Team6Biker) ProposeDirection() uuid.UUID {
	nearestSameColourBox := bb.nearestSameColourLoot(bb.nearestLootList()) // Get the nearest lootbox of the same colour as the biker.
	nearestBox := bb.nearestLoot(bb.nearestLootList())                     // Get the nearest lootbox of any colour

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

func (bb *Team6Biker) nearestLootList() []objects.ILootBox {
	currLocation := bb.GetLocation()
	var currDist float64
	var lbList []objects.ILootBox
	var disList []float64
	var sortedLoot []objects.ILootBox
	idToDis := make(map[objects.ILootBox]float64)

	for _, loot := range bb.GetGameState().GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		lbList = append(lbList, loot)
		disList = append(disList, currDist)
	}

	for i := 0; i < len(lbList) && i < len(disList); i++ {
		idToDis[lbList[i]] = disList[i]
	}

	for id := range idToDis {
		sortedLoot = append(sortedLoot, id)
	}

	sort.Slice(sortedLoot, func(i, j int) bool {
		return idToDis[sortedLoot[i]] < idToDis[sortedLoot[j]]
	})

	return sortedLoot
}

func (bb *Team6Biker) nearestLoot(sortedLoot []objects.ILootBox) uuid.UUID {
	var loot2Awdi float64
	var nearestBoxPos utils.Coordinates
	var nearestBox uuid.UUID
	awdiPos := bb.GetGameState().GetAwdi().GetPosition()

	for _, loot := range sortedLoot {
		nearestBoxPos = loot.GetPosition()
		loot2Awdi = math.Sqrt(math.Pow(awdiPos.X-nearestBoxPos.X, 2) + math.Pow(awdiPos.Y-nearestBoxPos.Y, 2))
		if loot2Awdi > utils.CollisionThreshold {
			nearestBox = loot.GetID()
			break
		}
	}
	return nearestBox
}

// Derive the UUID of the nearest lootbox of the same colour as the biker
func (bb *Team6Biker) nearestSameColourLoot(sortedLoot []objects.ILootBox) uuid.UUID {
	var loot2Awdi float64
	var nearestBoxPos utils.Coordinates
	var nearestSameColour uuid.UUID
	awdiPos := bb.GetGameState().GetAwdi().GetPosition()

	for _, loot := range sortedLoot {
		nearestBoxPos = loot.GetPosition()
		loot2Awdi = math.Sqrt(math.Pow(awdiPos.X-nearestBoxPos.X, 2) + math.Pow(awdiPos.Y-nearestBoxPos.Y, 2))
		if loot2Awdi > utils.CollisionThreshold {
			if loot.GetColour() == bb.GetColour() {
				nearestSameColour = loot.GetID()
				break
			}
		}
	}
	return nearestSameColour

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

func (bb *Team6Biker) ProposedVote() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	allLoot := bb.GetGameState().GetLootBoxes()
	nearestLootList := bb.nearestLootList()
	var i = len(nearestLootList)
	var k = 0
	var fl = 0 // I don't understand why
	for _, loot := range allLoot {
		fl = 0
		for _, j := range nearestLootList {
			if j == loot {
				fl = 1
				if j.GetColour() == bb.GetColour() {
					votes[loot.GetID()] = float64(i) * 2.0
					k += i
					//votes.append(i*1.2)
				} else {
					votes[loot.GetID()] = float64(i)
					//votes.append(i)
				}
				k += i
			}
		}
		if fl == 0 {
			votes[loot.GetID()] = 0.0
		}
		i -= 1
	}
	//softmax(voteRulerMap)
	for _, vote := range votes {
		vote = vote / float64(k)
	}
	return votes
}
