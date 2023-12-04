package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"maps"

	"github.com/google/uuid"
)

func (gs GameStateDump) GetLootBoxes() map[uuid.UUID]objects.ILootBox {
	result := make(map[uuid.UUID]objects.ILootBox)
	for id, lb := range gs.LootBoxes {
		result[id] = lb
	}
	return result
}

func (gs GameStateDump) GetMegaBikes() map[uuid.UUID]objects.IMegaBike {
	result := make(map[uuid.UUID]objects.IMegaBike)
	for id, mb := range gs.Bikes {
		result[id] = mb
	}
	return result
}

func (gs GameStateDump) GetAgents() map[uuid.UUID]objects.IBaseBiker {
	result := make(map[uuid.UUID]objects.IBaseBiker)
	for id, a := range gs.Agents {
		result[id] = a
	}
	return result
}

func (gs GameStateDump) GetAudi() objects.IAudi {
	return gs.Audi
}

func (o PhysicsObjectDump) GetID() uuid.UUID {
	return o.ID
}

func (o PhysicsObjectDump) GetPosition() utils.Coordinates {
	return o.PhysicalState.Position
}

func (o PhysicsObjectDump) GetVelocity() float64 {
	return o.PhysicalState.Velocity
}

func (o PhysicsObjectDump) GetOrientation() float64 {
	return o.Orientation
}

func (o PhysicsObjectDump) GetForce() float64 {
	return o.Force
}

func (o PhysicsObjectDump) GetPhysicalState() utils.PhysicalState {
	return o.PhysicalState
}

func (a AgentDump) GetID() uuid.UUID {
	return a.ID
}

func (a AgentDump) GetForces() utils.Forces {
	return a.Forces
}

func (a AgentDump) GetColour() utils.Colour {
	return a.Colour
}

func (a AgentDump) GetLocation() utils.Coordinates {
	return a.Location
}

func (a AgentDump) GetBike() uuid.UUID {
	return a.BikeID
}

func (a AgentDump) GetEnergyLevel() float64 {
	return a.EnergyLevel
}

func (a AgentDump) GetPoints() int {
	return a.Points
}

func (a AgentDump) GetBikeStatus() bool {
	return a.OnBike
}

func (a AgentDump) GetReputation() map[uuid.UUID]float64 {
	return maps.Clone(a.Reputation)
}

func (b BikeDump) GetAgents() []objects.IBaseBiker {
	result := make([]objects.IBaseBiker, 0, len(b.Agents))
	for i := range b.Agents {
		result = append(result, b.Agents[i])
	}
	return result
}

func (b BikeDump) GetGovernance() utils.Governance {
	return b.Governance
}

func (b BikeDump) GetRuler() uuid.UUID {
	return b.Ruler
}

func (l LootBoxDump) GetTotalResources() float64 {
	return l.TotalResources
}

func (l LootBoxDump) GetColour() utils.Colour {
	return l.Colour
}

func (a AudiDump) GetTargetID() uuid.UUID {
	return a.TargetBike
}
