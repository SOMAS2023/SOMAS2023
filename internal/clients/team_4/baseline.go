package team_4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"math"
	"sort"

	"github.com/google/uuid"
)

type IBaselineAgent interface {
	objects.IBaseBiker

	///////////////////////// INCOMPLETE FUNCTIONS /////////////////////////////
	CalculateReputation() map[uuid.UUID]float64    //calculate reputation matrix
	CalculateHonestyMatrix() map[uuid.UUID]float64 //calculate honesty matrix

	DecideAction() objects.BikerAction //determines what action the agent is going to take this round. (changeBike or Pedal)
	ChangeBike() uuid.UUID             //called when biker wants to change bike, it will choose which bike to try and join

	////////////////// CURRENTLY NOT CONSIDERED FUNCTIONS ///////////////////////
	DecideGovernance() utils.Governance //decide the governance system

	////////////////////////// IMPLEMENTED FUNCTIONS //////////////////////////
	ProposeDirection() uuid.UUID                                                //returns the id of the desired lootbox
	FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap //stage 3 of direction voting
	DecideAllocation() voting.IdVoteMap                                         //decide the allocation parameters
	DecideJoining(pendinAgents []uuid.UUID) map[uuid.UUID]bool                  //decide whether to accept or not accept bikers, ranks the ones
	nearestLoot() uuid.UUID                                                     //returns the id of the nearest lootbox
	DecideForce(direction uuid.UUID)                                            //defines the vector you pass to the bike: [pedal, brake, turning]
	VoteForKickout() map[uuid.UUID]int

	///////////////////////// DICATOR FUNCTIONS ///////////////////////////////////
	VoteDictator() voting.IdVoteMap
	DictateDirection() uuid.UUID                //called only when the agent is the dictator
	DecideKickOut() []uuid.UUID                 //decide which agents to kick out (dictator)
	DecideDictatorAllocation() voting.IdVoteMap //decide the allocation (dictator)

	///////////////////////// LEADER FUNCTIONS ///////////////////////////////////
	VoteLeader() voting.IdVoteMap
	DecideWeights(action utils.Action) map[uuid.UUID]float64 // decide on weights for various actions

	////////////////////////// HELPER FUNCTIONS////////////////////////////////////////
	UpdateDecisionData()           //updates all the data needed for the decision making process(call at the start of any decision making function)
	getHonestyAverage() float64    //returns the average honesty of all agents
	getReputationAverage() float64 //returns the average reputation of all agents

	rankFellowsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) //returns normal rank of fellow bikers reputation
	rankFellowsHonesty(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error)    //returns normal rank of fellow bikers honesty

	rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) //returns ranking of the proposed lootboxes

	IncreaseHonesty(agentID uuid.UUID, increaseAmount float64)
	DecreaseHonesty(agentID uuid.UUID, decreaseAmount float64)

	/////////////////////////// PRINT FUNCTIONS ///////////////////////////////////
	DisplayFellowsEnergyHistory()
	DisplayFellowsHonesty()
	DisplayFellowsReputation()
}

type BaselineAgent struct {
	objects.BaseBiker
	lootBoxColour     utils.Colour
	mylocationHistory []utils.Coordinates     //log location history for this agent
	energyHistory     map[uuid.UUID][]float64 //log energy level for all agents
	reputation        map[uuid.UUID]float64   //record reputation for other agents, 0-1
	honestyMatrix     map[uuid.UUID]float64   //record honesty for other agents, 0-1
}

// ////////////////////////////////////////////////// HELPER FUNCTIONS ////////////////////////////////////////////////////////
func (agent *BaselineAgent) UpdateDecisionData() {
	//Initialize mapping if not initialized yet (= nil)
	if agent.energyHistory == nil {
		agent.energyHistory = make(map[uuid.UUID][]float64)
	}
	if len(agent.mylocationHistory) == 0 {
		agent.mylocationHistory = make([]utils.Coordinates, 0)
	}
	if agent.honestyMatrix == nil {
		agent.honestyMatrix = make(map[uuid.UUID]float64)
	}
	if agent.reputation == nil {
		agent.reputation = make(map[uuid.UUID]float64)
	}
	fmt.Println("Updating decision data ...")
	//update location history for the agent
	agent.mylocationHistory = append(agent.mylocationHistory, agent.GetLocation())
	//get fellow bikers
	fellowBikers := agent.GetFellowBikers()
	//update energy history for each fellow biker
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		currentEnergyLevel := fellow.GetEnergyLevel()
		//Append bikers current energy level to the biker's history
		agent.energyHistory[fellowID] = append(agent.energyHistory[fellowID], currentEnergyLevel)
	}
	//call reputation and honesty matrix to calcuiate/update them
	//save updated reputation and honesty matrix
	agent.CalculateReputation()
	agent.CalculateHonestyMatrix()
	//agent.DisplayFellowsEnergyHistory()
	// agent.DisplayFellowsHonesty()
	// agent.DisplayFellowsReputation()
}

