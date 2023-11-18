package team_3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"math"
	"math/rand"
	"github.com/google/uuid"
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
	energyLevel := agent.GetEnergyLevel() // 当前能量

	randomPedalForce := rand.Float64() * energyLevel // 使用 rand 包生成随机的 pedal 力量，可以根据需要调整范围

	// 因为force是一个struct,包括pedal, brake,和turning，因此需要一起定义，不能够只有pedal
	forces := utils.Forces{
		Pedal:   randomPedalForce,
		Brake:   0.0, // 这里默认刹车为 0
		Turning: 0.0, // 这里默认转向为 0
	}

	// 将决定的力量设置给 BaseBiker 对象, 未做到,缺少函数
	println("forces for each round", forces)
	// 缺少给agent赋forces的函数,目前只有GetForces,只能从函数获取,没办法给函数赋新值
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

	agentLocation := agent.GetLocation() //agent location
	var nearestLootbox objects.ILootBox  //最近的一个lootbox
	shortestDistance := math.MaxFloat64  //最短距离一开始设置为正无穷

	for _, lootbox := range lootBoxes { //遍历每一个lootbox
		lootboxLocation := lootbox.GetPosition()
		distance := physics.ComputeDistance(agentLocation, lootboxLocation)

		if distance < shortestDistance {
			shortestDistance = distance
			nearestLootbox = lootbox
		}
	}
	return nearestLootbox, nil
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
