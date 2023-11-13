package server

import (
	obj "SOMAS2023/internal/common/objects"

	"github.com/google/uuid"
)

/*
The GameState is an implementation of the IGameState interface
*/
type GameState struct {
	BaseBikers map[uuid.UUID]obj.IBaseBiker // Map of IDs to IBaseBiker objects
	LootBoxes  map[uuid.UUID]obj.LootBox    // Map of IDs to LootBox objects
	MegaBikes  map[uuid.UUID]obj.MegaBike   // Map of IDs to LootBox objects
}

func (g GameState) GetGameState() obj.IGameState {
	//TODO implement me
	panic("implement me")
}

// The usage of GameState would be as follows:

// // Add a Biker to the GameState
// bikerID := "biker1"
// gameState.BaseBikers[bikerID] = /* Biker instance that implements IBaseBiker */

// // Add a LootBox to the GameState
// lootBoxID := "lootbox1"
// gameState.LootBoxes[lootBoxID] = obj.LootBox{ /* initializer fields */ }

// // Retrieve a Biker from the GameState
// biker, bikerExists := gameState.BaseBikers[bikerID]
// if bikerExists {
// 	// Use the biker
// }

// // Remove a Biker from the GameState
// delete(gameState.BaseBikers, bikerID)
