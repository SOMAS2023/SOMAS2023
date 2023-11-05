package lootbox

import (
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type LootBox struct {
	colour Colour
}

// NewLootBox is a constructor for LootBox that initializes it with a new UUID and default position.
func NewLootBox(colour Colour) *LootBox {
	return &LootBox{
		ID:       uuid.New(),          // Generate a new unique identifier
		Position: utils.Coordinates{}, // Initialize to default position, modify as necessary
		colour:   colour,
	}
}

// SetColour sets the color of the BikerAgent.
func (lb *LootBox) SetColour(lootBoxColour Colour) {
	lb.colour = lootBoxColour
}

// GetColour returns the color of the BikerAgent.
func (lb *LootBox) GetColour() Colour {
	return lb.colour
}
