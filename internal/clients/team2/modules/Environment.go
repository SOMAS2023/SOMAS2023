package modules

import (
	objects "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"fmt"
	"math"
	"math/rand"

	"github.com/google/uuid"
)

const (
	AudiRange = 10
)

type EnvironmentModule struct {
	AgentId   uuid.UUID
	GameState objects.IGameState
	BikeId    uuid.UUID
}

///
/// GameState
///

func (e *EnvironmentModule) SetGameState(gameState objects.IGameState) {
	e.GameState = gameState
}

///
/// Lootboxes
///

func (e *EnvironmentModule) GetLootBoxes() map[uuid.UUID]objects.ILootBox {
	return e.GameState.GetLootBoxes()
}

func (e *EnvironmentModule) GetLootBoxById(lootboxId uuid.UUID) objects.ILootBox {
	return e.GetLootBoxes()[lootboxId]
}

func (e *EnvironmentModule) GetLootboxPos(lootboxId uuid.UUID) utils.Coordinates {
	return e.GetLootBoxById(lootboxId).GetPosition()
}

func (e *EnvironmentModule) GetLootBoxesByColor(color utils.Colour) map[uuid.UUID]objects.ILootBox {
	lootboxes := e.GetLootBoxes()
	lootboxesFiltered := make(map[uuid.UUID]objects.ILootBox)
	for _, lootbox := range lootboxes {
		if lootbox.GetColour() == color {
			lootboxesFiltered[lootbox.GetID()] = lootbox
		}
	}
	return lootboxesFiltered
}

func (e *EnvironmentModule) GetNearestLootbox(agentId uuid.UUID) uuid.UUID {
	nearestLootbox := uuid.Nil
	minDist := math.MaxFloat64
	for _, lootbox := range e.GetLootBoxes() {
		dist := e.GetDistanceToLootbox(lootbox.GetID())
		if dist < minDist {
			minDist = dist
			nearestLootbox = lootbox.GetID()
		}
	}
	return nearestLootbox
}

func (e *EnvironmentModule) GetNearestLootboxByColor(agentId uuid.UUID, color utils.Colour) uuid.UUID {
	nearestLootbox := uuid.Nil
	minDist := math.MaxFloat64
	for _, lootbox := range e.GetLootBoxesByColor(color) {
		dist := e.GetDistanceToLootbox(lootbox.GetID())
		if dist < minDist {
			minDist = dist
			nearestLootbox = lootbox.GetID()
		}
	}
	if nearestLootbox == uuid.Nil {
		return e.GetNearestLootbox(e.AgentId)
	}
	return nearestLootbox
}

func (e *EnvironmentModule) GetDistanceToLootbox(lootboxId uuid.UUID) float64 {
	bikePos, agntPos := e.GetBikeById(e.BikeId).GetPosition(), e.GetLootBoxById(lootboxId).GetPosition()

	return e.GetDistance(bikePos, agntPos)
}

// Gets lootbox with the highest gain.
// We define gain as the distance to the lootbox divided by the total resources in the lootbox.
func (e *EnvironmentModule) GetHighestGainLootbox() uuid.UUID {
	bestGain := float64(0)
	bestLoot := uuid.Nil
	for _, lootboxId := range e.GetLootBoxes() {

		gain := lootboxId.GetTotalResources() / e.GetDistanceToLootbox(lootboxId.GetID())
		if gain > bestGain {
			bestGain = gain
			bestLoot = lootboxId.GetID()
		}
	}
	return bestLoot
}

func (e *EnvironmentModule) GetNearestLootboxAwayFromAudi() uuid.UUID {
	// Find positions.
	bikePos := e.GetBikeById(e.BikeId).GetPosition()
	audiPos := e.GetAudi().GetPosition()

	// Find position away from audi.
	deltaX := audiPos.X - bikePos.X
	deltaY := audiPos.Y - bikePos.Y

	awayX := bikePos.X - deltaX
	awayY := bikePos.Y - deltaY
	awayPos := utils.Coordinates{X: awayX, Y: awayY}

	// Find nearest lootbox away from audi.
	minLoot := uuid.Nil
	minDist := math.MaxFloat64
	for id, lootbox := range e.GetLootBoxes() {
		dist := e.GetDistance(awayPos, lootbox.GetPosition())
		if dist < minDist {
			minDist = dist
			minLoot = id
		}
	}
	return minLoot
}

///
/// Bikes
///

func (e *EnvironmentModule) GetAudi() objects.IAudi {
	return e.GameState.GetAudi()
}

