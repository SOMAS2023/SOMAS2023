package team_3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
	"math"
	"sort"
)

type BaselineAgent struct {
	objects.BaseBiker
	currentBike *objects.MegaBike
}

// DecideAction only pedal
func (agent *BaselineAgent) DecideAction() objects.BikerAction {
	return objects.Pedal
}

// DecideForces randomly based on current energyLevel
func (agent *BaselineAgent) DecideForces() {
	panic("DecideForces() to be implemented!")
}

// DecideJoining accept all
func (agent *BaselineAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

func (agent *BaselineAgent) ProposeDirection() utils.Coordinates {
	targetLootBox, e := agent.decideTargetLootBox(agent.GameState.GetLootBoxes())
	if e != nil {
		panic("unexpected error!")
	}
	return targetLootBox.GetPosition()
}

//func (agent *BaselineAgent) FinalDirectionVote([]utils.Coordinates) utils.PositionVoteMap{
//	panic("to be implemented!")
//}
//
//func (agent *BaselineAgent) DecideAllocationParameters(){
//	panic("to be implemented!")
//}

// decideTargetLootBox find closest lootBox
func (agent *BaselineAgent) decideTargetLootBox(lootBoxes map[uuid.UUID]objects.ILootBox) (objects.ILootBox, error) {
	var targetLootBox objects.ILootBox
	minDistance := math.MaxFloat64
	for _, lootBox := range lootBoxes {
		dist := physics.ComputeDistance(agent.currentBike.GetPosition(), lootBox.GetPosition())
		if dist < minDistance {
			minDistance = dist
			targetLootBox = lootBox
		}
	}
	return targetLootBox, nil
}

// rankTargetProposals rank by distance
func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.LootBox) (map[utils.Coordinates]float64, error) {
	// sort lootBox by distance
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(agent.currentBike.GetPosition(), proposedLootBox[i].GetPosition()) < physics.ComputeDistance(agent.currentBike.GetPosition(), proposedLootBox[j].GetPosition())
	})
	rank := make(map[utils.Coordinates]float64)
	for i, lootBox := range proposedLootBox {
		rank[lootBox.GetPosition()] = float64(i)
	}
	return rank, nil
}

// rankAgentReputation randomly rank agents
func (agent *BaselineAgent) rankAgentsReputation(agentsOnBike []objects.BaseBiker) (map[uuid.UUID]float64, error) {
	rank := make(map[uuid.UUID]float64)
	for i, agent := range agentsOnBike {
		rank[agent.GetID()] = float64(i)
	}
	return rank, nil
}
