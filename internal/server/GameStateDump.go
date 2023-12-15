package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"maps"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type GameStateDump struct {
	Iteration int                       `json:"iteration"`
	Agents    map[uuid.UUID]AgentDump   `json:"agents"`
	Bikes     map[uuid.UUID]BikeDump    `json:"bikes"`
	LootBoxes map[uuid.UUID]LootBoxDump `json:"loot_boxes"`
	Awdis     []AwdiDump                `json:"awdis"`
}

type PhysicsObjectDump struct {
	ID            uuid.UUID           `json:"-"`
	PhysicalState utils.PhysicalState `json:"physical_state"`
	Orientation   float64             `json:"orientation"`
	Force         float64             `json:"force"`
}

type BikeDump struct {
	PhysicsObjectDump
	Agents     []AgentDump      `json:"-"`
	AgentIDs   []uuid.UUID      `json:"agent_ids"`
	Governance utils.Governance `json:"governance"`
	Ruler      uuid.UUID        `json:"ruler"`
}

type AgentDump struct {
	ID           uuid.UUID             `json:"-"`
	Class        string                `json:"class"`
	Forces       utils.Forces          `json:"forces"`
	EnergyLevel  float64               `json:"energy_level"`
	Points       int                   `json:"points"`
	Colour       utils.Colour          `json:"-"`
	ColourString string                `json:"colour"`
	Location     utils.Coordinates     `json:"location"`
	OnBike       bool                  `json:"on_bike"`
	BikeID       uuid.UUID             `json:"bike_id"`
	Reputation   map[uuid.UUID]float64 `json:"reputation"`
	GroupID      int                   `json:"group_id"`
}

type LootBoxDump struct {
	PhysicsObjectDump
	TotalResources float64      `json:"total_resources"`
	Colour         utils.Colour `json:"-"`
	ColourString   string       `json:"colour"`
}

type AwdiDump struct {
	PhysicsObjectDump
	ID         uuid.UUID `json:"id"`
	TargetBike uuid.UUID `json:"target_bike"`
}

func newPhysicsObjectDump(physicsObject objects.IPhysicsObject) PhysicsObjectDump {
	return PhysicsObjectDump{
		ID:            physicsObject.GetID(),
		PhysicalState: physicsObject.GetPhysicalState(),
		Orientation:   physicsObject.GetOrientation(),
		Force:         physicsObject.GetForce(),
	}
}

// the updated game state dump will be used both by the visualiser (at the end of a round)
// and by the agents, as a representation of the state of the game (to be updated every time anything
// changes in the game, such as agents leaving a bike, bikes moving or lootboxes being looted)
func (s *Server) NewGameStateDump(iteration int) GameStateDump {
	agents := make(map[uuid.UUID]AgentDump, len(s.GetAgentMap()))
	for id, agent := range s.GetAgentMap() {
		var location utils.Coordinates
		if agent.GetBike() != uuid.Nil {
			location = s.megaBikes[agent.GetBike()].GetPosition()
		} else {
			location = utils.Coordinates{X: 0.0, Y: 0.0}
		}
		agents[id] = AgentDump{
			ID:           agent.GetID(),
			Class:        strings.TrimPrefix(reflect.TypeOf(agent).String(), "*"),
			Forces:       agent.GetForces(),
			EnergyLevel:  agent.GetEnergyLevel(),
			Points:       agent.GetPoints(),
			Colour:       agent.GetColour(),
			ColourString: agent.GetColour().String(),
			Location:     location,
			OnBike:       agent.GetBikeStatus(),
			BikeID:       agent.GetBike(),
			Reputation:   maps.Clone(agent.GetReputation()),
			GroupID:      agent.GetGroupID(),
		}
	}

	bikes := make(map[uuid.UUID]BikeDump, len(s.megaBikes))
	for id, bike := range s.megaBikes {
		agentDumps := make([]AgentDump, 0, len(bike.GetAgents()))
		agentIDs := make([]uuid.UUID, 0, len(bike.GetAgents()))
		for _, agent := range bike.GetAgents() {
			agentDumps = append(agentDumps, agents[agent.GetID()])
			agentIDs = append(agentIDs, agent.GetID())
		}
		bikes[id] = BikeDump{
			PhysicsObjectDump: newPhysicsObjectDump(bike),
			Agents:            agentDumps,
			AgentIDs:          agentIDs,
			Governance:        bike.GetGovernance(),
			Ruler:             bike.GetRuler(),
		}
	}

	lootBoxes := make(map[uuid.UUID]LootBoxDump, len(s.lootBoxes))
	for id, lootBox := range s.lootBoxes {
		lootBoxes[id] = LootBoxDump{
			PhysicsObjectDump: newPhysicsObjectDump(lootBox),
			TotalResources:    lootBox.GetTotalResources(),
			Colour:            lootBox.GetColour(),
			ColourString:      lootBox.GetColour().String(),
		}
	}

	return GameStateDump{
		Iteration: iteration,
		Agents:    agents,
		Bikes:     bikes,
		LootBoxes: lootBoxes,
		Awdis: []AwdiDump{{
			PhysicsObjectDump: newPhysicsObjectDump(s.awdi),
			ID:                s.awdi.GetID(),
			TargetBike:        s.awdi.GetTargetID(),
		}},
	}
}
