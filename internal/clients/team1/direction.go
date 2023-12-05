// DIRECTION DECISION FUNCTIONS

package team1

import (
	"SOMAS2023/internal/common/physics"
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"math"
	"github.com/google/uuid"
)
// ---------------DIRECTION DECISION FUNCTIONS------------------

// Simulates a step of the game, assuming all bikers pedal with the same force as us.
// Returns the distance travelled and the remaining energy
func (bb *Biker1) simulateGameStep(energy float64, velocity float64, force float64) (float64, float64) {
	bikerNum := len(bb.GetFellowBikers())
	totalBikerForce := force * float64(len(bb.GetFellowBikers()))
	totalMass := utils.MassBike + float64(bikerNum)*utils.MassBiker
	acceleration := physics.CalcAcceleration(totalBikerForce, totalMass, velocity)
	distance := velocity + 0.5*acceleration
	energy = energy - force*utils.MovingDepletion
	return distance, energy
}

// Calculates the approximate distance that can be travelled with the given energy
func (bb *Biker1) energyToReachableDistance(energy float64) float64 {
	distance := 0.0
	totalDistance := 0.0
	remainingEnergy := energy
	for remainingEnergy > 0 {
		distance, remainingEnergy = bb.simulateGameStep(remainingEnergy, bb.GetBikeInstance().GetVelocity(), bb.getPedalForce())
		totalDistance = totalDistance + distance
	}
	return totalDistance
}

// Calculates the energy remaining after travelling the given distance
func (bb *Biker1) distanceToEnergy(distance float64, initialEnergy float64) float64 {
	totalDistance := 0.0
	remainingEnergy := initialEnergy
	extraDist := 0.0
	for totalDistance < distance {
		extraDist, remainingEnergy = bb.simulateGameStep(remainingEnergy, bb.GetBikeInstance().GetPhysicalState().Mass, utils.BikerMaxForce*remainingEnergy)
		totalDistance = totalDistance + extraDist
	}

	return remainingEnergy
}

// Finds all boxes within our reachable distance
func (bb *Biker1) getAllReachableBoxes() []uuid.UUID {
	currLocation := bb.GetLocation()
	ourEnergy := bb.GetEnergyLevel()
	lootBoxes := bb.GetGameState().GetLootBoxes()
	reachableBoxes := make([]uuid.UUID, 0)
	var currDist float64
	for _, loot := range lootBoxes {
		lootPos := loot.GetPosition()
		currDist = physics.ComputeDistance(currLocation, lootPos)
		if currDist < bb.energyToReachableDistance(ourEnergy) {
			reachableBoxes = append(reachableBoxes, loot.GetID())
		}
	}
	return reachableBoxes
}

// Checks whether a box of the desired colour is within our reachable distance from a given box
func (bb *Biker1) checkBoxNearColour(box uuid.UUID, energy float64) bool {
	lootBoxes := bb.GetGameState().GetLootBoxes()
	boxPos := lootBoxes[box].GetPosition()
	var currDist float64
	for _, loot := range lootBoxes {
		lootPos := loot.GetPosition()
		currDist = physics.ComputeDistance(boxPos, lootPos)
		if currDist < bb.energyToReachableDistance(energy) && loot.GetColour() == bb.GetColour() {
			return true
		}
	}
	return false
}

