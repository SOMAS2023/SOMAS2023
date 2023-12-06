// DIRECTION DECISION FUNCTIONS

package team1

import (
	obj "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// ---------------DIRECTION DECISION FUNCTIONS------------------

// Simulates a step of the game, assuming all bikers pedal with the same force as us.
// Returns the distance travelled and the remaining energy
func (bb *Biker1) simulateGameStep(energy float64, velocity float64, force float64, numberOfBikers float64) (float64, float64) {
	totalBikerForce := force * numberOfBikers
	totalMass := utils.MassBike + float64(numberOfBikers)*utils.MassBiker
	acceleration := physics.CalcAcceleration(totalBikerForce, totalMass, velocity)
	distance := velocity + 0.5*acceleration
	energy = energy - force*utils.MovingDepletion
	return distance, energy
}

// Calculates the approximate distance that can be travelled with the given energy
func (bb *Biker1) energyToReachableDistance(energy float64, bike obj.IMegaBike) (float64, float64) {
	distance := 0.0
	totalDistance := 0.0
	remainingEnergy := energy
	var numberOfAgents float64
	if bike.GetID() == bb.GetBikeInstance().GetID() {
		numberOfAgents = float64(len(bike.GetAgents()))
	} else {
		numberOfAgents = float64(len(bike.GetAgents())) + 1
	}

	for remainingEnergy > 0 {
		distance, remainingEnergy = bb.simulateGameStep(remainingEnergy, bb.GetBikeInstance().GetVelocity(), bb.getPedalForce(), numberOfAgents)
		totalDistance = totalDistance + distance
	}
	return remainingEnergy, totalDistance
}

// Calculates the energy remaining after travelling the given distance
func (bb *Biker1) distanceToEnergy(distance float64, initialEnergy float64) float64 {
	totalDistance := 0.0
	remainingEnergy := initialEnergy
	extraDist := 0.0
	for totalDistance < distance {
		extraDist, remainingEnergy = bb.simulateGameStep(remainingEnergy, bb.GetBikeInstance().GetPhysicalState().Mass, utils.BikerMaxForce*remainingEnergy, float64(len(bb.GetFellowBikers())))
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
		_, distance := bb.energyToReachableDistance(ourEnergy, bb.GetBikeInstance())
		if currDist < distance {
			reachableBoxes = append(reachableBoxes, loot.GetID())
		}
	}
	return reachableBoxes
}

// // returns the nearest lootbox with respect to the agent's bike current position
// // in the MVP this is used to determine the pedalling forces as all agent will be
// // aiming to get to the closest lootbox by default
// func (bb *Biker1) nearestLoot() uuid.UUID {
// 	currLocation := bb.GetLocation()
// 	shortestDist := math.MaxFloat64
// 	var nearestBox uuid.UUID
// 	var currDist float64
// 	for _, loot := range bb.GetGameState().GetLootBoxes() {
// 		x, y := loot.GetPosition().X, loot.GetPosition().Y
// 		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
// 		if currDist < shortestDist {
// 			nearestBox = loot.GetID()
// 			shortestDist = currDist
// 		}
// 	}
// 	return nearestBox
// }

func (bb *Biker1) getNearestBox() uuid.UUID {
	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	//default to nearest lootbox
	nearestBox := uuid.Nil
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

func (bb *Biker1) nearestLootColour() (uuid.UUID, float64) {
	shortestDist := math.MaxFloat64
	nearestBox := uuid.Nil
	initialized := false
	currLocation := bb.GetLocation()
	//default to nearest lootbox
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

func (bb *Biker1) FindReachableBoxNearestToBox(nearestColourBox uuid.UUID) uuid.UUID {
	minDist := math.MaxFloat64
	nearestBox := uuid.Nil
	boxes := bb.GetGameState().GetLootBoxes()
	currBoxLocation := boxes[nearestColourBox].GetPosition()
	ourLocation := bb.GetLocation()
	var currDist float64
	var ourDist float64
	for _, loot := range boxes {
		lootPos := loot.GetPosition()
		currDist = physics.ComputeDistance(currBoxLocation, lootPos)
		ourDist = physics.ComputeDistance(lootPos, ourLocation)
		_, reachableDistance := bb.energyToReachableDistance(bb.GetEnergyLevel(), bb.GetBikeInstance())
		if ourDist < reachableDistance && currDist < minDist {
			minDist = currDist
			nearestBox = loot.GetID()
		}
	}
	return nearestBox
}

func (bb *Biker1) ProposeDirection() uuid.UUID {
	// get box of our colour
	nearestColourBox, distanceToNearestBox := bb.nearestLootColour()

	// if box of our colour exists
	if nearestColourBox != uuid.Nil {
		_, reachableDistance := bb.energyToReachableDistance(bb.GetEnergyLevel(), bb.GetBikeInstance())
		if distanceToNearestBox < reachableDistance {
			// if reachable, nominate C
			fmt.Printf("agent %v nominated nearest COLOUR %v box %v \n", bb.GetColour(), bb.GetGameState().GetLootBoxes()[nearestColourBox].GetColour(), nearestColourBox)
			return nearestColourBox
		} else {
			nearestBox := bb.FindReachableBoxNearestToBox(nearestColourBox)
			if nearestBox != uuid.Nil {
				fmt.Printf("agent %v nominated %v box nearest to COLOUR %v %v \n", bb.GetColour(), bb.GetGameState().GetLootBoxes()[nearestBox].GetColour(), bb.GetColour(), nearestBox)
				return nearestBox
			}
		}
	}

	// assumed that box always exists
	nearestBox := bb.getNearestBox()
	fmt.Printf("agent %v nominated nearest %v box %v \n", bb.GetColour(), bb.GetGameState().GetLootBoxes()[nearestBox].GetColour(), nearestBox)
	return nearestBox
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

// Checks whether a box of the desired colour is within our reachable distance from a given box
func (bb *Biker1) checkBoxNearColour(box uuid.UUID, energy float64) uuid.UUID {
	lootBoxes := bb.GetGameState().GetLootBoxes()
	boxPos := lootBoxes[box].GetPosition()
	var currDist float64
	for _, loot := range lootBoxes {
		lootPos := loot.GetPosition()
		currDist = physics.ComputeDistance(boxPos, lootPos)
		_, distance := bb.energyToReachableDistance(energy, bb.GetBikeInstance())
		if currDist < distance && loot.GetColour() == bb.GetColour() {
			return loot.GetID()
		}
	}
	return uuid.Nil
}

func (bb *Biker1) calculateCubeScoreForAgent(agent obj.IBaseBiker) float64 {
	agentPoints := float64(agent.GetPoints())
	ourPoints := float64(bb.GetPoints())
	agentOpinion := bb.opinions[agent.GetID()].opinion
	agentEnergy := agent.GetEnergyLevel()
	ourEnergy := bb.GetEnergyLevel()

	maxPoints := ourPoints + agentPoints // TODO use maxPoints
	relPoints := 0.0
	if maxPoints == 0 {
		relPoints = 0.5
	} else {
		relPoints = (((agentPoints - ourPoints) / (maxPoints + 0.00001)) + 1) / 2
	}

	relEnergy := ((agentEnergy - ourEnergy) + 1) / 2

	// Check Spec for cube explanation
	cubeScore := -0.3*relEnergy - 0.2*relPoints + 0.5*agentOpinion + 0.5
	return cubeScore
}

func (bb *Biker1) calcEnergyScore(destBoxID uuid.UUID, curBoxID uuid.UUID, curEnergy float64) float64 {
	boxes := bb.GetGameState().GetLootBoxes()
	destBox := boxes[destBoxID]
	curBox := boxes[curBoxID]

	destBoxLoc := destBox.GetPosition()
	curBoxLoc := curBox.GetPosition()

	dist := physics.ComputeDistance(destBoxLoc, curBoxLoc)

	energyAfterTravelling := bb.distanceToEnergy(dist, curEnergy)

	return energyAfterTravelling / curEnergy
}

func (bb *Biker1) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	votes := make(voting.LootboxVoteMap)
	_, maxDist := bb.energyToReachableDistance(bb.GetEnergyLevel(), bb.GetBikeInstance())

	// make map of proposal:numberofnoms and find maximum number of votes
	// make map of proposal:proposers
	proposalNoOfNoms := make(map[uuid.UUID]int)
	maxVotes := 1
	curVotes := 1
	for _, proposal := range proposals {
		// initialise final votes as 0.
		votes[proposal] = 0.0

		// add 1 to number of noms for each proposal
		_, exists := proposalNoOfNoms[proposal]
		if !exists {
			proposalNoOfNoms[proposal] = 1
		} else {
			proposalNoOfNoms[proposal] += 1
			if proposal != proposals[bb.GetID()] {
				curVotes = proposalNoOfNoms[proposal]
				if curVotes > maxVotes {
					maxVotes = curVotes
				}
			}
		}
	}
	fmt.Printf("Max Votes: %v\n", maxVotes)
	// if our proposal has majority noms, vote for it
	if proposalNoOfNoms[proposals[bb.GetID()]] > maxVotes {
		votes[proposals[bb.GetID()]] = 1
		fmt.Printf("%v votes for its own nomination %v\n", bb.GetColour(), votes)
		return votes
	}

	// for every nominated box (D)
	for proposer, proposal := range proposals {
		if maxDist < bb.distanceToBox(proposal) {
			// if it is not reachable, ignore
			continue
		}

		// calculate energy left if travelled to D
		remainingEnergy := bb.findRemainingEnergyAfterReachingBox(proposal)

		// check if our colour is reachable with this remaining energy from D
		nearColourBox := bb.checkBoxNearColour(proposal, remainingEnergy)
		// if no boxes of our colour are reachable from this box, assign vote of 0 and continue
		if nearColourBox == uuid.Nil {
			votes[proposal] = 0.0
			continue
		}

		// add this proposer's cube score to the votes map for this box
		votes[proposal] += bb.calculateCubeScoreForAgent(bb.GetAgentFromId(proposer))

		// calculate energy score to add
		energyScore := bb.calcEnergyScore(nearColourBox, proposal, remainingEnergy)

		// add (1/noms_for_D) * energyScore because we add this noms_for_D times
		votes[proposal] += energyScore / float64(proposalNoOfNoms[proposal])
	}

	// if all nominations have score 0, assign 1 to box we nominated
	allVotesZero := true
	for _, value := range votes {
		if value != 0 {
			allVotesZero = false
			break
		}
	}
	if allVotesZero {
		votes[proposals[bb.GetID()]] = 1.
	}

	// normalise values
	sum := 0.0
	for _, value := range votes {
		sum += value
	}
	for key := range votes {
		votes[key] /= sum
	}
	fmt.Printf("%v normalised votes pre-selection: %v\n", bb.GetColour(), votes)

	maxVote := 0.0
	var finalProposal uuid.UUID
	for proposal, value := range votes {
		if value >= maxVote {
			maxVote = value
			finalProposal = proposal
		}
		votes[proposal] = 0.0
	}
	votes[finalProposal] = 1.
	fmt.Printf("%v normalised votes post-selection: %v\n", bb.GetColour(), votes)

	return votes
}

// -----------------END OF DIRECTION DECISION FUNCTIONS------------------
