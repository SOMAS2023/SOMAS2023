// CALCULATES OPINIONS OF OTHER AGENTS

package team1

import (
	obj "SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type Opinion struct {
	effort          float64
	trust           float64
	fairness        float64
	relativeSuccess float64
	// forgiveness float64
	opinion float64 // cumulative result of all the above
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

	if finalTrust > 1 {
		finalTrust = 1
	} else if finalTrust < 0 {
		finalTrust = 0
	}
	newOpinion := Opinion{
		effort:          bb.opinions[id].effort,
		fairness:        bb.opinions[id].fairness,
		trust:           finalTrust,
		relativeSuccess: bb.opinions[id].relativeSuccess,
		opinion:         bb.opinions[id].opinion,
	}
	bb.opinions[id] = newOpinion
}

func (bb *Biker1) UpdateFairness(agentID uuid.UUID) {
	helpfulAllocation := bb.getHelpfulAllocation()
	//for now just implement for democracy
	agent := bb.GetAgentFromId(agentID)
	energyChange := agent.GetEnergyLevel() - bb.prevEnergy[agentID] //how much of lootx distribution they got
	finalFairness := bb.opinions[agent.GetID()].fairness

	if energyChange-helpfulAllocation[agentID] > 0 {
		//they have more than they should have fairly got
		fmt.Printf("UNFAIR\n")
		finalFairness -= (energyChange - helpfulAllocation[agentID]) * fairnessScaling
	} else {
		fmt.Printf("FAIR\n")
		finalFairness += ((1 - (energyChange - helpfulAllocation[agentID])) / 2) * fairnessScaling
	}

	if finalFairness > 1 {
		finalFairness = 1
	} else if finalFairness < 0 {
		finalFairness = 0
	}

	newOpinion := Opinion{
		effort:          bb.opinions[agentID].effort,
		fairness:        finalFairness,
		trust:           bb.opinions[agentID].trust,
		relativeSuccess: bb.opinions[agentID].relativeSuccess,
		opinion:         bb.opinions[agentID].opinion,
	}
	bb.opinions[agentID] = newOpinion
}

func (bb *Biker1) UpdateRelativeSuccess(agentID uuid.UUID) {
	relativeSuccess := bb.GetRelativeSuccess(bb.GetID(), agentID)
	finalRelativeSuccess := bb.opinions[agentID].relativeSuccess + (relativeSuccess-bb.opinions[agentID].relativeSuccess)*relativeSuccessScaling
	if finalRelativeSuccess > 1 {
		finalRelativeSuccess = 1
	}
	if finalRelativeSuccess < 0 {
		finalRelativeSuccess = 0
	}
	newOpinion := Opinion{
		effort:          bb.opinions[agentID].effort,
		fairness:        bb.opinions[agentID].fairness,
		trust:           bb.opinions[agentID].trust,
		relativeSuccess: finalRelativeSuccess,
		opinion:         bb.opinions[agentID].opinion,
	}
	bb.opinions[agentID] = newOpinion
}

// how well does agent 1 like agent 2 according to objective metrics
func (bb *Biker1) GetRelativeSuccess(id1 uuid.UUID, id2 uuid.UUID) float64 {
	agent1 := bb.GetAgentFromId(id1)
	agent2 := bb.GetAgentFromId(id2)
	relativeSuccess := 0.0
	if agent1.GetColour() == agent2.GetColour() {
		relativeSuccess = relativeSuccess + colorOpinionConstant
	}
	relativeSuccess = relativeSuccess + (agent1.GetEnergyLevel() - agent2.GetEnergyLevel())
	all_agents := bb.GetAllAgents()
	maxpoints := 0
	for _, agent := range all_agents {
		if agent.GetPoints() > maxpoints {
			maxpoints = agent.GetPoints()
		}
	}
	if maxpoints != 0 {
		relativeSuccess = relativeSuccess + float64((agent1.GetPoints()-agent2.GetPoints())/maxpoints)
	}
	relativeSuccess = math.Abs(relativeSuccess / (2.0 + colorOpinionConstant)) //normalise to 0-1
	return relativeSuccess
}

