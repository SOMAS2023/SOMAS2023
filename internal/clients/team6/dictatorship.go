package team6

import (
	voting "SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

func (bb *Team6Biker) DictateDirection() uuid.UUID {
	return bb.ProposeDirection()
}

// func (bb *Team6Biker) DecideKickOut() []uuid.UUID {
// 	return (make([]uuid.UUID, 0))
// }

func (bb *Team6Biker) DecideDictatorAllocation() voting.IdVoteMap {
	/*fellowBikers := bb.GetFellowBikers()
	distribution := make(voting.IdVoteMap)
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		distribution[agentID] = float64(1 / len(fellowBikers))
	}
	return distribution*/
	bikeID := bb.GetBike()
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	equalDist := 1.0 / float64(len(fellowBikers))
	for _, agent := range fellowBikers {
		distribution[agent.GetID()] = equalDist
	}
	return distribution
}
