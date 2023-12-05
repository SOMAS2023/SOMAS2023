package objects

import "github.com/google/uuid"

/*
IGameState is an interface for GameState that objects will use to get the current game state
*/
type IGameState interface {
	GetLootBoxes() map[uuid.UUID]ILootBox
	GetMegaBikes() map[uuid.UUID]IMegaBike
	GetAgents() map[uuid.UUID]IBaseBiker
	GetAudi() IAudi
}
