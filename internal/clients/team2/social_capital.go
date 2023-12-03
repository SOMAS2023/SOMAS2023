package team2

import (
	"SOMAS2023/internal/common/utils"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// TODO: function CalculateSocialCapital
func (a *AgentTwo) CalculateSocialCapital() {
	// Implement this method
	// Hardcode the weightings for now: Trust 1, Institution 0, Network 0
	// Calculate social capital of all agents
	// Calculate Reputation of all agents
	// Calculate social networks of all agents
	// Calculate institutions of all agents
	// Iterate over each agent and calculate their social capital

	for agentID := range a.Reputation {
		reputation := a.Reputation[agentID]
		institution := a.Institution[agentID]
		network := a.Network[agentID] // Assuming these values are already calculated

		newSocialCapital := TrustWeight*reputation + InstitutionWeight*institution + NetworkWeight*network

		// if the current socialCapital is smaller than the previous socialCapital AND the forgiveness counter is less than or equal to 3, then we increase the forgiveness counter
		if a.SocialCapital[agentID] > newSocialCapital && a.forgivenessCounter <= 3 { // If they were trustworthy in prev rounds, we feel remorse and we forgive them
			a.forgivenessCounter++
			a.SocialCapital[agentID] = (a.SocialCapital[agentID]*float64(a.GameIterations) + (newSocialCapital + forgivenessFactor*(newSocialCapital-a.SocialCapital[agentID]))) / (float64(a.GameIterations) + 1)
		} else if a.forgivenessCounter > 3 {
			// More than 3 rounds of BETRAYAL, we don't forgive them anymore...
			a.SocialCapital[agentID] = (a.SocialCapital[agentID]*float64(a.GameIterations) + newSocialCapital) / (float64(a.GameIterations) + 1)
		} else {
			// Good action with high trustworthiness
			a.forgivenessCounter = 0
			a.SocialCapital[agentID] = (a.SocialCapital[agentID]*float64(a.GameIterations) + newSocialCapital) / (float64(a.GameIterations) + 1)
		}
	}

	fmt.Println("Social Capital: ", a.SocialCapital)
}

func (a *AgentTwo) updateReputation(agentID uuid.UUID, ourDesiredLootbox uuid.UUID, theirDesiredLootbox uuid.UUID) {
	// Compare our desired lootbox with their desired lootbox
	// We retain a moving average of their reputation to not drastically make a change
	// If they are the same, we increase their reputation

	if ourDesiredLootbox == theirDesiredLootbox {
		// If they are the same, we increase their reputation
		a.Reputation[agentID] = (a.Reputation[agentID]*float64(a.GameIterations) + 1) / (float64(a.GameIterations) + 1)
	} else {
		// If they are different, we decrease their reputation
		a.Reputation[agentID] = (a.Reputation[agentID]*float64(a.GameIterations) - 1) / (float64(a.GameIterations) + 1)
	}

	fmt.Println("Reputation: ", a.Reputation)
}

// Get the direction to the voted lootbox
func (a *AgentTwo) GetVotedLootboxForces(lootboxID uuid.UUID) utils.Forces {
	lootbox := a.gameState.GetLootBoxes()[lootboxID]
	lootboxPositionX, lootboxPositionY := lootbox.GetPosition().X, lootbox.GetPosition().Y
	agentPositionX, agentPositionY := a.GetLocation().X, a.GetLocation().Y
	deltaX := lootboxPositionX - agentPositionX
	deltaY := lootboxPositionY - agentPositionY
	angle := math.Atan2(deltaY, deltaX)
	normalisedAngle := angle / math.Pi
	turningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: normalisedAngle - a.gameState.GetMegaBikes()[a.GetBike()].GetOrientation(),
	}
	return utils.Forces{
		Pedal:   utils.BikerMaxForce,
		Brake:   0.0,
		Turning: turningDecision,
	}
}

// Called by Events to obtain Event Value for update Institution
// Assume what they broadcast is the truth
// TODO: Obtain actual action performed from messaging
// 1. Rule Adhereance (Follow leader biker/ dictator)
func (a *AgentTwo) RuleAdhereanceValue(agentID uuid.UUID, expectedAction, actualAction utils.Forces) float64 {

	actualVector := forcesToVectorConversion(actualAction)
	expectedVector := forcesToVectorConversion(expectedAction)
	similarity := cosineSimilarity(actualVector, expectedVector)

	forceApplied := actualAction.Pedal
	return similarity * forceApplied
}

// Assume institution rules is only broadcasted within the same bike
// Events to update Institution
// 1. Rule Adhereance (Follow leader biker/ dictator )
// 2. Voting
// 3. Kicked out of bike
// 4. Accepted to bike
// 5. Role Assignment (Voted to be leader/ Dictator)
func (a *AgentTwo) updateInstitution(agentID uuid.UUID, weight float64, EventValue float64) {
	a.Institution[agentID] += EventValue * weight
}

// func (a *AgentTwo) updateNetwork(agentID uuid.UUID) float64 {
// 	// return 0.5 // This is just a placeholder value
// }

// func (a *AgentTwo) calculateTrustworthiness(agentID uuid.UUID) float64 {

// 	return 0.5 // This is just a placeholder value
// }

// func (a *AgentTwo) calculateInstitution(agentID uuid.UUID) float64 {

// 	// return 0.5 // This is just a placeholder value
// }

// func (a *AgentTwo) calculateNetwork(agentID uuid.UUID) float64 {
// 	// return 0.5 // This is just a placeholder value
// }
