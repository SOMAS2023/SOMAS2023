package team5Agent

import (
	"fmt"

	"github.com/google/uuid"
)

func (t5 *team5Agent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	bikeId := t5.BaseBiker.GetMegaBikeId()
	agentsOnBike := t5.BaseBiker.GetGameState().GetMegaBikes()[bikeId].GetAgents()
	decisions := make(map[uuid.UUID]bool)
	threshold := 0.5

	agentRep := NewRepSystem(t5.BaseBiker.GetGameState())
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
