package objects

import (
	"github.com/google/uuid"
)

/*
IGameState is an interface for GameState that objects will use to get the current game state
*/
type IGameState interface {
	GetLootBoxes() map[uuid.UUID]*LootBox
	GetMegaBikes() map[uuid.UUID]*MegaBike
	// GetMegaBikeRiders returns a mapping from Agent ID -> ID of the bike that they are riding
	GetMegaBikeRiders() map[uuid.UUID]uuid.UUID
}
