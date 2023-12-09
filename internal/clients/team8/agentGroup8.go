package team8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"fmt"

	"SOMAS2023/internal/common/voting"
	"math"

	"github.com/google/uuid"
)

type GP struct {
	EnergyThreshold              float64
	DistanceThresholdForVoting   float64
	ThresholdForJoiningDecision  float64
	ThresholdForChangingMegabike float64
}

var GlobalParameters GP = GP{
	EnergyThreshold:              0.6,
	DistanceThresholdForVoting:   (utils.GridHeight + utils.GridWidth) / 4,
	ThresholdForJoiningDecision:  0.2,
	ThresholdForChangingMegabike: 0.5,
}

type IBaselineAgent interface {
	objects.IBaseBiker
}

type Agent8 struct {
	*objects.BaseBiker
	overallLootboxPreferences voting.LootboxVoteMap         //rank score for the lootbox
	agentsActions             map[int]map[uuid.UUID]float64 //action score for each agent for the previous 10 loops (-1, 1)
	loopScore                 map[int]map[uuid.UUID]float64 //loop score for each loop for our megabike (-1, 1)
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DecideGovernance <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// base biker defaults to democracy
func (bb *Agent8) DecideGovernance() utils.Governance {
	// TODO: implement this function
	return utils.Democracy
}

// defaults to voting for first agent in the list
func (bb *Agent8) VoteDictator() voting.IdVoteMap {
	// TODO: implement this function
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()
	for _, fellowBiker := range fellowBikers {
		if fellowBiker.GetID() == bb.GetID() {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	fmt.Println(votes)
	return votes
}

// defaults to voting for first agent in the list
func (bb *Agent8) VoteLeader() voting.IdVoteMap {
	// TODO: implement this function
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()
	for _, fellowBiker := range fellowBikers {
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
func (bb *Agent8) VoteForKickout() map[uuid.UUID]int {
	// TODO: implement this function
	voteResults := make(map[uuid.UUID]int)
	bikeID := bb.GetBike()

	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if bb.QueryReputation(agentID) < 0.0 {
			// random votes to other agents
			voteResults[agentID] = 1 // randomly assigns 0 or 1 vote
		} else {
			voteResults[agentID] = 0
		}
	}

	return voteResults
}

// only called when the agent is the dictator
func (bb *Agent8) DecideKickOut() []uuid.UUID {
	// TODO: implement this function
	kickoutList := make([]uuid.UUID, 0.0)
	fellowBikers := bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if bb.QueryReputation(agentID) < 0.0 {
			// random votes to other agents
			kickoutList = append(kickoutList, agentID)
		}
	}
	return kickoutList
}

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *Agent8) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	threshold := GlobalParameters.ThresholdForJoiningDecision
	decision := make(map[uuid.UUID]bool)
	agentMap := bb.UuidToAgentMap(pendingAgents)

	for uuid, agent := range agentMap {
		var score float64
		if agent.GetColour() == bb.GetColour() {
			score = (agent.GetEnergyLevel() - bb.CalculateAverageEnergy(bb.GetBike())) / bb.CalculateAverageEnergy(bb.GetBike())
		} else {
			score = 0.5 * (agent.GetEnergyLevel() - bb.CalculateAverageEnergy(bb.GetBike())) / bb.CalculateAverageEnergy(bb.GetBike())
		}
		if score >= threshold {
			decision[uuid] = true
		} else {
			decision[uuid] = false
		}

	}

	return decision
}

func (bb *Agent8) ChangeBike() uuid.UUID {
	// Get all the bikes from the game state
	megaBikes := bb.GetGameState().GetMegaBikes()

	// Initialize a map to store Borda scores for each bike
	bordaScores := make(map[uuid.UUID]float64)
	acceptBool := make(map[uuid.UUID]bool)
	acceptBool[bb.GetBike()] = true

	// Iterate through each bike
	for bikeID := range megaBikes {
		// Calculate the Borda score for the current bike
		bordaScore := bb.CalculateAverageEnergy(bikeID) + float64(bb.CountAgentsWithSameColour(bikeID))

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
func (bb *Agent8) DecideAction() objects.BikerAction {
	var selfBikeId = bb.GetBike()
	var selfBikeScore = 0.0
	var loopNum = 0.0

	// calculate total reflection score for current bike
	for _, scoremap := range bb.loopScore {
		for bikeid, score := range scoremap {
			if bikeid == selfBikeId {
				selfBikeScore += score
				loopNum++
			}
		}
	}
	selfBikeScore = selfBikeScore / loopNum

	// check if we need to change bike
	if selfBikeScore < GlobalParameters.ThresholdForChangingMegabike {
		return objects.ChangeBike
	}

	// Default action
	return objects.Pedal
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 3 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) ProposeDirection() uuid.UUID {
	lootBoxes := bb.GetGameState().GetLootBoxes()
	preferences := make(map[uuid.UUID]float64)

	// Calculate preferences
	for _, lootBox := range lootBoxes {
		distance := calculateDistance(bb.GetLocation(), lootBox.GetPosition())
		colorPreference := calculateColorPreference(bb.GetColour(), lootBox.GetColour())
		energyWeighting := bb.GetEnergyLevel()
		// The higher energy, the higher weight for color
		distanceBoxAudi := calculateDistance(bb.GetGameState().GetAudi().GetPosition(), lootBox.GetPosition())
		if distanceBoxAudi > 20 {
			if energyWeighting > GlobalParameters.EnergyThreshold {
				preferences[lootBox.GetID()] = colorPreference*energyWeighting +
					(1-energyWeighting)*(GlobalParameters.DistanceThresholdForVoting-distance)/GlobalParameters.DistanceThresholdForVoting
			} else {
				preferences[lootBox.GetID()] = (GlobalParameters.DistanceThresholdForVoting - distance)
			}
		} else {
			preferences[lootBox.GetID()] = 0.0
		}

	}

	// Apply softmax to convert preferences to a probability distribution
	softmaxPreferences := softmax(preferences)

	// Rank loot boxes based on preferences
	rankedLootBoxes := rankByPreference(softmaxPreferences)

	bb.overallLootboxPreferences = softmaxPreferences

	return rankedLootBoxes[0]
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 4 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) DictateDirection() uuid.UUID {
	// TODO: implement this function
	return bb.ProposeDirection()
}

func (bb *Agent8) LeadDirection() uuid.UUID {
	// TODO: implement this function
	return bb.ProposeDirection()
}

// defaults to an equal distribution over all agents for all actions
func (bb *Agent8) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	// TODO: implement this function
	weights := make(map[uuid.UUID]float64)
	agents := bb.GetFellowBikers()
	for _, agent := range agents {
		weights[agent.GetID()] = bb.QueryReputation(agent.GetID())
	}
	return softmax(weights)
}

// Multi-voting system
func (bb *Agent8) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	// Calculate the biker's individual preference scores
	preferenceScores := make(map[uuid.UUID]float64)
	_ = bb.ProposeDirection()
	for _, lootboxid := range proposals {
		preferenceScores[lootboxid] = bb.overallLootboxPreferences[lootboxid]
	}
	softmaxScores := softmax(preferenceScores)

	return softmaxScores
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 5 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	// TODO: implement this function
	var forces utils.Forces
	forces.Brake = 0.0
	forces.Pedal = 1.0
	lootboxs := bb.GetGameState().GetLootBoxes()
	var target objects.ILootBox
	for key, value := range lootboxs {
		if key == direction {
			target = value
			break
		}
	}
	distanceAudiBike := calculateDistance(bb.GetLocation(), bb.GetGameState().GetAudi().GetPosition())
	var angle float64
	if distanceAudiBike > 10 {
		angle = math.Atan2(target.GetPosition().Y-bb.GetLocation().Y, target.GetPosition().X-bb.GetLocation().X)/math.Pi -
			bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
	} else {
		angle = math.Atan2(bb.GetLocation().Y-bb.GetGameState().GetAudi().GetPosition().Y, bb.GetLocation().X-bb.GetGameState().GetAudi().GetPosition().X)/math.Pi -
			bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
	}

	if angle > 1.0 {
		angle -= 2.0
	} else if angle < -1.0 {
		angle += 2.0
	}
	forces.Turning.SteerBike = true
	forces.Turning.SteeringForce = angle
	bb.SetForces(forces)
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 6 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	// TODO: implement this function
	bikeID := bb.GetBike()
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	for _, agent := range fellowBikers {
		if agent.GetID() == bb.GetID() {
			distribution[agent.GetID()] = 1.0
		} else {
			distribution[agent.GetID()] = 0.0
		}
	}
	return distribution
}

// only called when the agent is the dictator
func (bb *Agent8) DecideDictatorAllocation() voting.IdVoteMap {
	bikeID := bb.GetBike()
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	equalDist := 1.0 / float64(len(fellowBikers))
	for _, agent := range fellowBikers {
		distribution[agent.GetID()] = equalDist
	}
	return distribution
}

func (bb *Agent8) updateAgentActionMap() {

}

func (bb *Agent8) updateLoopScoreMap() {

}

// update the reputation for other agents
func (bb *Agent8) UpdateReputation() {
	// TODO: implement this function
}

// =========================================================================================================================================================

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIBaseBiker(baseBiker *objects.BaseBiker) objects.IBaseBiker {
	return &Agent8{
		BaseBiker: baseBiker,
	}
}
