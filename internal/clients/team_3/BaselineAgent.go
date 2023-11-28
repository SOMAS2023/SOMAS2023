package team_3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"sort"
)

type IBaselineAgent interface {
	objects.IBaseBiker
}

type BaselineAgent struct {
	objects.BaseBiker
	currentBike   *objects.MegaBike
	targetLootBox objects.ILootBox
}

// DecideAction only pedal
func (agent *BaselineAgent) DecideAction() objects.BikerAction {
	return objects.Pedal
}

// DecideForces randomly based on current energyLevel
func (agent *BaselineAgent) DecideForces(direction uuid.UUID) {
	energyLevel := agent.GetEnergyLevel() // 当前能量

	randomPedalForce := rand.Float64() * energyLevel // 使用 rand 包生成随机的 pedal 力量，可以根据需要调整范围

	// 因为force是一个struct,包括pedal, brake,和turning，因此需要一起定义，不能够只有pedal
	forces := utils.Forces{
		Pedal: randomPedalForce,
		Brake: 0.0, // 这里默认刹车为 0
		Turning: utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: physics.ComputeOrientation(agent.GetLocation(), agent.GetGameState().GetMegaBikes()[direction].GetPosition()) - agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()].GetOrientation(),
		}, // 这里默认转向为 0
	}

	agent.SetForces(forces)
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
	e := agent.decideTargetLootBox(agent.GetGameState().GetLootBoxes())
	if e != nil {
		panic("unexpected error!")
	}
	return agent.targetLootBox.GetPosition()
}

func (agent *BaselineAgent) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	boxesInMap := agent.GetGameState().GetLootBoxes()
	boxProposed := make([]objects.ILootBox, len(proposals))
	for i, pp := range proposals {
		boxProposed[i] = boxesInMap[pp]
	}
	rank, e := agent.rankTargetProposals(boxProposed)
	if e != nil {
		panic("unexpected error!")
	}
	return rank
}

func (agent *BaselineAgent) DecideAllocation() voting.IdVoteMap {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()]
	rank, e := agent.rankAgentsReputation(currentBike.GetAgents())
	if e != nil {
		panic("unexpected error!")
	}
	return rank
}

// decideTargetLootBox find closest lootBox
func (agent *BaselineAgent) decideTargetLootBox(lootBoxes map[uuid.UUID]objects.ILootBox) error {

	agentLocation := agent.GetLocation() //agent location
	shortestDistance := math.MaxFloat64  //最短距离一开始设置为正无穷

	for _, lootbox := range lootBoxes { //遍历每一个lootbox
		lootboxLocation := lootbox.GetPosition()
		distance := physics.ComputeDistance(agentLocation, lootboxLocation)

		if distance < shortestDistance {
			shortestDistance = distance
			agent.targetLootBox = lootbox
		}
	}
	return nil
}

// rankTargetProposals rank by distance
func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()]
	// sort lootBox by distance
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(currentBike.GetPosition(), proposedLootBox[i].GetPosition()) < physics.ComputeDistance(currentBike.GetPosition(), proposedLootBox[j].GetPosition())
	})
	rank := make(map[uuid.UUID]float64)
	for i, lootBox := range proposedLootBox {
		rank[lootBox.GetID()] = float64(i)
	}
	return rank, nil
}

// rankAgentReputation randomly rank agents
func (agent *BaselineAgent) rankAgentsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	rank := make(map[uuid.UUID]float64)
	for i, agent := range agentsOnBike {
		rank[agent.GetID()] = float64(i)
	}
	return rank, nil
}