func (bb *Biker1) UpdateOpinion(id uuid.UUID, multiplier float64) {
	_, ok := bb.opinions[id]
	if !ok {
		//if we have no data on an agent, initialise to neutral
		newOpinion := Opinion{
			effort:          0.5,
			trust:           0.5,
			fairness:        0.5,
			relativeSuccess: 0.5,
			opinion:         0.5,
		}
		bb.opinions[id] = newOpinion
	}

	newOpinion := Opinion{
		effort:          bb.opinions[id].effort,
		trust:           bb.opinions[id].trust,
		fairness:        bb.opinions[id].fairness,
		relativeSuccess: bb.opinions[id].relativeSuccess,
		opinion:         ((bb.opinions[id].trust*trustconstant + bb.opinions[id].effort*effortConstant + bb.opinions[id].fairness*fairnessConstant) / (trustconstant + effortConstant + fairnessConstant)) * multiplier,
	}

	if newOpinion.opinion > 1 {
		newOpinion.opinion = 1
	} else if newOpinion.opinion < 0 {
		newOpinion.opinion = 0
	}
	fmt.Printf("new opinion: %v\n", newOpinion)
	bb.opinions[id] = newOpinion

}

func (bb *Biker1) setOpinions() map[uuid.UUID]Opinion {
	if bb.opinions == nil {
		bb.opinions = make(map[uuid.UUID]Opinion)
		for _, agent := range bb.GetAllAgents() {
			agentId := agent.GetID()
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:          0.5,
				trust:           0.5,
				fairness:        0.5,
				relativeSuccess: 0.5,
				opinion:         0.5,
			}
			bb.opinions[agentId] = newOpinion
		}
	}
	return bb.opinions
}

// infer our reputation from the average relative success of agents in the current context
func (bb *Biker1) DetermineOurReputation() float64 {
	var agentsInContext []obj.IBaseBiker
	if bb.GetBike() == uuid.Nil {
		agentsInContext = bb.GetAllAgents()
	} else {
		agentsInContext = bb.GetFellowBikers()
	}

	reputation := 0.0
	for _, agent := range agentsInContext {
		reputation = reputation + bb.GetRelativeSuccess(bb.GetID(), agent.GetID())
	}
	fmt.Printf("Reputation: %v\n", reputation)
	reputation = reputation / float64(len(agentsInContext))
	return reputation
}

func (bb *Biker1) UpdateAllAgentsOpinions(agents_to_update []obj.IBaseBiker) {
	bb.setOpinions()
	for _, agent := range agents_to_update {
		id := agent.GetID()
		_, ok := bb.opinions[id]

		if !ok {
			agentId := agent.GetID()
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:          0.5,
				trust:           0.5,
				fairness:        0.5,
				relativeSuccess: 0.5,
				opinion:         0.5,
			}
			bb.opinions[agentId] = newOpinion
		}
		bb.UpdateOpinion(id, 1)
	}

}

