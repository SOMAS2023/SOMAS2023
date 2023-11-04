package lootbox

type LootBox struct {
	colour Colour
}

// SetColour sets the color of the BikerAgent.
func (lb *LootBox) SetColour(lootBoxColour Colour) {
	lb.colour = lootBoxColour
}

// GetColour returns the color of the BikerAgent.
func (lb *LootBox) GetColour() Colour {
	return lb.colour
}
