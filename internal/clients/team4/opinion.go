package team4

import (
	"math"

	"github.com/google/uuid"
)

// calc sigmoid
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// get reputation value of all other agents
func (agent *BaselineAgent) GetReputation() map[uuid.UUID]float64 {
	return agent.reputation
}

// query for reputation value of specific agent with UUID
func (agent *BaselineAgent) QueryReputation(agentId uuid.UUID) float64 {
	return agent.reputation[agentId]
}

func (agent *BaselineAgent) QueryHonesty(agentId uuid.UUID) float64 {
	return agent.honestyMatrix[agentId]
}

func (agent *BaselineAgent) GetHonestyMatrix() map[uuid.UUID]float64 {
	return agent.reputation
}

func (agent *BaselineAgent) CalculateReputation() {
	megaBikes := agent.GetGameState().GetMegaBikes()
	decay_factor := 0.1
	totalReputationSum := float64(0)
	for _, bike := range megaBikes {
		fellowBikers := bike.GetAgents()
		//epsilon := 1e-10

		for _, otherAgent := range fellowBikers {
			selfTest := otherAgent.GetID()
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
				continue // Skip the rest of the loop for the current agent
			}

			historyenergy := agent.energyHistory[otherAgent.GetID()]
			lastEnergy := 1.0
			if len(historyenergy) > 2 {
				lastEnergy = historyenergy[len(historyenergy)-2]
			} else {
				lastEnergy = 0.0
			}
			energyLevel := otherAgent.GetEnergyLevel()
			consumption := energyLevel - lastEnergy

			myhistoryenergy := agent.energyHistory[agent.GetID()]
			mylastEnergy := 1.0
			if len(myhistoryenergy) > 2 {
				mylastEnergy = myhistoryenergy[len(myhistoryenergy)-2]
			} else {
				mylastEnergy = 0.0
			}
			myenergyLevel := agent.GetEnergyLevel()
			myconsumption := myenergyLevel - mylastEnergy
			EnergyReputation := (consumption / (energyLevel + 0.001)) - (myconsumption / (myenergyLevel + 0.001))

			//consumption / (energyLevel + epsilon)

			// Check if ReputationEnergy is NaN or Inf before proceeding

			bikeStatus := otherAgent.GetBikeStatus()
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}

			OverallReputation := EnergyReputation * ReputationBikeShift

			// Check if OverallReputation is NaN or Inf before storing
			if math.IsNaN(OverallReputation) || math.IsInf(OverallReputation, 0) {
				agent.reputation[otherAgent.GetID()] = 0.0
				continue // Skip the rest of the loop for the current agent
			}
			OverallReputation = sigmoid(OverallReputation)
			// print((1-decay_factor)*(agent.reputation[otherAgent.GetID()])+decay_factor*(OverallReputation))
			// print("\n")
			finalReputation := (1-decay_factor)*(agent.reputation[otherAgent.GetID()]) + decay_factor*(OverallReputation)
			agent.reputation[otherAgent.GetID()] = finalReputation
			totalReputationSum += finalReputation
		}
	}
	//normalize the reputation
	for _, bike := range megaBikes {
		fellowBikers := bike.GetAgents()
		for _, otherAgent := range fellowBikers {
			agent.reputation[otherAgent.GetID()] = agent.reputation[otherAgent.GetID()] / totalReputationSum
		}
	}
}

func (agent *BaselineAgent) CalculateHonestyMatrix() {
	if agent.honestyMatrix == nil {
		agent.honestyMatrix = make(map[uuid.UUID]float64)
	}

	for _, bike := range agent.GetGameState().GetMegaBikes() {
		for _, biker := range bike.GetAgents() {
			bikerID := biker.GetID()

			if _, exists := agent.honestyMatrix[bikerID]; !exists {
				agent.honestyMatrix[bikerID] = 1.0
			}
		}
	}
}

func (agent *BaselineAgent) DecreaseHonesty(agentID uuid.UUID, decreaseAmount float64) {
	if currentHonesty, ok := agent.honestyMatrix[agentID]; ok {
		newHonesty := currentHonesty - decreaseAmount
		if newHonesty < 0 {
			newHonesty = 0
		}
		agent.honestyMatrix[agentID] = newHonesty
	}
}

func (agent *BaselineAgent) IncreaseHonesty(agentID uuid.UUID, increaseAmount float64) {
	if currentHonesty, ok := agent.honestyMatrix[agentID]; ok {
		newHonesty := currentHonesty + increaseAmount
		if newHonesty > 1 {
			newHonesty = 1
		}
		agent.honestyMatrix[agentID] = newHonesty
	}
}