// returns the nearest lootbox with respect to the agent's bike current position
// in the MVP this is used to determine the pedalling forces as all agent will be
// aiming to get to the closest lootbox by default
func (bb *Biker1) nearestLoot() uuid.UUID {
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

// Finds the nearest reachable box
func (bb *Biker1) getNearestReachableBox() uuid.UUID {
	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	//default to nearest lootbox
	nearestBox := bb.GetID()
	var currDist float64
	initialized := false
	for id, loot := range bb.GetGameState().GetLootBoxes() {
		if !initialized {
			nearestBox = id
			initialized = true
		}
		lootPos := loot.GetPosition()
		currDist = physics.ComputeDistance(currLocation, lootPos)
		if currDist < shortestDist {
			nearestBox = id
			shortestDist = currDist
		}
	}

	return nearestBox
}

// Finds the nearest lootbox of agent's colour
func (bb *Biker1) nearestLootColour() (uuid.UUID, float64) {
	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	//default to nearest lootbox
	nearestBox := bb.GetID()
	initialized := false
	var currDist float64
	for id, loot := range bb.GetGameState().GetLootBoxes() {
		if !initialized {
			nearestBox = id
			initialized = true
		}
		lootPos := loot.GetPosition()
		currDist = physics.ComputeDistance(currLocation, lootPos)
		if (currDist < shortestDist) && (loot.GetColour() == bb.GetColour()) {
			nearestBox = id
			shortestDist = currDist
		}
	}

	return nearestBox, shortestDist
}

func (bb *Biker1) ProposeDirection() uuid.UUID {
	// all logic for nominations goes in here
	// find nearest coloured box
	// if we can reach it, nominate it
	// if a box exists but we can't reach it, we nominate the box closest to that that we can reach
	// else, nominate nearest box TODO

	// necessary functions:
	// find nearest coloured box: DONE
	// for a box, see if we can reach it -> distance to box from us, our energy level -> function verifies if our energy means we can travel far enough to reach box
	// to do the above, need a function that converts energy to reachable distance
	// function to return nearest box in our reach to a box (our colour) that is out of reach
	// function that returns all the boxes we can reach

	nearestBox, distanceToNearestBox := bb.nearestLootColour()
	// TODO: check if nearestBox actually exists
	reachableDistance := bb.energyToReachableDistance(bb.GetEnergyLevel()) // TODO add all other biker energies
	if distanceToNearestBox < reachableDistance {
		return nearestBox
	}

	nearestReachableBox := bb.getNearestReachableBox()

	return nearestReachableBox
}
func (bb *Biker1) distanceToBox(box uuid.UUID) float64 {
	currLocation := bb.GetLocation()
	boxPos := bb.GetGameState().GetLootBoxes()[box].GetPosition()
	currDist := physics.ComputeDistance(currLocation, boxPos)
	return currDist
}

func (bb *Biker1) findRemainingEnergyAfterReachingBox(box uuid.UUID) float64 {
	dist := physics.ComputeDistance(bb.GetLocation(), bb.GetGameState().GetLootBoxes()[box].GetPosition())
	remainingEnergy := bb.distanceToEnergy(dist, bb.GetEnergyLevel())
	return remainingEnergy
}

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *Biker1) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	// add in voting logic using knowledge of everyone's nominations:

	// for all boxes, rule out any that you can't reach
	// if no boxes left, go for nearest one
	// else if you can reach a box, if someone else can't reach any boxes, vote the box nearest to them (altruistic - add later?)
	// else for every reachable box:
	// calculate energy left if you went there
	// function: calculate energy left given distance moved
	// scan area around box for other boxes based on energy left after reaching it
	// function: given energy and a coordinate on the map, get all boxes that are reachable from that coordinate
	// if our colour is in those boxes, assign the number of people who voted for that box as the score, else assign, 0
	// set highest score box to 1, rest to 0 (subject to change)
	votes := make(voting.LootboxVoteMap)
	maxDist := bb.energyToReachableDistance(bb.GetEnergyLevel())

	// pseudocode:
	// loop through proposals
	// for each box, add 1 to value of key=box_id in dic
	proposalVotes := make(map[uuid.UUID]int)

	maxVotes := 1
	curVotes := 1
	for _, proposal := range proposals {
		_, ok := proposalVotes[proposal]
		if !ok {
			proposalVotes[proposal] = 1
		} else {
			proposalVotes[proposal] += 1
			if proposal != proposals[bb.GetID()] {
				curVotes = proposalVotes[proposal]
				if curVotes > maxVotes {
					maxVotes = curVotes
				}
			}
		}
	}
	distToBoxMap := make(map[uuid.UUID]float64)
	for _, proposal := range proposals {
		distToBoxMap[proposal] = bb.distanceToBox(proposal)
		if distToBoxMap[proposal] <= maxDist { //if reachable
			// if box is our colour and number of proposals is majority, make it 1, rest 0, return
			if bb.GetGameState().GetLootBoxes()[proposal].GetColour() == bb.GetColour() {
				if proposalVotes[proposal] >= maxVotes { // to c
					for _, proposal1 := range proposals {
						if proposal1 == proposal {
							votes[proposal1] = 1
						} else {
							votes[proposal1] = 0
						}
					}
					break
				} else {
					votes[proposal] = float64(proposalVotes[proposal])
				}
			}
			// calculate energy left if we went here
			remainingEnergy := bb.findRemainingEnergyAfterReachingBox(proposal)
			// find nearest reachable boxes from current coordinate
			isColourNear := bb.checkBoxNearColour(proposal, remainingEnergy)
			// assign score of number of votes for this box if our colour is nearby
			if isColourNear {
				votes[proposal] = float64(proposalVotes[proposal])
			} else {
				votes[proposal] = 0.0
			}
		} else {
			votes[proposal] = 0.0
		}
	}

	// Check if all votes are 0
	allVotesZero := true
	for _, value := range votes {
		if value != 0 {
			allVotesZero = false
			break
		}
	}

	// If all votes are 0, nominate the nearest box
	// Maybe nominate our box?
	if allVotesZero {
		minDist := math.MaxFloat64
		var nearestBox uuid.UUID
		for _, proposal := range proposals {
			if distToBoxMap[proposal] < minDist {
				minDist = distToBoxMap[proposal]
				nearestBox = proposal
			}
		}
		votes[nearestBox] = 1
		return votes
	}

	// Normalize the values in votes so that the values sum to 1
	sum := 0.0
	for _, value := range votes {
		sum += value
	}
	for key := range votes {
		votes[key] /= sum
	}
	bb.recentVote = votes
	return votes
}

// -----------------END OF DIRECTION DECISION FUNCTIONS------------------