func (agent *BaselineAgent) rankFellowsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	totalsum := float64(0)
	rank := make(map[uuid.UUID]float64)

	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		totalsum += agent.reputation[fellowID]
	}
	//normalize the reputation
	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		rank[fellowID] = float64(agent.reputation[fellowID] / totalsum)
	}
	return rank, nil
}

func (agent *BaselineAgent) rankFellowsHonesty(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	totalsum := float64(0)
	rank := make(map[uuid.UUID]float64)

	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		totalsum += agent.honestyMatrix[fellowID]
	}
	//normalize the honesty
	for _, fellow := range agentsOnBike {
		fellowID := fellow.GetID()
		rank[fellowID] = float64(agent.honestyMatrix[fellowID] / totalsum)
	}
	return rank, nil
}

func (agent *BaselineAgent) getReputationAverage() float64 {
	sum := float64(0)
	//loop through all bikers find the average reputation
	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()
			sum += agent.reputation[bikerID]
		}
	}
	return sum / float64(len(agent.reputation))
}

func (agent *BaselineAgent) getHonestyAverage() float64 {
	sum := float64(0)
	//loop through all bikers find the average honesty
	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()
			sum += agent.honestyMatrix[bikerID]
		}
	}
	return sum / float64(len(agent.honestyMatrix))
}

func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (voting.LootboxVoteMap, error) {
	rank := make(voting.LootboxVoteMap) //make(map[uuid.UUID]float64)
	ranksum := make(map[uuid.UUID]float64)
	totalsum := float64(0)
	distanceRank := float64(0)
	w1 := float64(5.0)  //weight for distance
	w2 := float64(1.0)  //weight for reputation
	w3 := float64(1.0)  //weight for honesty
	w4 := float64(10.0) //weight for distance from Audi

	//if energy level is below threshold, increase weighting towards distance
	minEnergyThreshold := 0.2
	if agent.GetEnergyLevel() < minEnergyThreshold {
		w1 *= 2
	}
	totaloptions := len(proposedLootBox)
	audiPos := agent.GetGameState().GetAudi().GetPosition()

	fellowBikers := agent.GetFellowBikers()
	//This is the relavtive reputation and honest for bikers my bike
	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	//This is the absolute reputation and honest for bikers my bike
	// reputationRank  := agent.reputation
	// honestyRank  := agent.honestyMatrix
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}
	//sort proposed loot boxes by distance from agent
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(agent.GetLocation(), proposedLootBox[i].GetPosition()) < physics.ComputeDistance(agent.GetLocation(), proposedLootBox[j].GetPosition())
	})

	for i, lootBox := range proposedLootBox {
		lootboxID := lootBox.GetID()
		distanceFromAudi := physics.ComputeDistance(audiPos, lootBox.GetPosition())

		//loop through all fellow bikers and check if they have the same colour as the lootbox
		for _, fellow := range fellowBikers {
			distanceRank := float64(totaloptions - i)
			fellowID := fellow.GetID()
			if fellow.GetColour() == lootBox.GetColour() {
				weight := (w1 * distanceRank) + (w2 * reputationRank[fellowID]) + (w3 * honestyRank[fellowID]) + (w4 * distanceFromAudi)
				ranksum[lootboxID] += weight
				totalsum += weight
			}
		}

		if lootBox.GetColour() == agent.GetColour() {
			weight := (distanceRank * w1 * 1.25) + (w4 * distanceFromAudi)
			ranksum[lootboxID] += weight
			totalsum += weight
		}
		if ranksum[lootboxID] == 0 {
			weight := (distanceRank * w1 * 2.6) + (w4 * distanceFromAudi)
			ranksum[lootboxID] = weight
			totalsum += weight
		}
	}
	for _, lootBox := range proposedLootBox {
		rank[lootBox.GetID()] = ranksum[lootBox.GetID()] / totalsum
	}

	return rank, nil
}

