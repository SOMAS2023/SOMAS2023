package benchmark

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

type BenchmarkAgent struct {
	objects.BaseBiker
	currentBike     *objects.MegaBike
	PreferredColour utils.Colour
}

// DecideAction
func (agent *BenchmarkAgent) decideAction() objects.BikerAction {
	return objects.Pedal
}

// finding the closest same color lootBox
func (agent *BenchmarkAgent) decideTargetLootBox(lootBoxes map[uuid.UUID]objects.ILootBox) (objects.ILootBox, error) {
	agentLocation := agent.GetLocation()
	shortestDistance := math.MaxFloat64
	var nearestLootbox objects.ILootBox

	for _, lootbox := range lootBoxes {
		lootboxLocation := lootbox.GetPosition()
		distance := physics.ComputeDistance(agentLocation, lootboxLocation)
		if distance < shortestDistance {
			shortestDistance = distance
			nearestLootbox = lootbox
		}
	}
	return nearestLootbox, nil
}

func (agent *BenchmarkAgent) decideForces() {
	forces := utils.Forces{
		Pedal:   1.0,
		Brake:   0.0,
		Turning: 0.0,
	}

	println("forces for each round", forces)
}

func computeTurningAngle(agentPosition, lootBoxPosition utils.Coordinates) float64 {
	deltaX := lootBoxPosition.X - agentPosition.X
	deltaY := lootBoxPosition.Y - agentPosition.Y
	angle := math.Atan2(deltaY, deltaX) / math.Pi

	return angle * 2 / math.Pi
}

func (agent *BenchmarkAgent) voteOnJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agentID := range pendingAgents {
		decision[agentID] = true
	}
	return decision
}

func (b *BenchmarkAgent) VoteOnLootBox(lootBoxOptions []uuid.UUID) map[uuid.UUID]float64 {
	voteDistribution := make(map[uuid.UUID]float64)
	numOptions := len(lootBoxOptions)

	// If there are loot box options, distribute votes equally
	if numOptions > 0 {
		equalVote := 1.0 / float64(numOptions)
		for _, lootBoxID := range lootBoxOptions {
			voteDistribution[lootBoxID] = equalVote
		}
	}

	return voteDistribution
}

func (agent *BenchmarkAgent) voteOnTargetProposals(proposedLootBox []objects.LootBox) (map[utils.Coordinates]float64, error) {
	rank := make(map[utils.Coordinates]float64)
	equalRank := 1.0 // Assigning the same rank to all

	for _, lootBox := range proposedLootBox {
		rank[lootBox.GetPosition()] = equalRank
	}

	return rank, nil
}

func (agent *BenchmarkAgent) rankAgentsReputation(agentsOnBike []objects.BaseBiker) (map[uuid.UUID]float64, error) {
	rank := make(map[uuid.UUID]float64)
	equalRank := 1.0 // Assigning the same rank to all agents

	for _, agent := range agentsOnBike {
		rank[agent.GetID()] = equalRank
	}

	return rank, nil
}

func (b *BenchmarkAgent) ChooseBike(gameState objects.IGameState) uuid.UUID {
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

func (b *BenchmarkAgent) DecideBikeChangeBasedOnVote(votingResult uuid.UUID, gameState objects.IGameState) uuid.UUID {
	winningLootBox := gameState.GetLootBoxes()[votingResult]

	if winningLootBox.GetColour() != b.PreferredColour {
		return b.ChooseBike(gameState)
	}

	return uuid.Nil
}
