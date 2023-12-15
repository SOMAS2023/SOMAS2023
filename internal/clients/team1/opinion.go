// CALCULATES OPINIONS OF OTHER AGENTS

package team1

import (
	obj "SOMAS2023/internal/common/objects"
	"math"

	"github.com/google/uuid"
)

type Opinion struct {
	effort   float64
	trust    float64
	fairness float64
	opinion  float64 // cumulative result of all the above
}

// -----------------OPINION FUNCTIONS------------------

// Update an agent's trust value
func (bb *Biker1) UpdateTrust(agentID uuid.UUID) {
	id := agentID
	agent := bb.GetAgentFromId(id)
	finalTrust := bb.opinions[id].trust //nothing changes
	targetPos := bb.recentDecidedPosition
	currLocation := bb.GetLocation()
	deltaX := targetPos.X - currLocation.X
	deltaY := targetPos.Y - currLocation.Y
	angle := math.Atan2(deltaY, deltaX)
	normalisedAngle := angle / math.Pi
	steeringAngle := normalisedAngle - bb.GetBikeInstance().GetOrientation()
	if math.Abs(steeringAngle) < 0.01 { //we are headed in direction towards lootbox
		finalTrust = bb.opinions[id].trust + deviatePositive //will change to be based on weighting
	} else {
		//	need to estimate likelihood of each agent deviating from the correct steeringAngle
		if agent.GetColour() != bb.recentDecidedColour {
			//currently if its not the agent's colour then trust in them decreases
			//needs to include reputation somehow
			//needs to calculate orientation to their colour (is it closer to or further than (orientation wise) voted lootbox)
			finalTrust = bb.opinions[id].trust - deviateNegative
		}
	}

	// We update the trust value too of the agent from direct experience from messaging session

	if finalTrust > 1 {
		finalTrust = 1
	} else if finalTrust < 0 {
		finalTrust = 0
	}
	newOpinion := Opinion{
		effort:   bb.opinions[id].effort,
		fairness: bb.opinions[id].fairness,
		trust:    finalTrust,
		opinion:  bb.opinions[id].opinion,
	}
	bb.opinions[id] = newOpinion
}

// Update an agent's fairness value
func (bb *Biker1) UpdateFairness(agentID uuid.UUID) {
	helpfulAllocation := bb.getHelpfulAllocation()
	//for now just implement for democracy
	agent := bb.GetAgentFromId(agentID)
	energyChange := agent.GetEnergyLevel() - bb.prevEnergy[agentID] //how much of lootx distribution they got
	finalFairness := bb.opinions[agent.GetID()].fairness
	if energyChange-helpfulAllocation[agentID] > 0 {
		//they have more than they should have fairly got
		finalFairness -= (energyChange - helpfulAllocation[agentID]) * fairnessScaling
	} else {
		finalFairness += ((1 - (energyChange - helpfulAllocation[agentID])) / 2) * fairnessScaling
	}

	if finalFairness > 1 {
		finalFairness = 1
	} else if finalFairness < 0 {
		finalFairness = 0
	}

	newOpinion := Opinion{
		effort:   bb.opinions[agentID].effort,
		fairness: finalFairness,
		trust:    bb.opinions[agentID].trust,
		// relativeSuccess: bb.opinions[agentID].relativeSuccess,
		opinion: bb.opinions[agentID].opinion,
	}
	bb.opinions[agentID] = newOpinion
}

// how well does agent 1 like agent 2 according to objective metrics
func (bb *Biker1) GetRelativeSuccess(id1 uuid.UUID, id2 uuid.UUID, all_agents []obj.IBaseBiker) float64 {
	agent1 := bb.GetAgentFromId(id1)
	agent2 := bb.GetAgentFromId(id2)
	relativeSuccess := 0.0
	// Colour comparison
	if agent1.GetColour() == agent2.GetColour() {
		relativeSuccess = relativeSuccess + colorOpinionConstant
	}
	//Energy comparison
	relativeSuccess = relativeSuccess + (agent1.GetEnergyLevel() - agent2.GetEnergyLevel())
	maxpoints := 0
	for _, agent := range all_agents {
		if agent.GetPoints() > maxpoints {
			maxpoints = agent.GetPoints()
		}
	}
	// Points comparison
	if maxpoints != 0 {
		relativeSuccess = relativeSuccess + float64((agent1.GetPoints()-agent2.GetPoints())/maxpoints)
	}
	relativeSuccess = math.Abs(relativeSuccess / (2.0 + colorOpinionConstant)) //normalise to 0-1
	return relativeSuccess
}

// Update our agent's of another agent
func (bb *Biker1) UpdateOpinion(id uuid.UUID, multiplier float64) {
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
		// relativeSuccess: bb.opinions[id].relativeSuccess,
		opinion: ((bb.opinions[id].trust*trustconstant + bb.opinions[id].effort*effortConstant +
			bb.opinions[id].fairness*fairnessConstant) / (trustconstant + effortConstant + fairnessConstant)) * multiplier,
	}

	if newOpinion.opinion > 1 {
		newOpinion.opinion = 1
	} else if newOpinion.opinion < 0 {
		newOpinion.opinion = 0
	}
	bb.opinions[id] = newOpinion

}