/////////////////////////////////////////////// DECISION FUNCTIONS /////////////////////////////////////////////////////////

func (agent *BaselineAgent) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	fmt.Println("Final Direction Vote")
	agent.UpdateDecisionData()
	//We need to fix this ASAP
	boxesInMap := agent.GetGameState().GetLootBoxes()
	boxProposed := make([]objects.ILootBox, len(proposals))
	count := 0
	for _, i := range proposals {
		boxProposed[count] = boxesInMap[i]
		count++
	}

	rank, e := agent.rankTargetProposals(boxProposed)
	if e != nil {
		panic("unexpected error!")
	}
	return rank
}

func (agent *BaselineAgent) DecideAllocation() voting.IdVoteMap {
	fmt.Println("Decide Allocation")
	agent.UpdateDecisionData()
	distribution := make(voting.IdVoteMap) //make(map[uuid.UUID]float64)
	fellowBikers := agent.GetFellowBikers()
	totalEnergySpent := float64(0)
	totalAllocation := float64(0)

	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}

	for _, fellow := range fellowBikers {
		w1 := 2.0 //weight for reputation
		w2 := 2.0 //weight for honesty
		w3 := 1.0 //weight for energy spent
		w4 := 2.0 //weight for energy level
		fellowID := fellow.GetID()
		energyLog := agent.energyHistory[fellowID]
		energySpent := energyLog[len(energyLog)-2] - energyLog[len(energyLog)-1]
		totalEnergySpent += energySpent
		// In the case where the I am the same colour as the lootbox
		if fellowID == agent.GetID() {
			w4 = 3.0
			if agent.lootBoxColour == agent.GetColour() {
				w3 = 3.0
			}
		}
		distribution[fellow.GetID()] = float64((w1 * reputationRank[fellowID]) + (w2 * honestyRank[fellowID]) + (w3 * energySpent) + (w4 * fellow.GetEnergyLevel()))
		// distribution[fellow.GetID()] = energySpent * rand.Float64() // random for now
		totalAllocation += distribution[fellow.GetID()]
	}

	//normalize the distribution
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		distribution[fellowID] = distribution[fellowID] / totalAllocation
	}

	return distribution
}

func (agent *BaselineAgent) DecideForce(direction uuid.UUID) {

	currLocation := agent.GetLocation()
	targetLoot := direction
	currentLootBoxes := agent.GetGameState().GetLootBoxes()
	audiPos := agent.GetGameState().GetAudi().GetPosition()
	distanceFromAudi := physics.ComputeDistance(currLocation, audiPos)
	distanceThreshold := 20.0
	pedalForce := 1.0

	if distanceFromAudi < distanceThreshold {
		deltaX := audiPos.X - currLocation.X
		deltaY := audiPos.Y - currLocation.Y
		// Steer in opposite direction to audi
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi
		// Steer in opposite direction to audi
		var flipAngle float64
		if normalisedAngle < 0.0 {
			flipAngle = normalisedAngle + 1.0
		} else if normalisedAngle > 0.0 {
			flipAngle = normalisedAngle - 1.0
		}
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: flipAngle - agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetOrientation(),
		}
		escapeAudiForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		agent.SetForces(escapeAudiForces)
	} else {
		targetPos := currentLootBoxes[targetLoot].GetPosition()
		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetOrientation(),
		}
		if agent.GetEnergyLevel() <= 0.5 {
			pedalForce = pedalForce * agent.GetEnergyLevel()
		}
		nearestBoxForces := utils.Forces{
			Pedal:   pedalForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		agent.SetForces(nearestBoxForces)
	}
}

