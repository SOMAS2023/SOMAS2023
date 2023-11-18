package team_3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
)

type BaselineAgent struct {
	objects.BaseBiker
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
	panic("to be implemented!")
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
func (agent *BaselineAgent) decideTargetLootBox(lootBoxes map[uuid.UUID]objects.ILootBox) (objects.LootBox, error) {
	panic("to be implemented!")
}

// rankTargetProposals rank by distance
func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.LootBox) (map[utils.Coordinates]float64, error) {
	panic("to be implemented!")
}

// rankAgentReputation randomly rank agents
func (agent *BaselineAgent) rankAgentsReputation([]objects.BaseBiker) (map[uuid.UUID]float64, error) {
	panic("to be implemented!")
}
