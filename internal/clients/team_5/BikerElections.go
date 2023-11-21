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
	agentReputations := calculateMegaBikeReputation(bikeId) // Need to ask SG how it works after merging

	for key, value := range agentReputations {
		fmt.Println(value)

		if value > threshold {
			decisions[key] = true
		} else {
			decisions[key] = false
		}

	}

	return decisions

}
