package objects

import (
	utils "SOMAS2023/internal/common/utils"
)

type ILootBox interface {
	IPhysicsObject
	GetTotalResources() float64
	GetColour() utils.Colour
}

type LootBox struct {
	*PhysicsObject
	colour    utils.Colour
	totalLoot float64
}

// GetLootBox is a constructor for LootBox that initializes it with a new UUID and default position.
func GetLootBox() *LootBox {
	return &LootBox{
		PhysicsObject: GetPhysicsObject(0),
		colour:        utils.GenerateRandomColour(),    // Initialize to randomized colour
		totalLoot:     utils.GenerateRandomFloat(2, 4), // Initialize to randomized totalLoot
	}
}

// returns the total loot of the object
func (lb *LootBox) GetTotalResources() float64 {
	return lb.totalLoot
}

// GetColour returns the color of the BikerAgent.
func (lb *LootBox) GetColour() utils.Colour {
	return lb.colour
}
