// CALCULATES OPINIONS OF OTHER AGENTS

package team1

import (
	"math"

	"github.com/google/uuid"
)

type Opinion struct {
	effort   float64
	trust    float64
	fairness float64
	opinion  float64
}

// -----------------OPINION FUNCTIONS------------------

func (bb *Biker1) UpdateEffort(agentID uuid.UUID) {
	agent := bb.GetAgentFromId(agentID)
	fellowBikers := bb.GetFellowBikers()
	totalPedalForce := 0.0
	for _, agent := range fellowBikers {
		totalPedalForce = totalPedalForce + agent.GetForces().Pedal
	}
	avgForce := totalPedalForce / float64(len(fellowBikers))
	//effort expectation is scaled by their energy level -- should it be? (*agent.GetEnergyLevel())
	finalEffort := bb.opinions[agent.GetID()].effort + (agent.GetForces().Pedal-avgForce)*effortScaling

	if finalEffort > 1 {
		finalEffort = 1
	}
	if finalEffort < 0 {
		finalEffort = 0
	}
	newOpinion := Opinion{
		effort:   finalEffort,
		fairness: bb.opinions[agentID].fairness,
		trust:    bb.opinions[agentID].trust,
		opinion:  bb.opinions[agentID].opinion,
	}
	bb.opinions[agent.GetID()] = newOpinion
}

func (bb *Biker1) UpdateTrust(agentID uuid.UUID) {
	id := agentID
	agent := bb.GetAgentFromId(agentID)
	finalTrust := 0.5
	if agent.GetForces().Turning.SteeringForce == bb.GetForces().Turning.SteeringForce {
		finalTrust = bb.opinions[id].trust + deviatePositive
		if finalTrust > 1 {
			finalTrust = 1
		}
	} else {
		finalTrust := bb.opinions[id].trust + deviateNegative
		if finalTrust < 0 {
			finalTrust = 0
		}
	}
	newOpinion := Opinion{
		effort:   bb.opinions[id].effort,
		fairness: bb.opinions[id].fairness,
		trust:    finalTrust,
		opinion:  bb.opinions[id].opinion,
	}
	bb.opinions[id] = newOpinion
}

// func (bb *Biker1) UpdateFairness(agent obj.IBaseBiker) {
// 	difference := 0.0
// 	agentVote := agent.DecideAllocation()
// 	fairVote := bb.DecideAllocation()
// 	//If anyone has a better solution fo this please do it, couldn't find a better way to substract two maps in go
// 	for i, theirVote := range agentVote {
// 		for j, ourVote := range fairVote {
// 			if i == j {
// 				difference = difference + math.Abs(ourVote - theirVote)
// 			}
// 		}
// 	}
// 	finalFairness := bb.opinions[agent.GetID()].fairness + (fairnessDifference - difference/2)*fairnessScaling

// 	if finalFairness > 1 {
// 		finalFairness = 1
// 	}
// 	if finalFairness < 0 {
// 		finalFairness = 0
// 	}
// 	agentID := agent.GetID()
// 	newOpinion := Opinion{
// 		effort:   bb.opinions[agentID].effort,
// 		fairness: finalFairness,
// 		trust:    bb.opinions[agentID].trust,
// 		opinion:  bb.opinions[agentID].opinion,
// 	}
// 	bb.opinions[agent.GetID()] = newOpinion
// }

// how well does agent 1 like agent 2 according to objective metrics
func (bb *Biker1) GetObjectiveOpinion(id1 uuid.UUID, id2 uuid.UUID) float64 {
	agent1 := bb.GetAgentFromId(id1)
	agent2 := bb.GetAgentFromId(id2)
	objOpinion := 0.0
	if agent1.GetColour() == agent2.GetColour() {
		objOpinion = objOpinion + colorOpinionConstant
	}
	objOpinion = objOpinion + (agent1.GetEnergyLevel() - agent2.GetEnergyLevel())
	all_agents := bb.GetAllAgents()
	maxpoints := 0
	for _, agent := range all_agents {
		if agent.GetPoints() > maxpoints {
			maxpoints = agent.GetPoints()
		}
	}
	if maxpoints != 0 {
		objOpinion = objOpinion + float64((agent1.GetPoints()-agent2.GetPoints())/maxpoints)
	}
	objOpinion = math.Abs(objOpinion / (2.0 + colorOpinionConstant)) //normalise to 0-1
	return objOpinion
}

func (bb *Biker1) setOpinion() {
	if bb.opinions == nil {
		bb.opinions = make(map[uuid.UUID]Opinion)
	}
}

func (bb *Biker1) UpdateOpinions() {
	fellowBikers := bb.GetFellowBikers()
	bb.setOpinion()
	for _, agent := range fellowBikers {
		id := agent.GetID()
		_, ok := bb.opinions[agent.GetID()]

		if !ok {
			agentId := agent.GetID()
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:   0.5,
				trust:    0.5,
				fairness: 0.5,
				opinion:  0.5,
			}
			bb.opinions[agentId] = newOpinion
		}
		bb.UpdateTrust(id)
		bb.UpdateEffort(id)
		//bb.UpdateFairness(agent)
		bb.UpdateOpinion(id, 1.0)
	}
}

func (bb *Biker1) UpdateOpinion(id uuid.UUID, multiplier float64) {
	//Sorry no youre right, keep it, silly me
	bb.setOpinion()
	_, ok := bb.opinions[id]
	if !ok {
		//if we have no data on an agent, initialise to neutral
		newOpinion := Opinion{
			effort:   0.5,
			trust:    0.5,
			fairness: 0.5,
			opinion:  0.5,
		}
		bb.opinions[id] = newOpinion
	}

	newOpinion := Opinion{
		effort:   bb.opinions[id].effort,
		trust:    bb.opinions[id].trust,
		fairness: bb.opinions[id].fairness,
		opinion:  ((bb.opinions[id].trust*trustconstant + bb.opinions[id].effort*effortConstant + bb.opinions[id].fairness*fairnessConstant) / (trustconstant + effortConstant + fairnessConstant)) * multiplier,
	}

	if newOpinion.opinion > 1 {
		newOpinion.opinion = 1
	} else if newOpinion.opinion < 0 {
		newOpinion.opinion = 0
	}
	bb.opinions[id] = newOpinion

}

func (bb *Biker1) ourReputation() float64 {
	founding_agents := bb.GetAllAgents()
	reputation := 0.0
	for _, agent := range founding_agents {
		reputation = reputation + bb.GetObjectiveOpinion(bb.GetID(), agent.GetID())

	}
	reputation = reputation / float64(len(founding_agents))
	return reputation
}

// ----------------END OF OPINION FUNCTIONS--------------