// Initialise/Set the opinions of all agents
func (bb *Biker1) setOpinions() map[uuid.UUID]Opinion {
	for _, agent := range bb.GetAllAgents() {
		agentId := agent.GetID()
		_, ok := bb.opinions[agentId]
		if !ok {
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:   0.5,
				trust:    0.5,
				fairness: 0.5,
				opinion:  0.5,
			}
			bb.opinions[agentId] = newOpinion
		}
	}
	return bb.opinions
}

// infer our reputation from the average relative success of agents in the current context
func (bb *Biker1) DetermineOurAverageReputation() float64 {
	var agentsInContext []obj.IBaseBiker
	var numberOnBike float64
	if bb.GetBike() == uuid.Nil {
		agentsInContext = bb.GetAllAgents()
	} else {
		agentsInContext = bb.GetFellowBikers()
	}
	if len(agentsInContext) == 0 {
		numberOnBike = 1
	} else {
		numberOnBike = float64(len(agentsInContext))
	}

	reputation := 0.0
	for _, agent := range agentsInContext {
		// bb.UpdateRelativeSuccess(agent.GetID(), agentsInContext)
		// fmt.Printf("Agent %v relative success: %v\n", agent.GetID(), bb.GetRelativeSuccess(bb.GetID(), agent.GetID(), agentsInContext))
		reputation = reputation + bb.GetRelativeSuccess(bb.GetID(), agent.GetID(), agentsInContext)
	}
	reputation = reputation / numberOnBike
	return reputation
}

// Update all agents' opinions
func (bb *Biker1) UpdateAllAgentsOpinions(agents_to_update []obj.IBaseBiker) {
	bb.setOpinions()
	for _, agent := range agents_to_update {
		id := agent.GetID()
		_, ok := bb.opinions[id]

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
		bb.UpdateOpinion(id, 1)
	}

}

// Update all agents' efforts
func (bb *Biker1) UpdateAllAgentsEffort() {
	fellowBikers := bb.GetFellowBikers()

	// totalPedalForce := total_force + drag_force

	fellowBikersExpendedEnergy := make(map[uuid.UUID]float64)
	totalExpendedEnergy := 0.0
	for _, agent := range fellowBikers {
		AgentExpendedEnergy := (bb.prevEnergy[agent.GetID()] - agent.GetEnergyLevel())
		fellowBikersExpendedEnergy[agent.GetID()] = AgentExpendedEnergy
		totalExpendedEnergy += AgentExpendedEnergy
	}
	OurExpendedEnergy := bb.prevEnergy[bb.GetID()] - bb.GetEnergyLevel()

	effortProbability := make(map[uuid.UUID]float64) //probability that they are exc
	totalEffort := 0.0
	for _, agent := range fellowBikers {
		id := agent.GetID()
		_, ok := bb.opinions[id]

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

		colourProb := 0.0
		if agent.GetColour() != bb.recentDecidedColour {
			//probability should be high
			//for now set to 0.5 but later change based on how close the lootbox is to their colour lootbox
			colourProb += 0.35
		}
		energyProb := 1 - agent.GetEnergyLevel()
		//Will add weightings to this so that energy probability has a lower weighting than difference in colour for example
		//also plus reputation

		effortProb := 1 - (colourProb+energyProb)/2 //scales between 0 and 1 and then negative so that higher probabilities mean you are less likely to contribute to pedal force
		effortProbability[agent.GetID()] = effortProb
		totalEffort += effortProb
	}
	for agentId := range effortProbability {
		//normalise effort probabilities
		effortProbability[agentId] /= totalEffort
		effortProbability[agentId] *= totalExpendedEnergy
		agent := bb.GetAgentFromId(agentId)
		// Effort = Current effort + trust * effort probability * (fellow biker's expended energy - our expended energy)
		finalEffort := bb.opinions[agentId].effort + bb.opinions[agentId].trust*effortProbability[agentId]*(fellowBikersExpendedEnergy[agentId]-OurExpendedEnergy)*effortScaling
		// fmt.Printf("Final effort: %v\n", finalEffort)

		if finalEffort > 1 {
			finalEffort = 1
		}
		if finalEffort < 0 {
			finalEffort = 0
		}
		newOpinion := Opinion{
			effort:   finalEffort,
			fairness: bb.opinions[agentId].fairness,
			trust:    bb.opinions[agentId].trust,
			opinion:  bb.opinions[agentId].opinion,
		}
		bb.opinions[agent.GetID()] = newOpinion
	}

}

// Update all agents' trust
func (bb *Biker1) UpdateAllAgentsTrust(agents_to_update []obj.IBaseBiker) {
	bb.setOpinions()
	for _, agent := range agents_to_update {
		id := agent.GetID()
		_, ok := bb.opinions[id]

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
	}
}

// Update all agents' fairness
func (bb *Biker1) UpdateAllAgentsFairness(agents_to_update []obj.IBaseBiker) {
	bb.setOpinions()
	for _, agent := range agents_to_update {
		id := agent.GetID()
		_, ok := bb.opinions[id]

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

		bikeID := bb.GetBike()
		governance := bb.GetGameState().GetMegaBikes()[bikeID].GetGovernance()
		if governance == 0 {
			bb.UpdateFairness(id)
		} else {
			ruler := bb.GetGameState().GetMegaBikes()[bikeID].GetRuler()
			bb.UpdateFairness(ruler)
			return
		}
	}
}

// ----------------END OF OPINION FUNCTIONS--------------
