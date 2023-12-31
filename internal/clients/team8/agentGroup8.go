package team8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	"SOMAS2023/internal/common/voting"
	"math"

	"github.com/google/uuid"
)

type IBaselineAgent interface {
	objects.IBaseBiker
}

type Agent8 struct {
	*objects.BaseBiker
	overallLootboxPreferences voting.LootboxVoteMap         //rank score for the lootbox
	agentsActionsMap          map[int]map[uuid.UUID]float64 //action score for each agent for the previous 10 loops (-1, 1)
	loopScoreMap              map[int]map[uuid.UUID]float64 //loop score for each loop for our megabike (-1, 1)
	previousLocation          utils.Coordinates             // record the location of last loop
	previousTargetLocation    utils.Coordinates             //record the final target lootbox of last loop
	previousEnergy            float64                       // record the energy level of last loop
	previousPoints            int                           // record the point of last loop
	messageReputation         map[uuid.UUID]float64         // record the extra reputation from Message System
	//-----------------New parameters for recording social connection-----------------
	previousOwnColour    utils.Colour // record the colour of our bike of last loop
	previousTargetColour utils.Colour // record the colour of target lootbox of last loop
	utility              float64      // New parameters for recording utility
	fairness             float64      // New parameters for recording gini index
	satisfaction         float64      // New parameters for recording satisfaction(0,1)
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DecideGovernance <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// Decide which Governance to use
// New strategy: if the number of agents on our bike is more than 6 or the number of bad agents is more than half of the agents on our bike, we use Leadership
func (bb *Agent8) DecideGovernance() utils.Governance {
	fellowBikers := bb.GetFellowBikers()
	numberOfFellowBikers := len(fellowBikers)
	numberOfBadAgents := 0
	for _, fellowBiker := range fellowBikers {
		// logic of voting weight decision
		if bb.QueryReputation(fellowBiker.GetID()) < 0.0 {
			numberOfBadAgents++
		}
	}

	if numberOfFellowBikers > 6 || numberOfBadAgents > numberOfFellowBikers/2 {
		return utils.Leadership
	} else {
		return utils.Democracy
	}
}

// Decide the voting weight for each agent on the bike for dictator
func (bb *Agent8) VoteDictator() voting.IdVoteMap {
	// initialise the voteMap
	votes := make(voting.IdVoteMap)

	// get all the agent on our bike and iterate the agents
	fellowBikers := bb.GetFellowBikers()
	for _, fellowBiker := range fellowBikers {
		// logic of voting weight decision
		if fellowBiker.GetID() == bb.GetID() {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}

	return votes
}

// Decide the voting weight for each agent on the bike for leader
func (bb *Agent8) VoteLeader() voting.IdVoteMap {
	// initialise the voteMap
	votes := make(voting.IdVoteMap)

	// get all the agent on our bike and iterate the agents
	fellowBikers := bb.GetFellowBikers()
	for _, fellowBiker := range fellowBikers {
		// logic of voting weight decision
		if fellowBiker.GetID() == bb.GetID() {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}

	return votes
}

// ===============================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 1 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// Decide which agent to kickout
func (bb *Agent8) VoteForKickout() map[uuid.UUID]int {
	// initialise the kickoutVotingMap
	voteResults := make(map[uuid.UUID]int)

	// get all the agent on our bike and iterate the agents
	fellowBikers := bb.GetFellowBikers()
	for _, agent := range fellowBikers {
		// logic of kickout decision
		agentID := agent.GetID()
		// if our reputation of the agent is lower than baseline, kick the agent
		if bb.QueryReputation(agentID) < 0.0 {
			voteResults[agentID] = 1
		} else {
			voteResults[agentID] = 0
		}
	}

	return voteResults
}

// When we are dictator, we need to kick out bad agents
func (bb *Agent8) DecideKickOut() []uuid.UUID {
	// initialise the kickoutMap
	kickoutList := make([]uuid.UUID, 0.0)

	// get all the agent on our bike and iterate the agents
	fellowBikers := bb.GetFellowBikers()
	for _, agent := range fellowBikers {
		// logic of kickout decision
		agentID := agent.GetID()
		// if our reputation of the agent is much lower than baseline, kick the agent
		if bb.QueryReputation(agentID) < -0.2 {
			kickoutList = append(kickoutList, agentID)
		}
	}

	return kickoutList
}

// Decision for accept/reject the agent who want to join our bike
// New strategy: We choose the best n score agents to join our bike, where n is the number of vacancies
func (bb *Agent8) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	// initialise the Map
	threshold := GlobalParameters.ThresholdForJoiningDecision
	decision := make(map[uuid.UUID]bool)
	agentMap := bb.UuidToAgentMap(pendingAgents)
	// ------ calulate the vacancies and construct the scoremap ------
	numberOfFellowBikers := len(bb.GetFellowBikers())
	vacanicies := utils.BikersOnBike - numberOfFellowBikers
	scoreMap := make(map[uuid.UUID]float64)
	//----------------
	// iterate the agents who want to join
	for uuid, agent := range agentMap {
		// calculate the Score for each agent for dicision
		var score float64
		if agent.GetColour() == bb.GetColour() {
			score = (agent.GetEnergyLevel()-bb.CalculateAverageEnergy(bb.GetBike()))/bb.CalculateAverageEnergy(bb.GetBike()) +
				bb.QueryReputation(agent.GetID())
		} else {
			score = 0.5*(agent.GetEnergyLevel()-bb.CalculateAverageEnergy(bb.GetBike()))/bb.CalculateAverageEnergy(bb.GetBike()) +
				bb.QueryReputation(agent.GetID())
		}
		//store the score in the map
		scoreMap[uuid] = score

		if score >= threshold {
			decision[uuid] = true
		} else {
			decision[uuid] = false
		}

	}
	if len(scoreMap) == 0 {
		return decision
	}
	// sort the scoreMap
	rankedScoreMap := rankByPreference(scoreMap)
	// Rank the agents based on the score, and choose the best n agents to join our bike, where n is the vanancies of the bike
	for i := 0; i < vacanicies; i++ {
		if scoreMap[rankedScoreMap[i]] >= threshold {
			decision[rankedScoreMap[i]] = true
		} else {
			decision[rankedScoreMap[i]] = false
		}

	}

	return decision
}

// If we want to jump to another bike, we need call this function to find the best bike to join
func (bb *Agent8) ChangeBike() uuid.UUID {
	// Get all the bikes from the game state
	megaBikes := bb.GetGameState().GetMegaBikes()

	// Initialize a map to store Borda scores for each bike
	bordaScores := make(map[uuid.UUID]float64)
	acceptBool := make(map[uuid.UUID]bool)

	// At least our bike will accept us
	acceptBool[bb.GetBike()] = true

	// Iterate through each bike
	for bikeID, megaBike := range megaBikes {
		// Calculate the Borda score for the current bike
		bordaScore := bb.CalculateAverageEnergy(bikeID) + float64(bb.CountAgentsWithSameColour(bikeID)) + bb.countReputationScore(megaBike)

		// Store the Borda score in the map
		bordaScores[bikeID] = bordaScore

		// find the agents on each bike
		agentsOnBike := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()

		// iterate the agents and theck their reputation score to see if we could be accepted.
		var reputationSum = 0.0
		for _, agent := range agentsOnBike {
			reputationSum += agent.GetReputation()[bb.GetID()]
		}
		if reputationSum >= 0.0 {
			acceptBool[bikeID] = true
		}
	}

	// Find the bike with the highest Borda score
	var highestBordaScore float64
	var winningBikeID uuid.UUID
	for bikeID, score := range bordaScores {
		if score > highestBordaScore && acceptBool[bikeID] {
			highestBordaScore = score
			winningBikeID = bikeID
		}
	}

	return winningBikeID
}

//===============================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 2 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// Decide if we need to jump to other bike
func (bb *Agent8) DecideAction() objects.BikerAction {
	// initialise the parameters

	var selfBikeId = bb.GetBike()
	var selfBikeScore = 0.0
	var loopNum = 0.0

	// calculate total reflection score for current bike
	for i := 1; i <= 10; i++ {
		for bikeid, score := range bb.loopScoreMap[i] {
			if bikeid == selfBikeId {
				selfBikeScore += score
				loopNum++
			}
		}
	}
	selfBikeScore = selfBikeScore / loopNum

	// check if we need to change bike
	if selfBikeScore < GlobalParameters.ThresholdForChangingMegabike && bb.GetEnergyLevel() >= 0.7 {
		if bb.ChangeBike() != bb.GetBike() {
			return objects.ChangeBike
		}
	}

	// stay
	return objects.Pedal
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 3 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// Dicide which lootbox to vote for round 1
// New strategy: we propose two consider the angle between our bike and the lootbox and also the awdi
func (bb *Agent8) ProposeDirection() uuid.UUID {
	// Get all lootboxes and initialise the preferences map
	lootBoxes := bb.GetGameState().GetLootBoxes()
	preferences := make(map[uuid.UUID]float64)
	safeAngle := 20.0
	safe := false
	// Iterate the lootboxes and calculate preferences for each
	for boxId, lootBox := range lootBoxes {
		distance := 0.0
		angleBetweenBikeAndBox := 0.0
		angleBetweenBikeAndAwdi := 0.0
		//calcuate the angle between our bike and the lootbox
		// check if we are onbike
		if bb.GetBike() != uuid.Nil {
			// get distance between our bike and the lootbox
			angleBetweenBikeAndBox = math.Atan2(lootBox.GetPosition().Y-bb.GetLocation().Y, lootBox.GetPosition().X-bb.GetLocation().X)/math.Pi -
				bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
				//calculate the angle between our bike and the awdi
			angleBetweenBikeAndAwdi = math.Atan2(bb.GetGameState().GetAwdi().GetPosition().Y-bb.GetLocation().Y, bb.GetGameState().GetAwdi().GetPosition().X-bb.GetLocation().X)/math.Pi -
				bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
			distance = calculateDistance(bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetPosition(), lootBox.GetPosition())
		}
		// check if the color of box is our target and get energylevel
		colorPreference := calculateColorPreference(bb.GetColour(), lootBox.GetColour())
		energyWeighting := bb.GetEnergyLevel()

		// The higher energy, the higher weight for target color
		distanceBoxAwdi := calculateDistance(bb.GetGameState().GetAwdi().GetPosition(), lootBox.GetPosition())

		// check the angle
		if bb.GetBike() != uuid.Nil {
			if (angleBetweenBikeAndBox > 0 && angleBetweenBikeAndAwdi > 0) || (angleBetweenBikeAndBox < 0 && angleBetweenBikeAndAwdi < 0) {
				if math.Abs(angleBetweenBikeAndBox-angleBetweenBikeAndAwdi) > safeAngle {
					safe = true
				}
			} else {
				if (angleBetweenBikeAndBox + angleBetweenBikeAndAwdi) > 0 {
					safe = true
				}
			}

		}
		// if the lootbox is near awdi, igore this box
		if distanceBoxAwdi > 20 && safe {
			// if distanceBoxAwdi > 20 {
			// check our energylevel and calculate the preference of lootbox
			if energyWeighting > GlobalParameters.EnergyThreshold {
				// colorPreference + distancePreference
				preferences[boxId] = colorPreference*energyWeighting +
					(1-energyWeighting)*(GlobalParameters.DistanceThresholdForVoting-distance)/GlobalParameters.DistanceThresholdForVoting
			} else {
				// when the energyLevel is low, just consider the energy and try to survive
				preferences[boxId] = (GlobalParameters.DistanceThresholdForVoting - distance)
			}
		} else {
			preferences[boxId] = 0.0
		}

	}

	// Apply softmax to convert preferences to a probability distribution
	softmaxPreferences := softmax(preferences)

	// Rank loot boxes based on preferences
	rankedLootBoxes := rankByPreference(softmaxPreferences)

	// store the preferencesMap
	bb.overallLootboxPreferences = softmaxPreferences

	return rankedLootBoxes[0]
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 4 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// Dicide the final target when we are the dictator
func (bb *Agent8) DictateDirection() uuid.UUID {
	// find the lootbox with highest preference
	return bb.ProposeDirection()
}

// Dicide the VotingWeight for each agent on our bike
func (bb *Agent8) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	// initialise the weightsMap
	weights := make(map[uuid.UUID]float64)

	// iterate all agents on our bike
	agents := bb.GetFellowBikers()
	for _, agent := range agents {
		// give good agent a higher weighting

		weights[agent.GetID()] = bb.QueryReputation(agent.GetID())

		// give ourself a high weighting
		if agent.GetID() == bb.GetID() {
			weights[agent.GetID()] = 1
		}
	}

	return softmax(weights)
}

// In Democracy, dicide the VotingMap for lootboxList
func (bb *Agent8) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	// initialise the PreferenceScoreMap
	preferenceScores := make(map[uuid.UUID]float64)

	// rerank the all lootboxes in MVP
	_ = bb.ProposeDirection()

	// iterate the input lootboxList and give them preference based on our rankingMap
	for _, lootboxid := range proposals {
		preferenceScores[lootboxid] = bb.overallLootboxPreferences[lootboxid]
	}

	// apply the softmax function
	softmaxScores := softmax(preferenceScores)

	return softmaxScores
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 5 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// Decide the force for current loop with a target lootbox
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	// initialise the parameters
	var forces utils.Forces

	// decide the brake and pedal depends on our energylevel and the speed of bike
	forces.Brake = 0.0
	if bb.GetBike() != uuid.Nil {
		if bb.GetEnergyLevel() > GlobalParameters.EnergyThreshold {
			forces.Pedal = math.Max(0, bb.GetEnergyLevel()*(1-bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetVelocity())/2)
		} else {
			forces.Pedal = math.Max(0, 0.1*(1-bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetVelocity()))
		}
	} else {
		forces.Pedal = 1.0
	}
	lootboxs := bb.GetGameState().GetLootBoxes()

	// --- decide the steering force ---
	// Get the target lootbox object
	var target objects.ILootBox
	for key, value := range lootboxs {
		if key == direction {
			target = value
			break
		}
	}

	// initialise the distance between our bike and the awdi to check our risky score
	distanceAwdiBike := 0.0

	// if we are onbike, calculate the distance between our bike and the awdi
	if bb.GetBike() != uuid.Nil {
		distanceAwdiBike = calculateDistance(bb.GetLocation(), bb.GetGameState().GetAwdi().GetPosition())
	}

	// intialise the angel for tuning
	var angle float64

	// check if we are in danger
	if distanceAwdiBike > 15 {
		angle = math.Atan2(target.GetPosition().Y-bb.GetLocation().Y, target.GetPosition().X-bb.GetLocation().X)/math.Pi -
			bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
	} else {
		// if we are in danger, run away
		if bb.GetBike() != uuid.Nil {
			angle = math.Atan2(bb.GetLocation().Y-bb.GetGameState().GetAwdi().GetPosition().Y, bb.GetLocation().X-bb.GetGameState().GetAwdi().GetPosition().X)/math.Pi -
				bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
		}
	}

	// change the angle to range -1 and 1
	if angle > 1.0 {
		angle -= 2.0
	} else if angle < -1.0 {
		angle += 2.0
	}

	// set turning angle
	forces.Turning.SteerBike = true
	forces.Turning.SteeringForce = angle

	// set forces
	bb.SetForces(forces)

	/*
		The rest code in this function is to record the self-reflection parameter in our bike.
		Since this function will be called in every loop, we put these code here.
	*/

	// update the self-reflection parameter of last loop
	bb.updateAgentActionMap()
	bb.updateLoopScoreMap()
	bb.UpdateReputation()
	// bb.UpdateUtility()

	// store the target and location of current loop for self-reflection parameter calculation
	bb.previousTargetLocation = bb.GetGameState().GetLootBoxes()[direction].GetPosition()
	if bb.GetBike() != uuid.Nil {
		bb.previousLocation = bb.GetLocation()
	}
	bb.previousTargetColour = bb.GetGameState().GetLootBoxes()[direction].GetColour()
	bb.previousOwnColour = bb.GetColour()
	bb.previousEnergy = bb.GetEnergyLevel()
	bb.previousPoints = bb.GetPoints()
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 6 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// Decide the resource allocation if our bike get a lootbox in current loop
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	// get all agents on our bike and initialise the allocationMap
	fellowBikers := bb.GetFellowBikers()
	allocationMap := make(voting.IdVoteMap)

	// iterate the agents on our bike and update the allocationMap
	for _, agent := range fellowBikers {
		if agent.GetID() == bb.GetID() {
			allocationMap[agent.GetID()] = math.Exp(5 - 5*bb.GetEnergyLevel())
		} else {
			allocationMap[agent.GetID()] = bb.QueryReputation(agent.GetID())*0.5 + 0.5
		}
	}

	return softmax(allocationMap)
}

// Decide the resource allocation if we are the dicator
func (bb *Agent8) DecideDictatorAllocation() voting.IdVoteMap {
	return bb.DecideAllocation()
}

// update the score to each agent on our bike for the previous 10 loops for self-reflection
func (bb *Agent8) updateAgentActionMap() {
	currentLoopAgentActionMap := make(map[uuid.UUID]float64)
	agents := bb.GetFellowBikers()
	for _, agent := range agents {
		// agentForce := agent.GetForces()
		// if agentForce.Turning.SteerBike {
		// 	if agentForce.Turning.SteeringForce == bb.GetForces().Turning.SteeringForce {
		// 		currentLoopAgentActionMap[agent.GetID()] = 1 * math.Max(0, agentForce.Pedal-agentForce.Brake)
		// 	} else {
		// 		currentLoopAgentActionMap[agent.GetID()] = -1 * math.Min(1, agentForce.Pedal+agentForce.Brake)
		// 	}
		// } else {
		// 	currentLoopAgentActionMap[agent.GetID()] = 0.7 * math.Max(0, agentForce.Pedal-agentForce.Brake)
		// }
		currentLoopAgentActionMap[agent.GetID()] = agent.GetEnergyLevel()
	}
	if bb.agentsActionsMap == nil {
		bb.agentsActionsMap = make(map[int]map[uuid.UUID]float64)
	}
	for i := 1; i < 10; i++ {
		bb.agentsActionsMap[i] = bb.agentsActionsMap[i+1]
	}
	bb.agentsActionsMap[10] = currentLoopAgentActionMap
}

// update the score for each loop for the previous 10 loops for self-reflection
func (bb *Agent8) updateLoopScoreMap() {
	previousDistanceBikeBox := calculateDistance(bb.previousLocation, bb.previousTargetLocation)
	currentDistanceBikeBox := 0.0
	if bb.GetBike() != uuid.Nil {
		currentDistanceBikeBox = calculateDistance(bb.GetLocation(), bb.previousTargetLocation)
	}
	loopScore := 0.0
	if bb.GetEnergyLevel() < bb.previousEnergy {
		loopScore = (previousDistanceBikeBox - currentDistanceBikeBox) / previousDistanceBikeBox / math.Max(0.01, (bb.previousEnergy-bb.GetEnergyLevel()))
	} else {
		if bb.GetPoints() > bb.previousPoints {
			loopScore = 1 * 5 * (bb.GetEnergyLevel() - bb.previousEnergy)
		} else {
			loopScore = 1 * 1 * (bb.GetEnergyLevel() - bb.previousEnergy)
		}
	}

	if bb.loopScoreMap == nil {
		bb.loopScoreMap = make(map[int]map[uuid.UUID]float64)
	}
	for i := 1; i < 10; i++ {
		bb.loopScoreMap[i] = bb.loopScoreMap[i+1]
	}
	bb.loopScoreMap[10] = make(map[uuid.UUID]float64)
	bb.loopScoreMap[10][bb.GetBike()] = loopScore
}

// update the reputation for other agents
func (bb *Agent8) UpdateReputation() {
	agentCount := make(map[uuid.UUID]float64)
	agentScore := make(map[uuid.UUID]float64)
	for i := 1; i <= 10; i++ {
		for agentId, Score := range bb.agentsActionsMap[i] {
			if Score != 0.0 {
				agentScore[agentId] += Score * math.Pow(0.8, float64(i))
				agentCount[agentId]++
			}
		}
	}
	for agentId, scoreSum := range agentScore {
		bb.SetReputation(agentId, math.Min(1, scoreSum/agentCount[agentId]+bb.messageReputation[agentId]))
		if agentId == bb.GetID() {
			bb.SetReputation(agentId, 1)
		}
	}
}

func (bb *Agent8) UpdateSatisfaction() {
	alpha := 0.8
	beta := 0.5
	distance_prev := calculateDistance(bb.previousLocation, bb.previousTargetLocation)
	distance_curr := calculateDistance(bb.GetLocation(), bb.previousTargetLocation)
	distance_diff := distance_prev - distance_curr
	energy_diff := bb.previousEnergy - bb.GetEnergyLevel()
	if energy_diff < 0 && distance_diff > 0 {
		bb.satisfaction = bb.satisfaction + alpha*(1-bb.satisfaction)
		if bb.previousOwnColour == bb.previousTargetColour {
			bb.satisfaction = 1.1 * bb.satisfaction
		}
	} else {
		bb.satisfaction = bb.satisfaction - beta*bb.satisfaction

	}

}

// =========================================================================================================================================================

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIBaseBiker(baseBiker *objects.BaseBiker) objects.IBaseBiker {
	pointer := &Agent8{
		BaseBiker: baseBiker,
	}
	pointer.GroupID = 8
	pointer.satisfaction = 0.5 // initialise the satisfaction
	return pointer
}
