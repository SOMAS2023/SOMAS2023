package objects

import (
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type LootBox struct {
	colour      utils.Colour
	coordinates utils.Coordinates
	id          uuid.UUID
}

// NewLootBox is a constructor for LootBox that initializes it with a new UUID and default position.
func NewLootBox(colour utils.Colour) *LootBox {
	return &LootBox{
		id:          uuid.New(),                        // Generate a new unique identifier
		coordinates: utils.GenerateRandomCoordinates(), // Initialize to randomized position
		colour:      utils.GenerateRandomColour(),      // Initialize to randomized colour
	}
}

// returns the unique ID of the object
func (lb *LootBox) GetID() uuid.UUID {
	return lb.id
}

// returns the current coordinates of the object
func (lb *LootBox) GetPosition() utils.Coordinates {
	return lb.coordinates
}

// SetColour sets the color of the BikerAgent.
func (lb *LootBox) SetColour(lootBoxColour utils.Colour) {
	lb.colour = lootBoxColour
}

// GetColour returns the color of the BikerAgent.
func (lb *LootBox) GetColour() utils.Colour {
	return lb.colour
}
