package team_3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"github.com/google/uuid"
	"math"
	"sort"
)

type ISmartAgent interface {
	objects.IBaseBiker
}

type SmartAgent struct {
	objects.BaseBiker
	currentBike   *objects.MegaBike
	targetLootBox objects.ILootBox
	reputationMap map[uuid.UUID]reputation

	lootBoxCnt                     float64
	energySpent                    float64
	lastEnergyLevel                float64
	satisfactionOfRecentAllocation float64
}

// DecideAction only pedal
func (agent *SmartAgent) DecideAction() objects.BikerAction {
	if agent.GetEnergyLevel() < agent.lastEnergyLevel {
		agent.energySpent += agent.lastEnergyLevel - agent.GetEnergyLevel()
	} else {
		agent.recalculateSatisfaction()
	}
	agent.lastEnergyLevel = agent.GetEnergyLevel()

	agent.updateRepMap()

	return objects.Pedal
}

// DecideForces considering Hegselmann-Krause model, Ramirez-Cano-Pitt model and Satisfaction
func (agent *SmartAgent) DecideForces(direction uuid.UUID) {
	agentsOnBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()].GetAgents()
	scores := make(map[uuid.UUID]float64)
	totalScore := 0.0
	for _, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		score := rep.isSameColor/ // Cognitive dimension: is same belief?
			+rep.historyContribution + rep.lootBoxGet/ // Pareto principle: give more energy to those with more outcome
			+rep.recentContribution // Forgiveness: forgive agents pedal harder recently
		scores[others.GetID()] = score
		totalScore += score
	}

	for id, score := range scores {
		scores[id] = score / totalScore
	}

	pedalForce := 0.0
	for id, weight := range scores {
		pedalForce += weight * agent.reputationMap[id].lastPedal
	}
	pedalForce /= agent.satisfactionOfRecentAllocation

	// 因为force是一个struct,包括pedal, brake,和turning，因此需要一起定义，不能够只有pedal
	forces := utils.Forces{
		Pedal: pedalForce,
		Brake: 0.0, // 这里默认刹车为 0
		Turning: utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: physics.ComputeOrientation(agent.GetLocation(), agent.GetGameState().GetMegaBikes()[direction].GetPosition()) - agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()].GetOrientation(),
		}, // 这里默认转向为 0
	}

	agent.SetForces(forces)
}

// DecideJoining accept all
func (agent *SmartAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

func (agent *SmartAgent) ProposeDirection() utils.Coordinates {
	e := agent.decideTargetLootBox(agent.GetGameState().GetLootBoxes())
	if e != nil {
		panic("unexpected error!")
	}
	return agent.targetLootBox.GetPosition()
}

func (agent *SmartAgent) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	boxesInMap := agent.GetGameState().GetLootBoxes()
	boxProposed := make([]objects.ILootBox, len(proposals))
	for i, pp := range proposals {
		boxProposed[i] = boxesInMap[pp]
	}
	rank, _ := agent.rankTargetProposals(boxProposed)
	return rank
}

func (agent *SmartAgent) DecideAllocation() voting.IdVoteMap {
	agent.lootBoxCnt += 1
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()]
	vote, _ := agent.scoreAgentsForAllocation(currentBike.GetAgents())
	return vote
}

// decideTargetLootBox find closest lootBox
func (agent *SmartAgent) decideTargetLootBox(lootBoxes map[uuid.UUID]objects.ILootBox) error {

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
func (agent *SmartAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error) {
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

// scoreAgentsForAllocation if self energy level is low (below average cost for a lootBox), we follow 'Smallest First', else 'Ration'
func (agent *SmartAgent) scoreAgentsForAllocation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	scores := make(map[uuid.UUID]float64)
	totalScore := 0.0
	if agent.energySpent/agent.lootBoxCnt > agent.GetEnergyLevel() {
		// go 'Smallest First' strategy, only take energyRemain into consideration
		for _, others := range agentsOnBike {
			id := others.GetID()
			score := agent.reputationMap[id].energyRemain
			scores[others.GetID()] = score
			totalScore += score
		}
	} else {
		// go 'Ration' strategy, considering all facts
		for _, others := range agentsOnBike {
			id := others.GetID()
			rep := agent.reputationMap[id]
			score := rep.isSameColor/ // Cognitive dimension: is same belief?
				+rep.historyContribution + rep.lootBoxGet/ // Pareto principle: give more energy to those with more outcome
				+rep.recentContribution/ // Forgiveness: forgive agents pedal harder recently
				-rep.energyGain/ // Equality: Agents received more energy before should get less this time
				+rep.energyRemain // Need: Agents with lower energyLevel require more, try to meet their need
			scores[others.GetID()] = score
			totalScore += score
		}
	}

	// normalize scores
	for id, score := range scores {
		scores[id] = score / totalScore
	}

	return scores, nil
}

func (agent *SmartAgent) updateRepMap() {
	if agent.reputationMap == nil {
		agent.reputationMap = make(map[uuid.UUID]reputation)
	}
	for _, bikes := range agent.GetGameState().GetMegaBikes() {
		for _, otherAgent := range bikes.GetAgents() {
			rep, exist := agent.reputationMap[otherAgent.GetID()]
			if !exist {
				rep = reputation{}
			}
			rep.updateScore(otherAgent, agent.GetColour())
			agent.reputationMap[otherAgent.GetID()] = rep
		}
	}
}

func (agent *SmartAgent) recalculateSatisfaction() {
	scores := []float64{}

	totalScore := 0.0
	agentsOnBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()].GetAgents()
	for _, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		score := rep.isSameColor/ // Cognitive dimension: is same belief?
			+rep.historyContribution + rep.lootBoxGet/ // Pareto principle: give more energy to those with more outcome
			+rep.recentContribution/ // Forgiveness: forgive agents pedal harder recently
			-rep.energyGain/ // Equality: Agents received more energy before should get less this time
			+rep.energyRemain // Need: Agents with lower energyLevel require more, try to meet their need
		scores[others.GetID()] = score
		totalScore += score
	}
}

func measureOrder(input []float64) float64 {
	inversionCnt := 0.0
	size := len(input)
	for i, n := range input {
		j := i + 1
		for j < size {
			if n > input[j] { // 升序为正序
				inversionCnt += 1
			}
			j += 1
		}
	}
	return 1.0 - 2.0*inversionCnt/float64(size*(size-1))
}
