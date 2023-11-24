package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"math/rand"

	"github.com/google/uuid"
)

type Agent8 struct {
	objects.BaseBiker
	currentBike *objects.MegaBike
}

func (agent *Agent8) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	trust := make(map[uuid.UUID]int)
	for _, agent := range pendingAgents {
		//trust
		if trust[agent] >= 5 {
			decision[agent] = true
		}
	}
	return decision
}

func (agent *Agent8) ChangeBike() {
	return
}

func (agent *Agent8) DecideAction() {

	return objects.ChangeBike
} //BikerAction

func (agent *Agent8) ProposeDirection() {
	//
	return
}

func (agent *Agent8) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	return LootboxVoteMap
}

func (agent *Agent8) DecideForces() {
	energyLevel := agent.GetEnergyLevel()

	randomPedalForce := rand.Float64() * energyLevel
	forces := utils.Forces{
		Pedal:   randomPedalForce,
		Brake:   0.0,
		Turning: 0.0,
	}
	println("forces for each round", forces)
}

func (agent *Agent8) DecideAllocation() voting.IdVoteMap {

}
