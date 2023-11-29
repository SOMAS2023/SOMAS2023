package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
)

type GameStateDump struct {
	Agents    map[uuid.UUID]AgentDump   `json:"agents"`
	Bikes     map[uuid.UUID]BikeDump    `json:"bikes"`
	LootBoxes map[uuid.UUID]LootBoxDump `json:"loot_boxes"`
	Audi      AudiDump                  `json:"audi"`
}

type PhysicsObjectDump struct {
	PhysicalState utils.PhysicalState `json:"physical_state"`
}

type BikeDump struct {
	PhysicsObjectDump
	AgentIDs []uuid.UUID `json:"agent_ids"`
}

type AgentDump struct {
	Forces                   utils.Forces                     `json:"forces"`
	EnergyLevel              float64                          `json:"energy_level"`
	Points                   int                              `json:"points"`
	ResourceAllocationParams objects.ResourceAllocationParams `json:"resource_allocation_params"`
	Colour                   string                           `json:"colour"`
	Location                 utils.Coordinates                `json:"location"`
	OnBike                   bool                             `json:"on_bike"`
	BikeID                   uuid.UUID                        `json:"bike_id"`
}

type LootBoxDump struct {
	PhysicsObjectDump
	TotalResources float64 `json:"total_resources"`
	Colour         string  `json:"colour"`
}

type AudiDump struct {
	PhysicsObjectDump
	ID         uuid.UUID `json:"id"`
	TargetBike uuid.UUID `json:"target_bike"`
}

func newPhysicsObjectDump(physicsObject objects.IPhysicsObject) PhysicsObjectDump {
	return PhysicsObjectDump{
		PhysicalState: physicsObject.GetPhysicalState(),
	}
}

func (s *Server) NewGameStateDump() GameStateDump {
	agents := make(map[uuid.UUID]AgentDump, len(s.GetAgentMap()))
	for id, agent := range s.GetAgentMap() {
		agents[id] = AgentDump{
			Forces:                   agent.GetForces(),
			EnergyLevel:              agent.GetEnergyLevel(),
			Points:                   agent.GetPoints(),
			ResourceAllocationParams: agent.GetResourceAllocationParams(),
			Colour:                   agent.GetColour().String(),
			Location:                 agent.GetLocation(),
			OnBike:                   agent.GetBikeStatus(),
			BikeID:                   agent.GetBike(),
		}
	}

	bikes := make(map[uuid.UUID]BikeDump, len(s.megaBikes))
	for id, bike := range s.megaBikes {
		agentIDs := make([]uuid.UUID, 0, len(bike.GetAgents()))
		for _, agent := range bike.GetAgents() {
			agentIDs = append(agentIDs, agent.GetID())
		}
		bikes[id] = BikeDump{
			PhysicsObjectDump: newPhysicsObjectDump(bike),
			AgentIDs:          agentIDs,
		}
	}

	lootBoxes := make(map[uuid.UUID]LootBoxDump, len(s.lootBoxes))
	for id, lootBox := range s.lootBoxes {
		lootBoxes[id] = LootBoxDump{
			PhysicsObjectDump: newPhysicsObjectDump(lootBox),
			TotalResources:    lootBox.GetTotalResources(),
			Colour:            lootBox.GetColour().String(),
		}
	}

	return GameStateDump{
		Agents:    agents,
		Bikes:     bikes,
		LootBoxes: lootBoxes,
		Audi: AudiDump{
			PhysicsObjectDump: newPhysicsObjectDump(s.audi),
			ID:                s.audi.GetID(),
			TargetBike:        s.audi.GetTargetID(),
		},
	}
}
