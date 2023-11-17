package benchmark

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

type BenchmarkAgent struct {
	// Embed BaseBiker or any relevant struct
	*objects.PhysicsObject
	PreferredColour utils.Colour
	CurrentBikeID   uuid.UUID
	OnBike          bool
}

func NewBenchmarkAgent(preferredColour utils.Colour) *BenchmarkAgent {
	physicsObject := objects.GetPhysicsObject(utils.MassBiker)
	return &BenchmarkAgent{
		PhysicsObject:   physicsObject,
		PreferredColour: preferredColour,
		CurrentBikeID:   uuid.Nil,
		OnBike:          true,
	}
}

// choose the target bike
func (b *BenchmarkAgent) ChooseBike(gameState objects.IGameState) uuid.UUID {
	closestBikeID := uuid.Nil
	minDistance := math.MaxFloat64
	targetLootBox := gameState.GetLootBoxes()[b.ChooseLootBox(gameState)]

	for bikeID, bike := range gameState.GetMegaBikes() {
		distance := physics.ComputeDistance(bike.GetPosition(), targetLootBox.GetPosition())
		if distance < minDistance {
			minDistance = distance
			closestBikeID = bikeID
		}
	}
	return closestBikeID
}

// communication between the server and the agent, if have the leaving signal then leave
func (b *BenchmarkAgent) VoteOnJoiningRequests(joiningRequests []uuid.UUID, leavingSignal bool, gameState objects.IGameState) map[uuid.UUID]float64 {
	if leavingSignal {
		//leav the current bike
		targetBikeID := b.ChooseBike(gameState)
		b.CurrentBikeID = targetBikeID
		b.OnBike = false

		// communication

		return nil // no voting necessary
	}

	// If not leaving, vote on joining requests
	votes := make(map[uuid.UUID]float64)
	for _, requestID := range joiningRequests {
		votes[requestID] = 1.0 / float64(len(joiningRequests)) // Fair voting
	}
	return votes
}

// choose a target lootbox
func (b *BenchmarkAgent) ChooseLootBox(gameState objects.IGameState) uuid.UUID {
	closestLootBoxID := uuid.Nil
	minDistance := math.MaxFloat64
	agentPosition := b.GetPosition()

	for _, lootBox := range gameState.GetLootBoxes() {
		if lootBox.GetColour() == b.PreferredColour {
			lootBoxPosition := lootBox.GetPosition()
			distance := math.Sqrt(math.Pow(lootBoxPosition.X-agentPosition.X, 2) + math.Pow(lootBoxPosition.Y-agentPosition.Y, 2))
			if distance < minDistance {
				minDistance = distance
				closestLootBoxID = lootBox.GetID()
			}
		}
	}
	return closestLootBoxID
}

func (b *BenchmarkAgent) GetAcceleration(force float64) float64 {
	return force / b.GetMass()
}

func (b *BenchmarkAgent) GetPosition() utils.Coordinates {
	return b.PhysicsObject.GetPosition()
}

// vote equally on ranking, vote for the lootbox choosing ranking
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

// decide to choose pedal or change the bike, return teh action and the calues
func (b *BenchmarkAgent) DecideActions(action string, gameState objects.IGameState) map[string]float64 {
	actions := make(map[string]float64)

	if action == "Pedal" {
		actions["Pedal"] = 1.0
		actions["Brake"] = 0.0
		actions["Turning"] = 0.0
	} else if action == "ChangeBike" {
		//choose the closest lootbox of the same color as the future target
		closestLootBoxID := b.ChooseLootBox(gameState)
		closestLootBox := gameState.GetLootBoxes()[closestLootBoxID]
		angle := computeSteeringAngle(b.GetPosition(), closestLootBox.GetPosition())

		actions["Brake"] = 0.0
		actions["Turning"] = normalizeAngleToRange(angle)
	}

	return actions
}

// get the angle aof a closest lootbox of the same color
func computeSteeringAngle(currentPosition, targetPosition utils.Coordinates) float64 {
	// Compute the angle from the current position to the target position
	xDiff := targetPosition.X - currentPosition.X
	yDiff := targetPosition.Y - currentPosition.Y
	return math.Atan2(yDiff, xDiff) / math.Pi
}

// steering angle
func normalizeAngleToRange(angle float64) float64 {
	// Normalize the angle to the range of [-1, 1]
	normalizedAngle := angle / math.Pi
	if normalizedAngle > 1 {
		return 1
	} else if normalizedAngle < -1 {
		return -1
	}
	return normalizedAngle
}

// vote resource allocation, still equally distributed
func (b *BenchmarkAgent) VoteResourceAllocation(onBikeAgents []uuid.UUID) map[uuid.UUID]float64 {
	allocation := make(map[uuid.UUID]float64)
	for _, agentID := range onBikeAgents {
		allocation[agentID] = 1.0 / float64(len(onBikeAgents)) // Fair allocation
	}
	return allocation
}
