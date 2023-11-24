package agents

import (
	objects "SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

type EnvironmentHandler struct {
	GameState     objects.IGameState // Game state to be updated in each round
	CurrentBikeId uuid.UUID          // unique ID of current bike
	AgentId       uuid.UUID          // unique ID of agent
}

func (env *EnvironmentHandler) GetLootBoxesByColour(colour utils.Colour) []objects.ILootBox {
	lootBoxes := env.GameState.GetLootBoxes()
	var matchingLootBoxes []objects.ILootBox
	for _, lootBox := range lootBoxes {
		if lootBox.GetColour() == colour {
			matchingLootBoxes = append(matchingLootBoxes, lootBox)
		}
	}
	return matchingLootBoxes
}

func (env *EnvironmentHandler) GetCurrentBike() objects.IMegaBike {
	return env.GameState.GetMegaBikes()[env.CurrentBikeId]
}

func (env *EnvironmentHandler) GetAgentsOnCurrentBike() []objects.IBaseBiker {
	return env.GetBikeAgentsByBikeId(env.CurrentBikeId)
}

func (env *EnvironmentHandler) GetBikeAgentsByBikeId(bikeId uuid.UUID) []objects.IBaseBiker {
	megaBikes := env.GameState.GetMegaBikes()
	bike := megaBikes[bikeId]
	return bike.GetAgents()
}

func (env *EnvironmentHandler) GetNearestLootBox() objects.ILootBox {
	X, Y := env.GetCurrentBike().GetPosition().X, env.GetCurrentBike().GetPosition().Y
	lootBoxes := env.GameState.GetLootBoxes()
	var nearestLootBox objects.ILootBox
	var nearestDistance float64
	for _, lootBox := range lootBoxes {
		x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
		distance := math.Sqrt(math.Pow(X-x, 2) + math.Pow(Y-y, 2))
		if nearestLootBox == nil || distance < nearestDistance {
			nearestLootBox = lootBox
			nearestDistance = distance
		}
	}

	return nearestLootBox
}

func NewEnvironmentHandler(gameState objects.IGameState, bikeId uuid.UUID, agentId uuid.UUID) *EnvironmentHandler {
	return &EnvironmentHandler{
		GameState:     gameState,
		CurrentBikeId: bikeId,
		AgentId:       agentId,
	}
}
