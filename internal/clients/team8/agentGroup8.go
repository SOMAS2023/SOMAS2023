package team_8

import (
	"SOMAS2023/internal/common/objects"
	voting "SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

type Agent8 struct {
	*objects.BaseBiker
}

// in the MVP the biker's action defaults to pedaling (as it won't be able to change bikes)
// in future implementations this function will be overridden by the agent's specific strategy
// which will be used to determine whether to pedalor try to change bike
func (bb *Agent8) DecideAction() objects.BikerAction {
	//TODO： need to be implemented
	return 0
}

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	//TODO： need to be implemented
	return
}

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *Agent8) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	//TODO： need to be implemented
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

// decide which bike to go to
// for now it just returns a random uuid
func (bb *Agent8) ChangeBike() uuid.UUID {
	//TODO： need to be implemented
	return uuid.New()
}

// default implementation returns the id of the nearest lootbox
func (bb *Agent8) ProposeDirection() uuid.UUID {
	//TODO： need to be implemented
	return uuid.New()
}

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *Agent8) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	//TODO： need to be implemented
	votes := make(map[uuid.UUID]float64)
	return votes
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	//TODO： need to be implemented
	distribution := make(map[uuid.UUID]float64)
	return distribution
}