func (agent *BaselineAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	agent.UpdateDecisionData()
	fellowBikers := agent.GetFellowBikers()
	spare := 8 - len(fellowBikers)
	decision := make(map[uuid.UUID]bool)
	w1 := 3.0 //weight for reputation
	w2 := 1.0 //weight for honesty

	// Temporary slice to hold combined reputation and honesty scores along with agent ID
	type agentScore struct {
		ID    uuid.UUID
		Score float64
	}
	var scoredAgents []agentScore

	for _, pendingAgent := range pendingAgents {
		reputation := agent.reputation[pendingAgent]
		honesty := agent.honestyMatrix[pendingAgent]
		scoredAgents = append(scoredAgents, agentScore{ID: pendingAgent, Score: ((w1 * reputation) + (w2 * honesty))})
	}
	// Sort the slice based on the combined score
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].Score > scoredAgents[j].Score
	})

	// Make decisions based on the sorted slice
	for i, scoredAgent := range scoredAgents {
		// Example decision making logic
		if i < spare {
			decision[scoredAgent.ID] = true // Accept if there's spare capacity
		} else {
			decision[scoredAgent.ID] = false // Reject if no capacity
		}
	}
	return decision
}

func (agent *BaselineAgent) DecideGovernance() utils.Governance {
	// Change behaviour here to return different governance
	return utils.Democracy
}

func (agent *BaselineAgent) VoteForKickout() map[uuid.UUID]int {
	agent.UpdateDecisionData()
	fmt.Println("Vote for Kickout")
	voteResults := make(map[uuid.UUID]int)

	fellowBikers := agent.GetFellowBikers()
	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)

	if e1 != nil || e2 != nil {
		panic("unexpected error!")

	}
	combined := make(map[uuid.UUID]float64)
	worstRank := float64(1)

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		if combined[fellowID] == worstRank && fellowID != uuid.Nil {

			if fellowID != agent.GetID() {
				combined[fellowID] = reputationRank[fellowID] * honestyRank[fellowID]
				if combined[fellowID] < worstRank {
					worstRank = combined[fellowID]
				}
			} else {
				combined[fellowID] = 1.0
			}
		}
	}

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		if len(fellowBikers) > 5 {
			if combined[fellowID] == worstRank && fellowID != uuid.Nil {
				if agent.reputation[fellowID] < agent.getReputationAverage() || agent.honestyMatrix[fellowID] < agent.getHonestyAverage() {
					if len(fellowBikers) > 4 {
						voteResults[fellowID] = 1
					} else {
						voteResults[fellowID] = 0
					}
				}
			} else {
				voteResults[fellowID] = 0
			}
		} else {
			voteResults[fellowID] = 0
		}
	}
	voteResults[agent.GetID()] = 0
	println("the voting results are:", voteResults)
	return voteResults
}

func (agent *BaselineAgent) nearestLoot() uuid.UUID {
	currLocation := agent.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestBox uuid.UUID
	var currDist float64
	for _, loot := range agent.GetGameState().GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return nearestBox
}

func (agent *BaselineAgent) ProposeDirection() uuid.UUID {
	fmt.Println("Propose Direction")
	agent.UpdateDecisionData()

	var lootBoxesWithinThreshold []objects.ILootBox
	distanceThresholdFromAudi := 20.0 // adjust this value as needed
	audiPos := agent.GetGameState().GetAudi().GetPosition()
	agentLocation := agent.GetLocation() // agent's location

	for _, lootbox := range agent.GetGameState().GetLootBoxes() {
		if physics.ComputeDistance(lootbox.GetPosition(), audiPos) > distanceThresholdFromAudi {
			lootBoxesWithinThreshold = append(lootBoxesWithinThreshold, lootbox)
		}
	}

	// Sort the lootboxes within threshold by distance from the agent
	sort.Slice(lootBoxesWithinThreshold, func(i, j int) bool {
		return physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[i].GetPosition()) <
			physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[j].GetPosition())
	})

	// Select the closest lootbox if any are within the threshold
	if len(lootBoxesWithinThreshold) > 0 {
		closestLootBox := lootBoxesWithinThreshold[0]
		return closestLootBox.GetID()
	} else {
		return agent.nearestLoot()
	}
}

