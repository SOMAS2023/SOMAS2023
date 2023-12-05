package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"math/rand"

	"SOMAS2023/internal/common/voting"
	"math"

	"github.com/MattSScott/basePlatformSOMAS/messaging"
	"github.com/google/uuid"
)

type GP struct {
	EnergyThreshold              float64
	DistanceThresholdForVoting   float64
	ThresholdForJoiningDecision  float64
	ThresholdForChangingMegabike float64
}

var GlobalParameters GP = GP{
	EnergyThreshold:              0.5,
	DistanceThresholdForVoting:   30,
	ThresholdForJoiningDecision:  0.2,
	ThresholdForChangingMegabike: 0.3,
}

type IBaselineAgent interface {
	objects.IBaseBiker
}

type Agent8 struct {
	*objects.BaseBiker
	overallScores voting.LootboxVoteMap //rank score for the lootbox
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
	for i, fellowBiker := range fellowBikers {
		if i == 0 {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	return votes
}

// defaults to voting for first agent in the list
func (bb *Agent8) VoteLeader() voting.IdVoteMap {
	// TODO: implement this function
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()
	for i, fellowBiker := range fellowBikers {
		if i == 0 {
			votes[fellowBiker.GetID()] = 1.0
		} else {
			votes[fellowBiker.GetID()] = 0.0
		}
	}
	return votes
}

// ===============================================================================================================================================================
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Message System <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// This function updates all the messages for that agent i.e. both sending and receiving.
// And returns the new messages from other agents to your agent
func (bb *Agent8) GetAllMessages([]objects.IBaseBiker) []messaging.IMessage[objects.IBaseBiker] {
	// For team's agent add your own logic on chosing when your biker should send messages
	wantToSendMsg := false
	if wantToSendMsg {
		reputationMsg := bb.CreateReputationMessage()
		kickOffMsg := bb.CreateKickOffMessage()
		lootboxMsg := bb.CreateLootboxMessage()
		joiningMsg := bb.CreateJoiningMessage()
		governceMsg := bb.CreateGoverenceMessage()
		forcesMsg := bb.CreateForcesMessage()
		return []messaging.IMessage[objects.IBaseBiker]{reputationMsg, kickOffMsg, lootboxMsg, joiningMsg, governceMsg, forcesMsg}
	}
	return []messaging.IMessage[objects.IBaseBiker]{}
}

func (bb *Agent8) HandleKickOffMessage(msg objects.KickOffAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// kickOff := msg.KickOff
}

func (bb *Agent8) HandleReputationMessage(msg objects.ReputationOfAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// reputation := msg.Reputation
}

func (bb *Agent8) HandleJoiningMessage(msg objects.JoiningAgentMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// bikeId := msg.BikeId
}

func (bb *Agent8) HandleLootboxMessage(msg objects.LootboxMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// lootboxId := msg.LootboxId
}

func (bb *Agent8) HandleGovernanceMessage(msg objects.GovernanceMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// bikeId := msg.BikeId
	// governanceId := msg.GovernanceId
}

func (bb *Agent8) HandleForcesMessage(msg objects.ForcesMessage) {
	// Team's agent should implement logic for handling other biker messages that were sent to them.

	// sender := msg.BaseMessage.GetSender()
	// agentId := msg.AgentId
	// agentForces := msg.AgentForces

}

//===============================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 1 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) VoteForKickout() map[uuid.UUID]int {
	// TODO: implement this function
	voteResults := make(map[uuid.UUID]int)
	bikeID := bb.GetBike()

	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if agentID != bb.GetID() {
			// random votes to other agents
			voteResults[agentID] = rand.Intn(2) // randomly assigns 0 or 1 vote
		}
	}

	return voteResults
}

