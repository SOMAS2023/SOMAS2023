package team_4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"math"

	"sort"

	"github.com/google/uuid"
)

type team4Agent struct {
	objects.BaseBiker
	currentBike     *objects.MegaBike
	PreferredColour utils.Colour
	GameState       objects.IGameState
}

// DecideAction
func (agent *team4Agent) decideAction() objects.BikerAction {
	return objects.Pedal
}

// finding the closest same color lootBox
func (agent *team4Agent) decideTargetLootBox(lootBoxes map[uuid.UUID]objects.ILootBox) (objects.ILootBox, error) {
	agentLocation := agent.GetLocation()
	shortestDistance := math.MaxFloat64
	var nearestLootbox objects.ILootBox

	for _, lootbox := range lootBoxes {
		lootboxLocation := lootbox.GetPosition()
		distance := physics.ComputeDistance(agentLocation, lootboxLocation)
		if distance <= shortestDistance && lootbox.GetColour() == agent.PreferredColour {
			shortestDistance = distance
			nearestLootbox = lootbox
		}
	}
	return nearestLootbox, nil
}

func (agent *team4Agent) decideForces() {
	var distance float64
	energyLevel := agent.GetEnergyLevel()

	//decide the pedal force based on our logic
	currentPedalForce := energyLevel / distance
	agentPosition := agent.GetLocation()

	// Find the nearest lootbox of the same color as the agent
	nearestLootBox, err := agent.decideTargetLootBox(agent.GameState.GetLootBoxes())
	if err != nil || nearestLootBox == nil {
		panic("unexpected error!")
	}

	lootBoxPosition := nearestLootBox.GetPosition()

	distance = physics.ComputeDistance(agentPosition, lootBoxPosition)

	forces := utils.Forces{
		Pedal: currentPedalForce,
		Brake: 0.0, // 这里默认刹车为 0
		//Turning: 0.0, // 这里默认转向为 0 // 不知道还有没有这一部分strreing的问题
	}

	// Calculate the distance to the nearest lootbox

	// Set the turning angle based on the distance
	//var turningAngle float64
	//不知道能不能或的距离，到时候再看能不拿到决定的

	//if distance > 10 {
	//	turningAngle = computeTurningAngle(agentPosition, lootBoxPosition)
	//} else {
	//	turningAngle = 0.0
	//}

	//agent.SetForces(forces)
	//没找到在哪写，估计还没这个函数呢
	println("forces for each round", forces)
}

func computeTurningAngle(agentPosition, lootBoxPosition utils.Coordinates) float64 {
	deltaX := lootBoxPosition.X - agentPosition.X
	deltaY := lootBoxPosition.Y - agentPosition.Y
	angle := math.Atan2(deltaY, deltaX) / math.Pi

	return angle * 2 / math.Pi
}

func (agent *team4Agent) voteOnJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agentID := range pendingAgents {
		decision[agentID] = true
	}
	return decision
}

func (b *team4Agent) VoteOnLootBox(lootBoxOptions []uuid.UUID) map[uuid.UUID]float64 {
	voteDistribution := make(map[uuid.UUID]float64)

	// First, find the closest loot box of the preferred color
	closestLootBox, err := b.decideTargetLootBox(b.GameState.GetLootBoxes())
	if err != nil {
		panic("unexpected error!")
	}

	// If a closest loot box is found, give it the highest vote
	for _, lootBoxID := range lootBoxOptions {
		if lootBoxID == closestLootBox.GetID() {
			voteDistribution[lootBoxID] = 1.0 // Full vote for the closest loot box
		} else {
			voteDistribution[lootBoxID] = 0.0 // No vote for other loot boxes
		}
	}

	return voteDistribution
}

func (agent *team4Agent) voteOnTargetProposals(proposedLootBox []objects.LootBox) (map[utils.Coordinates]float64, error) {
	rank := make(map[utils.Coordinates]float64)
	equalRank := 1.0 // Assigning the same rank to all

	for _, lootBox := range proposedLootBox {
		rank[lootBox.GetPosition()] = equalRank
	}

	return rank, nil
}

// 所以到底怎么排序，选领导还是选排序，还是要都写，麻了
func (agent *team4Agent) rankTargetProposals(proposedLootBox []objects.LootBox) (map[utils.Coordinates]float64, error) {
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

func (b *team4Agent) ChooseBike(gameState objects.IGameState) uuid.UUID {
	closestBikeID := uuid.Nil
	minDistance := math.MaxFloat64

	// Find the nearest loot box of the preferred color
	nearestLootBox, err := b.decideTargetLootBox(gameState.GetLootBoxes())
	if err != nil || nearestLootBox == nil {
		return closestBikeID // Return nil if no suitable loot box is found
	}

	lootBoxPosition := nearestLootBox.GetPosition()

	// Iterate through all MegaBikes and find the one closest to the loot box
	for bikeID, bike := range gameState.GetMegaBikes() {
		bikePosition := bike.GetPosition()
		distance := physics.ComputeDistance(bikePosition, lootBoxPosition)
		if distance < minDistance {
			minDistance = distance
			closestBikeID = bikeID
		}
	}

	return closestBikeID
}

func (b *team4Agent) DecideBikeChangeBasedOnVote(votingResult uuid.UUID, gameState objects.IGameState) uuid.UUID {
	winningLootBox := gameState.GetLootBoxes()[votingResult]

	if winningLootBox.GetColour() != b.PreferredColour {
		return b.ChooseBike(gameState)
	}

	return uuid.Nil
}

// VoteForResourceAllocation - Distributes vote for resource allocation among all agents on the bike
func (b *team4Agent) VoteForResourceAllocation(agentsOnBike []uuid.UUID) map[uuid.UUID]float64 {
	vote := make(map[uuid.UUID]float64)
	numAgents := len(agentsOnBike)

	if numAgents == 0 {
		return vote
	}

	equalVote := 1.0 / float64(numAgents)
	for _, agentID := range agentsOnBike {
		vote[agentID] = equalVote
	}

	return vote
}

// 就大命的reputation，完全不知道怎么算，排个序笑笑算了
func (agent *team4Agent) rankReputation(agentsOnBike []objects.BaseBiker) (map[uuid.UUID]float64, error) {
	rank := make(map[uuid.UUID]float64)
	for i, agent := range agentsOnBike {
		rank[agent.GetID()] = float64(i)
	}
	return rank, nil
}

/*
// An example method where you might want to calculate honesty
func (agent *team4Agent) SomeDecisionMethod() {
	actionData := matrixCalculation.ActionData{
		// Populate with relevant data for actionData
	}

	// Assuming you have an instance of HonestyMatrix available,
	// you could also consider having a single instance that's passed around or referenced
	honestyMatrix := matrixCalculation.NewHonestyMatrix()
	honestyMatrix.CalculateHonestyBasedOnActions(agent.GetID(), actionData)

	// Optionally, access updated honesty records if needed
	honestyRecords := honestyMatrix.GetRecords(agent.GetID())
}
*/