func (e *EnvironmentModule) GetBikes() map[uuid.UUID]objects.IMegaBike {
	return e.GameState.GetMegaBikes()
}

func (e *EnvironmentModule) GetBikeById(bikeId uuid.UUID) objects.IMegaBike {
	return e.GetBikes()[bikeId]
}

func (e *EnvironmentModule) GetBike() objects.IMegaBike {
	return e.GetBikeById(e.BikeId)
}

func (e *EnvironmentModule) GetBikeOrientation() float64 {
	return e.GetBikeById(e.BikeId).GetOrientation()
}

func (e *EnvironmentModule) GetBikerWithMaxSocialCapital(sc *SocialCapital) (uuid.UUID, float64) {
	fellowBikers := e.GetBikerAgents()
	maxSCAgentId := uuid.Nil
	maxSC := 0.0
	for _, fellowBiker := range fellowBikers {
		if sc, ok := sc.SocialCapital[e.AgentId]; ok {
			if sc >= maxSC {
				maxSCAgentId = fellowBiker.GetID()
				maxSC = sc
			}
		}
	}
	return maxSCAgentId, maxSC
}

func (e *EnvironmentModule) GetBikerWithMinSocialCapital(sc *SocialCapital) (uuid.UUID, float64) {
	fellowBikers := e.GetBikerAgents()
	minSCAgentId := uuid.Nil
	minSC := math.MaxFloat64
	for _, fellowBiker := range fellowBikers {
		if sc, ok := sc.SocialCapital[e.AgentId]; ok {
			if sc < minSC {
				minSCAgentId = fellowBiker.GetID()
				minSC = sc
			}
		}
	}
	return minSCAgentId, minSC
}

func (e *EnvironmentModule) GetBikeWithMaximumSocialCapital(sc *SocialCapital) uuid.UUID {
	maxAverage := float64(0)
	maxBikeId := uuid.Nil

	bikes := e.GetBikes()
	for bikeId, bike := range bikes {
		totalSocialCapital := float64(0)
		agentCount := float64(len(bike.GetAgents()))

		// Sum up the social capital of all agents on this bike
		for _, agent := range bike.GetAgents() {
			agentId := agent.GetID()
			totalSocialCapital += sc.SocialCapital[agentId]
		}

		// Calculate average social capital for this bike, Assume we don't swtich to a bike with 0 agents
		if agentCount > 0 {
			averageSocialCapital := totalSocialCapital / agentCount
			if averageSocialCapital > maxAverage {
				maxAverage = averageSocialCapital
				maxBikeId = bikeId
			}
		}
	}

	if maxBikeId != uuid.Nil || maxBikeId == e.BikeId {
		// If found, change to that bike.
		return maxBikeId
	} else {
		// Otherwise, change to a random bike.
		i, targetI := 0, rand.Intn(len(bikes))
		for id := range bikes {
			if i == targetI {
				return id
			}
			i++
		}
		panic("No bikes found to change to.")
	}
}

func (e *EnvironmentModule) GetDistanceToAudi() float64 {
	bikePos, audiPos := e.GetBikeById(e.BikeId).GetPosition(), e.GetAudi().GetPosition()

	fmt.Printf("[GetDistanceToAudi] Pos of bike: %f\n", bikePos)
	fmt.Printf("[GetDistanceToAudi] Pos of Audi: %f\n", audiPos)

	return e.GetDistance(bikePos, audiPos)
}

func (e *EnvironmentModule) IsAudiNear() bool {
	fmt.Printf("[IsAudiNear] Distance to audi: %f\n", e.GetDistanceToAudi())
	return e.GetDistanceToAudi() <= AudiRange
}

func (e *EnvironmentModule) GetBikerAgents() map[uuid.UUID]objects.IBaseBiker {
	bikes := e.GetBikes()
	bikerAgents := make(map[uuid.UUID]objects.IBaseBiker)
	for _, bike := range bikes {
		for _, biker := range bike.GetAgents() {
			bikerAgents[biker.GetID()] = biker
		}
	}
	return bikerAgents
}

///
/// Utils
///

func (e *EnvironmentModule) GetDistance(pos1, pos2 utils.Coordinates) float64 {

	return math.Sqrt(math.Pow(pos1.X-pos2.X, 2) + math.Pow(pos1.Y-pos2.Y, 2))
}

func GetEnvironmentModule(agentId uuid.UUID, gameState objects.IGameState, bikeId uuid.UUID) *EnvironmentModule {
	return &EnvironmentModule{
		AgentId:   agentId,
		GameState: gameState,
		BikeId:    bikeId,
	}
}