func (agent *BaselineAgent) ChangeBike() uuid.UUID {
	megaBikes := agent.GetGameState().GetMegaBikes()
	optimalBike := agent.GetBike()
	weight := float64(-99)
	for _, bike := range megaBikes {
		if bike.GetID() != megaBikes[agent.GetBike()].GetID() && bike.GetID() != uuid.Nil { //get all bikes apart from our agent's bike
			bikeWeight := float64(0)

			for _, biker := range bike.GetAgents() {
				if biker.GetColour() == agent.GetColour() {
					bikeWeight += 1.8
				} else {
					bikeWeight += 1
				}
			}

			if bikeWeight > weight {
				optimalBike = bike.GetID()
			}
		}
	}
	return optimalBike
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//changed version

func (agent *BaselineAgent) CalculateReputation() {
	////////////////////////////
	//  As the program I used for debugging invoked "padal" and "break" with values of 0, I conducted tests using random numbers.
	// In case of an updated main program, I will need to adjust the parameters and expressions of the reputation matrix.
	// The current version lacks real data during the debugging process.
	////////////////////////////
	megaBikes := agent.GetGameState().GetMegaBikes()

	for _, bike := range megaBikes {
		// Get all agents on MegaBike
		fellowBikers := bike.GetAgents()

		// Iterate over each agent on MegaBike, generate reputation assessment
		for _, otherAgent := range fellowBikers {
			// Exclude self
			selfTest := otherAgent.GetID() //nolint
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
			}

			// Monitor otherAgent's location
			// location := otherAgent.GetLocation()
			// RAP := otherAgent.GetResourceAllocationParams()
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Location:", location, "ResourceAllocationParams:", RAP)

			// Monitor otherAgent's forces
			historyenergy := agent.energyHistory[otherAgent.GetID()]
			lastEnergy := 1.0
			if len(historyenergy) >= 2 {
				lastEnergy = historyenergy[len(historyenergy)-2]
				// rest of your code
			} else {
				lastEnergy = 0.0
			}
			energyLevel := otherAgent.GetEnergyLevel()
			ReputationEnergy := float64((lastEnergy)) / energyLevel //CAUTION: REMOVE THE RANDOM VALUE
			//print("我是大猴子")
			//fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Forces:", ReputationEnergy)

			// Monitor otherAgent's bike status
			bikeStatus := otherAgent.GetBikeStatus()
			// Convert the boolean value to float64 and print the result
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}
			//fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Bike_Shift", float64(ReputationBikeShift))

			// Calculate Overall_reputation
			OverallReputation := ReputationEnergy * ReputationBikeShift
			//fmt.Println("Agent ID:", otherAgent.GetID(), "Overall Reputation:", OverallReputation)

			// Store Overall_reputation in the reputation map
			agent.reputation[otherAgent.GetID()] = OverallReputation
		}
	}
	/* 	for agentID, agentReputation := range agent.reputation {
		print("Agent ID: ", agentID.String(), ", Reputation: ", agentReputation, "\n")
	} */

}

/* // Reputation and Honesty Matrix Teams Must Implement these or similar functions

func (agent *BaselineAgent) CalculateReputation() {
	////////////////////////////
	//  As the program I used for debugging invoked "padal" and "break" with values of 0, I conducted tests using random numbers.
	// In case of an updated main program, I will need to adjust the parameters and expressions of the reputation matrix.
	// The current version lacks real data during the debugging process.
	////////////////////////////
	megaBikes := agent.GetGameState().GetMegaBikes()

	for _, bike := range megaBikes {
		// Get all agents on MegaBike
		fellowBikers := bike.GetAgents()

		// Iterate over each agent on MegaBike, generate reputation assessment
		for _, otherAgent := range fellowBikers {
			// Exclude self
			selfTest := otherAgent.GetID() //nolint
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
			}

			// Monitor otherAgent's location
			// location := otherAgent.GetLocation()
			// RAP := otherAgent.GetResourceAllocationParams()
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Location:", location, "ResourceAllocationParams:", RAP)

			// Monitor otherAgent's forces
			forces := otherAgent.GetForces()
			energyLevel := otherAgent.GetEnergyLevel()
			ReputationForces := float64(forces.Pedal+forces.Brake+rand.Float64()) / energyLevel //CAUTION: REMOVE THE RANDOM VALUE
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Forces:", ReputationForces)

			// Monitor otherAgent's bike status
			bikeStatus := otherAgent.GetBikeStatus()
			// Convert the boolean value to float64 and print the result
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Bike_Shift", float64(ReputationBikeShift))

			// Calculate Overall_reputation
			OverallReputation := ReputationForces * ReputationBikeShift
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Overall Reputation:", OverallReputation)

			// Store Overall_reputation in the reputation map
			agent.reputation[otherAgent.GetID()] = OverallReputation
		}
	}
	// for agentID, agentReputation := range agent.reputation {
	// 	print("Agent ID: ", agentID.String(), ", Reputation: ", agentReputation, "\n")
	// }
} */

