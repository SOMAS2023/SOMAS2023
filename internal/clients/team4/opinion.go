package team4

import (
	"github.com/google/uuid"
)

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

// changed version
func (agent *BaselineAgent) CalculateReputation() {
	////////////////////////////
	//  As the program I used for debugging invoked "padal" and "break" with values of 0, I conducted tests using random numbers.
	// In case of an updated main program, I will need to adjust the parameters and expressions of the reputation matrix.
	// The current version lacks real data during the debugging process.
	////////////////////////////
	megaBikes := agent.GetGameState().GetMegaBikes()

	for _, bike := range megaBikes {
		// Get all agents on MegaBike
		fellowBikers := bike.GetAgents()

		// Iterate over each agent on MegaBike, generate reputation assessment
		for _, otherAgent := range fellowBikers {
			// Exclude self
			selfTest := otherAgent.GetID() //nolint
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
			}

			// Monitor otherAgent's location
			// location := otherAgent.GetLocation()
			// RAP := otherAgent.GetResourceAllocationParams()
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Location:", location, "ResourceAllocationParams:", RAP)

			// Monitor otherAgent's forces
			historyenergy := agent.energyHistory[otherAgent.GetID()]
			lastEnergy := 1.0
			if len(historyenergy) >= 2 {
				lastEnergy = historyenergy[len(historyenergy)-2]
				// rest of your code
			} else {
				lastEnergy = 0.0
			}
			energyLevel := otherAgent.GetEnergyLevel()
			ReputationEnergy := float64((lastEnergy)) / energyLevel //CAUTION: REMOVE THE RANDOM VALUE
			//print("我是大猴子")
			//fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Forces:", ReputationEnergy)

			// Monitor otherAgent's bike status
			bikeStatus := otherAgent.GetBikeStatus()
			// Convert the boolean value to float64 and print the result
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}
			//fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Bike_Shift", float64(ReputationBikeShift))

			// Calculate Overall_reputation
			OverallReputation := ReputationEnergy * ReputationBikeShift
			//fmt.Println("Agent ID:", otherAgent.GetID(), "Overall Reputation:", OverallReputation)

			// Store Overall_reputation in the reputation map
			agent.reputation[otherAgent.GetID()] = OverallReputation
		}
	}
	/* 	for agentID, agentReputation := range agent.reputation {
		print("Agent ID: ", agentID.String(), ", Reputation: ", agentReputation, "\n")
	} */

}

/* // Reputation and Honesty Matrix Teams Must Implement these or similar functions

func (agent *BaselineAgent) CalculateReputation() {
	////////////////////////////
	//  As the program I used for debugging invoked "padal" and "break" with values of 0, I conducted tests using random numbers.
	// In case of an updated main program, I will need to adjust the parameters and expressions of the reputation matrix.
	// The current version lacks real data during the debugging process.
	////////////////////////////
	megaBikes := agent.GetGameState().GetMegaBikes()

	for _, bike := range megaBikes {
		// Get all agents on MegaBike
		fellowBikers := bike.GetAgents()

		// Iterate over each agent on MegaBike, generate reputation assessment
		for _, otherAgent := range fellowBikers {
			// Exclude self
			selfTest := otherAgent.GetID() //nolint
			if selfTest == agent.GetID() {
				agent.reputation[otherAgent.GetID()] = 1.0
			}

			// Monitor otherAgent's location
			// location := otherAgent.GetLocation()
			// RAP := otherAgent.GetResourceAllocationParams()
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Location:", location, "ResourceAllocationParams:", RAP)

			// Monitor otherAgent's forces
			forces := otherAgent.GetForces()
			energyLevel := otherAgent.GetEnergyLevel()
			ReputationForces := float64(forces.Pedal+forces.Brake+rand.Float64()) / energyLevel //CAUTION: REMOVE THE RANDOM VALUE
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Forces:", ReputationForces)

			// Monitor otherAgent's bike status
			bikeStatus := otherAgent.GetBikeStatus()
			// Convert the boolean value to float64 and print the result
			ReputationBikeShift := 0.2
			if bikeStatus {
				ReputationBikeShift = 1.0
			}
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Reputation_Bike_Shift", float64(ReputationBikeShift))

			// Calculate Overall_reputation
			OverallReputation := ReputationForces * ReputationBikeShift
			// fmt.Println("Agent ID:", otherAgent.GetID(), "Overall Reputation:", OverallReputation)

			// Store Overall_reputation in the reputation map
			agent.reputation[otherAgent.GetID()] = OverallReputation
		}
	}
	// for agentID, agentReputation := range agent.reputation {
	// 	print("Agent ID: ", agentID.String(), ", Reputation: ", agentReputation, "\n")
	// }
} */

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
