package team5Agent

import (
	"fmt"

	"github.com/google/uuid"
)

func (t5 *team5Agent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	agentsOnBike := t5.GetAgentsOnMegaBike()
	decisions := make(map[uuid.UUID]bool)
	threshold := 0.5

	agentRep := NewRepSystem(t5.GetGameState())
	agentRep.updateReputationOfAllAgents()

	for _, agent := range agentsOnBike {

		key := agent.GetID()
		value := agentRep.calculateReputationOfAgent(key)

		fmt.Println(value)

		if value <= threshold {
			decisions[key] = false
		} else {
			decisions[key] = true
		}

	}

	return decisions

}
