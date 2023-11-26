package server_test

import (
	"SOMAS2023/internal/server"
	"testing"

	"github.com/google/uuid"
)

func TestInitialize(t *testing.T) {

	it := 3
	s := server.Initialize(it)

	if len(s.GetAgentMap()) != server.BikerAgentCount {
		t.Error("Agents not properly instantiated")
	}

	if len(s.GetMegaBikes()) != server.MegaBikeCount {
		t.Error("mega bikes not properly instantiated")
	}

	if len(s.GetLootBoxes()) != server.LootBoxCount {
		t.Error("Mega bikes not properly instantiated")
	}

	if s.GetAudi().GetID() == uuid.Nil {
		t.Error("audi not properly instantiated")
	}

	s.RunGameLoop()
	s.Start()
}

func TestGetJoiningRequests(t *testing.T) {
	it := 3
	s := server.Initialize(it)

	// 1: get two bike ids
	targetBikes := make([]uuid.UUID, 2)

	i := 0
	for bikeId, _ := range s.GetMegaBikes() {
		if i == 2 {
			break
		}
		targetBikes[i] = bikeId
		i += 1
	}

	// 2: set one agent requesting the first bike and two other requesting the second one
	i = 0
	requests := make(map[uuid.UUID][]uuid.UUID)
	requests[targetBikes[0]] = make([]uuid.UUID, 1)
	requests[targetBikes[1]] = make([]uuid.UUID, 2)
	for _, agent := range s.GetAgentMap() {
		if i == 0 {
			agent.ToggleOnBike()
			agent.SetBike(targetBikes[0])
			requests[targetBikes[0]][0] = agent.GetID()
		} else if i <= 2 {
			agent.ToggleOnBike()
			agent.SetBike(targetBikes[1])
			requests[targetBikes[1]][i-1] = agent.GetID()
		} else {
			break
		}
		i += 1
	}

	// 3. check that joining requests reflect the previous actions
	bikeRequests := s.GetJoiningRequests()
	if len(bikeRequests) != len(requests) {
		t.Error("bike requests processed incorrectly: empty")
	}

	for bikeId, agentIds := range bikeRequests {
		if len(agentIds) != len(requests[bikeId]) {
			t.Error("bike requests processed incorrectly: wrong number of agents for given bike")
		}
	}
}