func (bb *Biker1) UpdateAllAgentsEffort() {
	fellowBikers := bb.GetFellowBikers()
	fmt.Printf("Fellow bikers: %v\n", len(fellowBikers))
	bikeId := bb.GetBike()
	gs := bb.GetGameState()
	totalMass := utils.MassBike + float64(len(fellowBikers))*utils.MassBiker
	velocity := gs.GetMegaBikes()[bikeId].GetPhysicalState().Velocity
	drag_force := utils.DragCoefficient * math.Pow(velocity, 2)
	acceleration := gs.GetMegaBikes()[bikeId].GetPhysicalState().Acceleration
	total_force := acceleration * totalMass

	//resultant =

	// bike_force := utils.MassBike * acceleration
	// fmt.Printf("Bike force: %v\n", bike_force)
	//resultant = force - drag
	//force = mass*unknown acceleration?
	//actual acceleration is calculated from (totalpedalforce - CalcDrag(v)) / m
	//resultant = mass*actual acceleration
	totalPedalForce := total_force + drag_force
	fmt.Printf("Acceleration: %v\n", acceleration)
	fmt.Printf("Total pedal force: %v\n", totalPedalForce)
	fmt.Printf("Drag force: %v\n", drag_force)
	fmt.Printf("Total force: %v\n", total_force)
	fmt.Printf("ACTUAL BIKE FORCE: %v\n", gs.GetMegaBikes()[bikeId].GetForce())

	// Calculate force pedalled by everyone else
	remainingForce := totalPedalForce - bb.getPedalForce()
	effortProbability := make(map[uuid.UUID]float64) //probability that they are exc
	lootBoxes := bb.GetGameState().GetLootBoxes()
	totalEffort := 0.0
	for _, agent := range fellowBikers {
		id := agent.GetID()
		_, ok := bb.opinions[id]

		if !ok {
			agentId := agent.GetID()
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:          0.5,
				trust:           0.5,
				fairness:        0.5,
				relativeSuccess: 0.5,
				opinion:         0.5,
			}
			bb.opinions[agentId] = newOpinion
		}

		colourProb := 0.0
		fmt.Printf("Agent colour: %v\n", agent.GetColour())
		fmt.Printf("Recent decided: %v\n", bb.recentDecided)
		fmt.Printf("Lootboxes %v\n", lootBoxes)
		if agent.GetColour() != bb.recentDecidedColour {
			//probability should be high
			//for now set to 0.5 but later change based on how close the lootbox is to their colour lootbox
			colourProb += 0.3
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
		effortProbability[agentId] *= remainingForce
		agent := bb.GetAgentFromId(agentId)

		//effort expectation is scaled by their energy and compare to our effort
		fmt.Printf("Remaining force: %v\n", remainingForce)
		fmt.Printf("Total effort: %v\n", totalEffort)
		fmt.Printf("Current effort: %v\n", bb.opinions[agentId].effort)
		fmt.Printf("verified effort prob: %v\n", agent.GetForces())
		fmt.Printf("Effort probability: %v\n", effortProbability[agentId])
		fmt.Printf("Pedal force: %v\n", bb.getPedalForce())
		finalEffort := bb.opinions[agentId].effort + (effortProbability[agentId]-bb.getPedalForce())*effortScaling

		if finalEffort > 1 {
			finalEffort = 1
		}
		if finalEffort < 0 {
			finalEffort = 0
		}
		newOpinion := Opinion{
			effort:          finalEffort,
			fairness:        bb.opinions[agentId].fairness,
			trust:           bb.opinions[agentId].trust,
			relativeSuccess: bb.opinions[agentId].relativeSuccess,
			opinion:         bb.opinions[agentId].opinion,
		}
		bb.opinions[agent.GetID()] = newOpinion
	}

}

func (bb *Biker1) UpdateAllAgentsTrust(agents_to_update []obj.IBaseBiker) {
	bb.setOpinions()
	for _, agent := range agents_to_update {
		id := agent.GetID()
		_, ok := bb.opinions[id]

		if !ok {
			agentId := agent.GetID()
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:          0.5,
				trust:           0.5,
				fairness:        0.5,
				relativeSuccess: 0.5,
				opinion:         0.5,
			}
			bb.opinions[agentId] = newOpinion
		}
		bb.UpdateTrust(id)
		fmt.Printf("Agent %v trust: %v\n", id, bb.opinions[id].trust)
	}
}

func (bb *Biker1) UpdateAllAgentsFairness(agents_to_update []obj.IBaseBiker) {
	bb.setOpinions()
	for _, agent := range agents_to_update {
		id := agent.GetID()
		_, ok := bb.opinions[id]

		if !ok {
			agentId := agent.GetID()
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:          0.5,
				trust:           0.5,
				fairness:        0.5,
				relativeSuccess: 0.5,
				opinion:         0.5,
			}
			bb.opinions[agentId] = newOpinion
		}

		bikeID := bb.GetBike()
		governance := bb.GetGameState().GetMegaBikes()[bikeID].GetGovernance()
		fmt.Printf("Governance %v\n", governance)
		if governance == 0 {
			bb.UpdateFairness(id)
		} else {
			ruler := bb.GetGameState().GetMegaBikes()[bikeID].GetRuler()
			bb.UpdateFairness(ruler)
			return
		}
		fmt.Printf("Agent %v fairness: %v\n", id, bb.opinions[id].fairness)
	}
}

func (bb *Biker1) UpdateAllAgentsRelativeSuccess(agents_to_update []obj.IBaseBiker) {
	bb.setOpinions()
	for _, agent := range agents_to_update {
		id := agent.GetID()
		_, ok := bb.opinions[id]

		if !ok {
			agentId := agent.GetID()
			//if we have no data on an agent, initialise to neutral
			newOpinion := Opinion{
				effort:          0.5,
				trust:           0.5,
				fairness:        0.5,
				relativeSuccess: 0.5,
				opinion:         0.5,
			}
			bb.opinions[agentId] = newOpinion
		}
		bb.UpdateRelativeSuccess(id)
		fmt.Printf("Agent %v relative success: %v\n", id, bb.opinions[id].relativeSuccess)
	}
}

// ----------------END OF OPINION FUNCTIONS--------------
