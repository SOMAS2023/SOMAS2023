package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
)

type GameStateDump struct {
	Bikes     []BikeDump    `json:"bikes"`
	LootBoxes []LootBoxDump `json:"loot_boxes"`
	Audi      AudiDump      `json:"audi"`
}

type PhysicsObjectDump struct {
	ID            uuid.UUID           `json:"id"`
	PhysicalState utils.PhysicalState `json:"physical_state"`
}

type BikeDump struct {
	PhysicsObjectDump
	Agents []AgentDump `json:"agents"`
}

type AgentDump struct {
	ID                       uuid.UUID                        `json:"id"`
	Forces                   utils.Forces                     `json:"forces"`
	EnergyLevel              float64                          `json:"energy_level"`
	ResourceAllocationParams objects.ResourceAllocationParams `json:"resource_allocation_params"`
	Colour                   string                           `json:"colour"`
	Location                 utils.Coordinates                `json:"location"`
}

type LootBoxDump struct {
	PhysicsObjectDump
	TotalResources float64 `json:"total_resources"`
	Colour         string  `json:"colour"`
}

type AudiDump struct {
	PhysicsObjectDump
	TargetBike uuid.UUID `json:"target_bike"`
}

func newPhysicsObjectDump(physicsObject objects.IPhysicsObject) PhysicsObjectDump {
	return PhysicsObjectDump{
		ID:            physicsObject.GetID(),
		PhysicalState: physicsObject.GetPhysicalState(),
	}
}

func (s *Server) NewGameStateDump() GameStateDump {
	bikes := make([]BikeDump, 0, len(s.megaBikes))
	for _, bike := range s.megaBikes {
		agents := make([]AgentDump, 0, len(bike.GetAgents()))
		for _, agent := range bike.GetAgents() {
			agents = append(agents, AgentDump{
				ID:                       agent.GetID(),
				Forces:                   agent.GetForces(),
				EnergyLevel:              agent.GetEnergyLevel(),
				ResourceAllocationParams: agent.GetResourceAllocationParams(),
				Colour:                   agent.GetColour().String(),
				Location:                 agent.GetLocation(),
			})
		}
		bikes = append(bikes, BikeDump{
			PhysicsObjectDump: newPhysicsObjectDump(bike),
			Agents:            agents,
		})
	}

	lootBoxes := make([]LootBoxDump, 0, len(s.lootBoxes))
	for _, lootBox := range s.lootBoxes {
		lootBoxes = append(lootBoxes, LootBoxDump{
			PhysicsObjectDump: newPhysicsObjectDump(lootBox),
			TotalResources:    lootBox.GetTotalResources(),
			Colour:            lootBox.GetColour().String(),
		})
	}

	return GameStateDump{
		Bikes:     bikes,
		LootBoxes: lootBoxes,
		Audi: AudiDump{
			PhysicsObjectDump: newPhysicsObjectDump(s.audi),
			TargetBike:        s.audi.GetTargetID(),
		},
	}
}
