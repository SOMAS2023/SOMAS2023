package team5Agent

import (
	// Assuming this package contains the IMegaBike interface

	"math"

	"github.com/google/uuid"
)

func (t5 *team5Agent) InitialiseReputation() {
	megaBikes := t5.GetGameState().GetMegaBikes()
	for _, mb := range megaBikes {
		// Iterate through all agents on each MegaBike
		for _, agent := range mb.GetAgents() {
			// Set initial reputation to 0.5 for each agent
			t5.SetReputation(agent.GetID(), 0.5)
		}
	}
}

// Most important 3 functions:

// Reputation calculation currently just based on energy and force
func (t5 *team5Agent) calculateReputationOfAgent(agentID uuid.UUID, currentRep float64) float64 {
	averagePedalForce := t5.getAverageForceOfAgents()
	averageEnergy := t5.getAverageEnergyOfAgents()

	agentPedalForce := t5.getForceOfOneAgent(agentID)
	agentEnergy := t5.getEnergyOfOneAgent(agentID)

	forceDeviation := agentPedalForce / averagePedalForce //fraction of agentMetric/averageMetric
	energyDeviation := agentEnergy / averageEnergy

	combinedDeviation := (forceDeviation + energyDeviation) / 2 // keeps it in range [0,1]

	// get current reputation of the agent

	weight := 0.2 //maximum change per round
	newRep := currentRep + (combinedDeviation-1)*weight
	return math.Min(math.Max(newRep, 0), 1) //capped at 0 and 1
}

func (t5 *team5Agent) updateReputationOfAllAgents() {
	for agentID, reputaion := range t5.GetReputation() {
		newRep := t5.calculateReputationOfAgent(agentID, reputaion)
		t5.SetReputation(agentID, newRep)
	}
}

//Useful helper functions:

func (t5 *team5Agent) getAveragePedalSpeedOfMegaBike(megaBikeID uuid.UUID) float64 {
	megaBikes := t5.GetGameState().GetMegaBikes()
	megaBike, exists := megaBikes[megaBikeID]
	if !exists {
		return 0
	}
	agents := megaBike.GetAgents()
	var totalPedalSpeed float64
	for _, agent := range agents {
		totalPedalSpeed += agent.GetForces().Pedal
	}
	return totalPedalSpeed / float64(len(agents))
}

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
	return totalForce / totalAgents
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
