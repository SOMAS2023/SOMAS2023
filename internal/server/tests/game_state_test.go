package server_test

import (
	"SOMAS2023/internal/server"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

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
	fmt.Printf("\n Joining request passed \n")
}

func TestGetRandomID(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	bike := s.GetRandomBikeId()
	_, exists := s.GetMegaBikes()[bike]
	if !exists {
		t.Error("returned bike is not in ")
	}
	fmt.Printf("\n Get random ID passed \n")
}

func TestSetBikerBike(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	bike := s.GetRandomBikeId()
	var changedAgent uuid.UUID
	for agentID, agent := range s.GetAgentMap() {
		s.SetBikerBike(agent, bike)
		changedAgent = agentID
		break
	}

	agentToCheck := s.GetAgentMap()[changedAgent]
	if agentToCheck.GetBike() != bike {
		t.Error("agent's bike is not as expected")
	}
	fmt.Printf("\n Set biker bike passed \n")
}