func (agent *BaselineAgent) CalculateHonestyMatrix() {
	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()
			agent.honestyMatrix[bikerID] = 1.0
		}
	}

}

///////////////////////////////////// LEADER FUNCTIONS ///////////////////////////////////////

// defaults to an equal distribution over all agents for all actions
func (agent *BaselineAgent) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	weights := make(map[uuid.UUID]float64)
	fellows := agent.GetFellowBikers()
	for _, fellow := range fellows {
		if fellow.GetID() != uuid.Nil {
			weights[fellow.GetID()] = 1.0
		} else {
			weights[fellow.GetID()] = 0.0
		}
	}
	return weights
}

func (agent *BaselineAgent) VoteLeader() voting.IdVoteMap {
	agent.UpdateDecisionData()
	votes := make(voting.IdVoteMap)
	fellowBikers := agent.GetFellowBikers()
	totalsum := float64(0)
	w1 := 3.0 //weight for reputation
	w2 := 1.0 //weight for honesty

	// Temporary slice to hold combined reputation and honesty scores along with agent ID
	type agentScore struct {
		ID    uuid.UUID
		Score float64
	}
	var scoredAgents []agentScore

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		reputation := agent.reputation[fellowID]
		honesty := agent.honestyMatrix[fellowID]
		scoredAgents = append(scoredAgents, agentScore{ID: fellowID, Score: ((w1 * reputation) + (w2 * honesty))})
	}
	// Sort the slice based on the combined score
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].Score > scoredAgents[j].Score
	})

	for i, scoredAgent := range scoredAgents {
		weight := float64(len(scoredAgents) - i)
		votes[scoredAgent.ID] = weight
		totalsum += weight
	}
	votes[agent.GetID()] = 20.0
	totalsum += 20.0
	//normalize the vote
	for _, scoredAgent := range scoredAgents {
		votes[scoredAgent.ID] = votes[scoredAgent.ID] / totalsum
	}
	return votes
}

/////////////////////////////////// DICATOR FUNCTIONS /////////////////////////////////////

func (agent *BaselineAgent) VoteDictator() voting.IdVoteMap {
	agent.UpdateDecisionData()
	votes := make(voting.IdVoteMap)
	fellowBikers := agent.GetFellowBikers()
	totalsum := float64(0)
	w1 := 3.0 //weight for reputation
	w2 := 1.0 //weight for honesty

	// Temporary slice to hold combined reputation and honesty scores along with agent ID
	type agentScore struct {
		ID    uuid.UUID
		Score float64
	}
	var scoredAgents []agentScore

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		reputation := agent.reputation[fellowID]
		honesty := agent.honestyMatrix[fellowID]
		scoredAgents = append(scoredAgents, agentScore{ID: fellowID, Score: ((w1 * reputation) + (w2 * honesty))})
	}
	// Sort the slice based on the combined score
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].Score > scoredAgents[j].Score
	})

	// Make decisions based on the sorted slice
	for i, scoredAgent := range scoredAgents {
		weight := float64(len(scoredAgents) - i)
		votes[scoredAgent.ID] = weight
		totalsum += weight
	}
	votes[agent.GetID()] = 20.0
	totalsum += 20.0
	//normalize the vote
	for _, scoredAgent := range scoredAgents {
		votes[scoredAgent.ID] = votes[scoredAgent.ID] / totalsum
	}
	return votes
}

func (agent *BaselineAgent) DictateDirection() uuid.UUID {
	fmt.Println("Dictate Direction")
	agent.UpdateDecisionData()

	var lootBoxesWithinThreshold []objects.ILootBox
	distanceThresholdFromAudi := 20.0 // adjust this value as needed
	audiPos := agent.GetGameState().GetAudi().GetPosition()
	agentLocation := agent.GetLocation() // agent's location

	for _, lootbox := range agent.GetGameState().GetLootBoxes() {
		if physics.ComputeDistance(lootbox.GetPosition(), audiPos) > distanceThresholdFromAudi {
			lootBoxesWithinThreshold = append(lootBoxesWithinThreshold, lootbox)
		}
	}

	// Sort the lootboxes within threshold by distance from the agent
	sort.Slice(lootBoxesWithinThreshold, func(i, j int) bool {
		return physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[i].GetPosition()) <
			physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[j].GetPosition())
	})

	// Select the closest lootbox if any are within the threshold
	if len(lootBoxesWithinThreshold) > 0 {
		closestLootBox := lootBoxesWithinThreshold[0]
		return closestLootBox.GetID()
	} else {
		return agent.nearestLoot()
	}
}

