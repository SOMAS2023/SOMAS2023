package team6

import (
	"github.com/google/uuid"
)

func (bb *Team6Biker) VoteForKickout() map[uuid.UUID]int {
	voteResults := make(map[uuid.UUID]int)
	bikeID := bb.GetBike()

	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if agentID != bb.GetID() {
			voteResults[agentID] = 0
			if agent.GetColour() == bb.GetColour() || (agent.GetReputation()[agentID] > reputationThreshold) {
				voteResults[agentID] = 0
			}

		}
	}
	return voteResults
}
