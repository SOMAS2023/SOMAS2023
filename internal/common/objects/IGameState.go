package objects

import "github.com/google/uuid"

/*
IGameState is an interface for GameState that objects will use to get the current game state
*/
type IGameState interface {
	GetLootBoxes() map[uuid.UUID]*LootBox
	GetMegaBikes() map[uuid.UUID]*MegaBike
}