// only called when the agent is the dictator
func (bb *Agent8) DecideKickOut() []uuid.UUID {
	// TODO: implement this function
	return (make([]uuid.UUID, 0))
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

	// Iterate through each bike
	for bikeID, megabike := range megaBikes {
		// Calculate the Borda score for the current bike
		bordaScore := bb.CalculateAverageEnergy(bikeID) +
			float64(bb.CountAgentsWithSameColour(bikeID)) +
			CalculateGiniIndexFromAB(float64(bb.CountAgentsWithSameColour(bikeID)), float64(len(megabike.GetAgents())))

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

//===============================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 2 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) DecideAction() objects.BikerAction {

	var energyLevels []float64
	var target_goal int
	energy_threshold := GlobalParameters.EnergyThreshold
	changingbike_threshold := GlobalParameters.ThresholdForChangingMegabike

	// utility should be calculated by the fomula outlined on page7 of Lec6
	utilityLevels := []float64{80.0, 90.0, 75.0, 85.0, 45.0, 35.0, 60.0, 70.0, 65.0}
	turns := []bool{true, false, true, true, false, true, true, false, true}
	decisions := []bool{true, false, false, true, true, false, false, true, false}

	// get the energy level of all agents in the megabike
	fellowBikers := bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetAgents()

	for _, agent := range fellowBikers {
		if !(agent.GetID() == bb.GetID()) {
			energy_level := agent.GetEnergyLevel()
			energyLevels = append(energyLevels, energy_level)
		}
	}

	goalPreferenceList := make([]int, len(energyLevels))
	if bb.GetEnergyLevel() >= energy_threshold {
		target_goal = 1
	} else {
		target_goal = 0
	}

	// Convert energyLevels to 0 or 1 based on the threshold
	for i, energy := range energyLevels {
		if energy >= energy_threshold {
			goalPreferenceList[i] = 1
		} else {
			goalPreferenceList[i] = 0
		}
	}

	// Find quantified ‘Value-judgement’
	valueJudgement := bb.calculateValueJudgement(utilityLevels, goalPreferenceList, target_goal, turns)

	// Scale the ‘Cost in the collective improvement’
	AverageOfCost := bb.calculateAverageOfCostAndPercentage(decisions, energyLevels, energy_threshold)

	// Find the overall ‘changeBike’ coefficient
	changeBikeCoefficient := 0.6*valueJudgement - 0.4*AverageOfCost

	// Make a decision based on the calculated coefficients
	if changingbike_threshold > changeBikeCoefficient {
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
		energyWeighting := calculateEnergyWeighting(bb.GetEnergyLevel())

		preferences[lootBox.GetID()] = colorPreference + (GlobalParameters.DistanceThresholdForVoting-distance)*energyWeighting
	}

	// Apply softmax to convert preferences to a probability distribution
	softmaxPreferences := softmax(preferences)

	// Rank loot boxes based on preferences
	rankedLootBoxes := rankByPreference(softmaxPreferences)

	bb.overallScores = softmaxPreferences

	return rankedLootBoxes[0]
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 4 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (bb *Agent8) DictateDirection() uuid.UUID {
	// TODO: implement this function
	return bb.GetID()
}

func (bb *Agent8) LeadDirection() uuid.UUID {
	// TODO: implement this function
	return bb.GetID()
}

// defaults to an equal distribution over all agents for all actions
func (bb *Agent8) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	weights := make(map[uuid.UUID]float64)
	agents := bb.GetFellowBikers()
	for _, agent := range agents {
		weights[agent.GetID()] = 1.0
	}
	return weights
}

// Multi-voting system
func (bb *Agent8) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	// Calculate the biker's individual preference scores
	preferenceScores := bb.calculatePreferenceScores(proposals)

	combinedScores := make(map[uuid.UUID]float64)
	for _, proposal := range proposals {
		combinedScore := preferenceScores[proposal] + bb.overallScores[proposal]
		combinedScores[proposal] = combinedScore
	}
	softmaxScores := softmax(combinedScores)

	return softmaxScores
}

// =========================================================================================================================================================

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> stage 5 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	// TODO: implement this function
	var forces utils.Forces
	forces.Brake = 0
	forces.Pedal = 1
	lootboxs := bb.GetGameState().GetLootBoxes()
	var target objects.ILootBox
	for key, value := range lootboxs {
		if key == direction {
			target = value
			break
		}
	}
	angle := math.Atan2(target.GetPosition().Y-bb.GetLocation().Y, target.GetPosition().X-bb.GetLocation().X)/math.Pi - bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
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

// update the reputation for other agents
func (bb *Agent8) UpdateReputation() {
	// TODO: implement this function
}

// =========================================================================================================================================================

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIBaseBiker(totColours utils.Colour, bikeId uuid.UUID) objects.IBaseBiker {
	return &Agent8{
		BaseBiker: objects.GetBaseBiker(totColours, bikeId),
	}
}
