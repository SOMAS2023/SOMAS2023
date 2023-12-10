package team5Agent

import (
	// Assuming this package contains the IMegaBike interface

	"math"

	"github.com/google/uuid"
)

func (t5 *team5Agent) InitialiseReputation() {
	//fmt.Println("HAHAHA: ", t5.GetReputation())
	megaBikes := t5.GetGameState().GetMegaBikes()
	for _, mb := range megaBikes {
		// Iterate through all agents on each MegaBike
		for _, agent := range mb.GetAgents() {
			// Set initial reputation to 0.5 for each agent
			t5.SetReputation(agent.GetID(), 0.5)
		}
	}
	//fmt.Println("HAHAHA22: ", t5.GetReputation())

}

// Most important 3 functions:

// Reputation calculation currently just based on energy and force
func (t5 *team5Agent) calculateReputationOfAgent(agentID uuid.UUID, currentRep float64) float64 {
	//Colour of agent
	//check energy allocation -> change of energy in each agent
	//if bike speed slow - lower everyone by small amount
	//if direction wrong a lot - lower everyone by small amount
	//Increase forgivenesss rate if in ultristic state
	//Msging stuff - increase rep of people who msg
	averageEnergy := t5.getAverageEnergyOfAgents()

	forgivenessRate := 0.0005 //Reputation slowly goes back to average over time = forgiveness.
	if t5.state == 3 {
		forgivenessRate += 0.0003
	}
	colourRep := 0.0
	//get all agent colours of all agents and check if they are the same as the agentID
	//if agent exists on map and agent colour is the same as the agentID then add 0.01 to colourRep
	if agent, ok := t5.GetGameState().GetAgents()[agentID]; ok {
		if (t5.GetColour()) == agent.GetColour() {
			colourRep = 0.01
		}
	}

	agentEnergy := t5.getEnergyOfOneAgent(agentID)
	energyDeviation := agentEnergy - averageEnergy
	//fmt.Println("anget energy: ", agentEnergy, " average energy: ", averageEnergy, " energy deviation: ", energyDeviation)
	combinedDeviation := energyDeviation //(forceDeviation + energyDeviation) / 2 // keeps it in range [0,1]

	// get current reputation of the agent

	weight := 0.2
	newRep := currentRep + combinedDeviation*weight + colourRep
	if newRep > 0.5 {
		newRep = newRep - forgivenessRate
	} else if newRep < 0.5 {
		newRep = newRep + forgivenessRate
	}
	rValue := math.Min(math.Max(newRep, 0.0), 1.0)
	//fmt.Println("Reputation of agent: ", agentID, " is: ", rValue)
	return rValue //capped at 0 and 1 (our internal reputation system is 0 to 1 not -1 to 1)

}

func (t5 *team5Agent) updateReputationOfAllAgents() {
	// if all agents have a reputation of 0 then update all to have a reputation of 0.5
	reputationMap := t5.GetReputation()

	if len(reputationMap) == 0 {
		t5.InitialiseReputation()
	}

	for agentID, reputation := range reputationMap {

		//if reuptation is NaN then set to 0.5
		if !(0 <= reputation && reputation <= 1) {
			t5.SetReputation(agentID, 0.5)
			reputation = 0.5
		}
		reputation = t5.calculateReputationOfAgent(agentID, reputation)
		if reputation > 0.5 {
			reputation -= 0.01
		} else if reputation < 0.5 {
			reputation += 0.01
		}
		t5.SetReputation(agentID, reputation)
	}
	t5.determineGreed()
}

func (t5 *team5Agent) determineGreed() {
	repMap := t5.CalculateEnergyChange(t5.GetBike())
	myEnergyChange := repMap[t5.GetID()]
	for agentID, energyChange := range repMap {
		if energyChange > myEnergyChange {
			t5.SetReputation(agentID, t5.QueryReputation(agentID)-0.05)
		} else if energyChange < myEnergyChange {
			t5.SetReputation(agentID, t5.QueryReputation(agentID)+0.05)
		}
	}
}

//Useful helper functions:

// func (t5 *team5Agent) getAveragePedalSpeedOfMegaBike(megaBikeID uuid.UUID) float64 {
// 	megaBikes := t5.GetGameState().GetMegaBikes()
// 	megaBike, exists := megaBikes[megaBikeID]
// 	if !exists {
// 		return 0
// 	}
// 	agents := megaBike.GetAgents()
// 	var totalPedalSpeed float64
// 	for _, agent := range agents {
// 		totalPedalSpeed += agent.GetForces().Pedal
// 	}
// 	return totalPedalSpeed / float64(len(agents))
// }

// Functions used in calculating the reputation value:

func (t5 *team5Agent) getReputationOfSingleBike(megaBikeID uuid.UUID) float64 {
	megaBikes := t5.GetGameState().GetMegaBikes()
	megaBike, exists := megaBikes[megaBikeID]
	if !exists {
		return 0
	}
	agents := megaBike.GetAgents()
	var totalReputation float64
	for _, agent := range agents {
		totalReputation += t5.GetReputation()[agent.GetID()]
	}
	return totalReputation / float64(len(agents))
}

func (t5 *team5Agent) getReputationOfAllBikes() map[uuid.UUID]float64 {
	megaBikes := t5.GetGameState().GetMegaBikes()
	reputations := make(map[uuid.UUID]float64)
	for megaBikeID := range megaBikes {
		reputations[megaBikeID] = t5.getReputationOfSingleBike(megaBikeID)
	}
	return reputations
}

func (t5 *team5Agent) getAverageEnergyOfAgents() float64 {
	megaBikes := t5.GetGameState().GetMegaBikes()
	var totalEnergy float64
	var totalAgents float64
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			totalEnergy += agent.GetEnergyLevel()
			totalAgents++
		}
	}
	return totalEnergy / totalAgents
}

func (t5 *team5Agent) getAverageForceOfAgents() float64 {
	megaBikes := t5.GetGameState().GetMegaBikes()
	var totalForce float64
	var totalAgents float64
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			forceOfAgent := agent.GetForces().Pedal
			if forceOfAgent > 0 { //only add force if agent is pedalling
				totalForce += forceOfAgent
				totalAgents++
			}
		}
	}
	//print("totalForce: ", totalForce, "totalAgents: ", totalAgents)
	//if naan then return 0
	// avgForce := totalForce / totalAgents
	// if avgForce > 0 {
	// 	return avgForce
	// }
	return 1
}

func (t5 *team5Agent) getEnergyOfOneAgent(agentID uuid.UUID) float64 {
	megaBikes := t5.GetGameState().GetMegaBikes()
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			if agent.GetID() == agentID {
				return agent.GetEnergyLevel()
			}
		}
	}
	return 0
}

func (t5 *team5Agent) getForceOfOneAgent(agentID uuid.UUID) float64 {
	megaBikes := t5.GetGameState().GetMegaBikes()
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			if agent.GetID() == agentID {
				return agent.GetForces().Pedal
			}
		}
	}
	return 0
}
