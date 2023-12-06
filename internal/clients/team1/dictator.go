// DICTATOR FUNCTIONS

package team1

import (
	voting "SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

//--------------------DICTATOR FUNCTIONS------------------

// ** called only when the agent is the dictator
func (bb *Biker1) DictateDirection() uuid.UUID {
	// TODO: make more sophisticated
	tmp, _ := bb.nearestLootColour()
	return tmp
}

// ** decide which agents to kick out (dictator)
func (bb *Biker1) DecideKickOut() []uuid.UUID {

	// TODO: make more sophisticated
	tmp := []uuid.UUID{}
	agent := bb.lowestOpinionKick()
	if agent != uuid.Nil {
		tmp = append(tmp, agent)
	}
	//tmp = append(tmp, bb.lowestOpinionKick())
	return tmp
}

// ** decide the allocation (dictator)
func (bb *Biker1) DecideDictatorAllocation() voting.IdVoteMap {
	return bb.DecideAllocation()
}

//--------------------END OF DICTATOR FUNCTIONS------------------
