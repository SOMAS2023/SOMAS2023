package utils

import (
	"SOMAS2023/internal/common/objects"
	"github.com/google/uuid"
)

type Colour int

const (
	Red Colour = iota
	Green
	Blue
	Yellow
	Orange
	Purple
	Pink
	Brown
	Gray
	White
	NumOfColours // add a sentinel for counting the number of colours
)

type Forces struct {
	Pedal   float64 // Pedal is a force from 0-1 where 1 is 100% power
	Brake   float64 // Brake is a force from 0-1 opposing the direction of travel (bike cannot go backwards)
	Turning float64 // Turning is a force from -1 to 1 which maps to -180° to 180°
}

type Coordinates struct {
	X float64
	Y float64
}

/*
IGameState is an interface for GameState that objects will use to get the current game state
*/
type IGameState interface {
	GetLootBoxes() map[uuid.UUID]*objects.LootBox
	GetMegaBikes() map[uuid.UUID]*objects.MegaBike
	// GetMegaBikeRiders returns a mapping from Agent ID -> ID of the bike that they are riding
	GetMegaBikeRiders() map[uuid.UUID]uuid.UUID
}