func (agent *BaselineAgent) DecideKickOut() []uuid.UUID {
	fmt.Println("Decide Kickout")
	kickoutResults := make([]uuid.UUID, 0)
	agent.UpdateDecisionData()

	fellowBikers := agent.GetFellowBikers()
	if len(fellowBikers) > 2 {

		reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
		honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
		if e1 != nil || e2 != nil {
			panic("unexpected error!")
		}
		combined := make(map[uuid.UUID]float64)
		worstRank := float64(1)

		for _, fellow := range fellowBikers {
			fellowID := fellow.GetID()
			if combined[fellowID] == worstRank && fellowID != uuid.Nil {

				if fellowID != agent.GetID() {
					combined[fellowID] = reputationRank[fellowID] * honestyRank[fellowID]
					if combined[fellowID] < worstRank {
						worstRank = combined[fellowID]
					}
				} else {
					combined[fellowID] = 1.0
				}
			}
		}
		for _, fellow := range fellowBikers {
			fellowID := fellow.GetID()
			if fellowID != agent.GetID() {
				if combined[fellowID] == worstRank && fellowID != uuid.Nil {
					if agent.reputation[fellowID] < agent.getReputationAverage() || agent.honestyMatrix[fellowID] < agent.getHonestyAverage() {
						kickoutResults = append(kickoutResults, fellowID)
					}
				}
			}
		}
	}
	return kickoutResults

}

func (agent *BaselineAgent) DecideDictatorAllocation() voting.IdVoteMap {
	fmt.Println("Dictate Allocation")
	agent.UpdateDecisionData()
	distribution := make(voting.IdVoteMap)
	fellowBikers := agent.GetFellowBikers()
	totalEnergySpent := float64(0)
	totalAllocation := float64(0)

	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}

	for _, fellow := range fellowBikers {
		w1 := 2.0 //weight for reputation
		w2 := 2.0 //weight for honesty
		w3 := 1.0 //weight for energy spent
		w4 := 2.0 //weight for energy level
		fellowID := fellow.GetID()
		energyLog := agent.energyHistory[fellowID]
		energySpent := energyLog[len(energyLog)-2] - energyLog[len(energyLog)-1]
		totalEnergySpent += energySpent
		// In the case where the I am the same colour as the lootbox
		if fellowID == agent.GetID() && fellowID != uuid.Nil {
			w4 = 3.0
			if agent.lootBoxColour == agent.GetColour() {
				w3 = 3.0
			}
		}
		distribution[fellow.GetID()] = float64((w1 * reputationRank[fellowID]) + (w2 * honestyRank[fellowID]) + (w3 * energySpent) + (w4 * fellow.GetEnergyLevel()))
		// distribution[fellow.GetID()] = energySpent * rand.Float64() // random for now
		totalAllocation += distribution[fellow.GetID()]
	}

	//normalize the distribution
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		distribution[fellowID] = distribution[fellowID] / totalAllocation
	}

	return distribution
}

// //////////////////////////// DISPLAY FUNCTIONS ////////////////////////////////////////
func (agent *BaselineAgent) DisplayFellowsEnergyHistory() {
	fellowBikers := agent.GetFellowBikers()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Energy history for: ", fellowID)
		fmt.Print(agent.energyHistory[fellowID])
		fmt.Println("")
	}
}
func (agent *BaselineAgent) DisplayFellowsHonesty() {
	fellowBikers := agent.GetFellowBikers()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Honesty Matrix for: ", fellowID)
		fmt.Print(agent.honestyMatrix[fellowID])
		fmt.Println("")
	}
}
func (agent *BaselineAgent) DisplayFellowsReputation() {
	fellowBikers := agent.GetFellowBikers()
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		fmt.Println("")
		fmt.Println("Reputation Matrix for: ", fellowID)
		fmt.Print(agent.reputation[fellowID])
		fmt.Println("")
	}
}
