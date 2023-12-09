package server_test

import (
	"SOMAS2023/internal/server"
	"fmt"
	"testing"

	"slices"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetJoiningRequests(t *testing.T) {
	it := 3
	s := server.Initialize(it)

	// 1: get two bike ids
	targetBikes := make([]uuid.UUID, 2)

	i := 0
	for bikeId := range s.GetMegaBikes() {
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
	bikeRequests := s.GetJoiningRequests(make([]uuid.UUID, 0))
	if len(bikeRequests) != len(requests) {
		t.Error("bike requests processed incorrectly: empty")
	}

	for bikeId, agentIds := range bikeRequests {
		if len(agentIds) != len(requests[bikeId]) {
			t.Error("bike requests processed incorrectly: wrong number of agents for given bike")
		}
	}
	fmt.Printf("\nJoining request passed \n")
}

func TestGetJoiningRequestsWithLimbo(t *testing.T) {
	it := 3
	s := server.Initialize(it)

	// 1: get two bike ids
	targetBikes := make([]uuid.UUID, 2)

	i := 0
	for bikeId := range s.GetMegaBikes() {
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
	requests[targetBikes[1]] = make([]uuid.UUID, 1)
	limbo := make([]uuid.UUID, 1)
	for _, agent := range s.GetAgentMap() {
		if i == 0 {
			agent.ToggleOnBike()
			agent.SetBike(targetBikes[0])
			requests[targetBikes[0]][0] = agent.GetID()
		} else if i == 1 {
			// add it to second bike for request
			agent.ToggleOnBike()
			agent.SetBike(targetBikes[1])
			requests[targetBikes[1]][i-1] = agent.GetID()
		} else if i == 2 {
			//remove it from bike but add it to limbo (to mimick request made in this turn)
			agent.ToggleOnBike()
			agent.SetBike(targetBikes[1])
			limbo[0] = agent.GetID()
		} else {
			break
		}
		i += 1
	}

	// 3. check that joining requests reflect the previous actions
	bikeRequests := s.GetJoiningRequests(limbo)
	assert.Equal(t, len(bikeRequests), len(requests), "bike requests processed incorrectly: empty")

	for bikeId, agentIds := range bikeRequests {
		assert.Equal(t, len(agentIds), len(requests[bikeId]), "bike requests processed incorrectly: wrong number of agents for given bike")
		assert.False(t, slices.Contains(agentIds, limbo[0]), "bike requests processed incorrectly: agent in limbo is requesting a bike")
	}

	fmt.Printf("\nJoining request passed \n")
}

func TestGetRandomID(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	bike := s.GetRandomBikeId()
	_, exists := s.GetMegaBikes()[bike]
	if !exists {
		t.Error("returned bike is not in ")
	}
	fmt.Printf("\nGet random ID passed \n")
}

func TestAddAgentToBike(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	bike := s.GetRandomBikeId()
	var changedAgent uuid.UUID
	for agentID, agent := range s.GetAgentMap() {
		agent.SetBike(bike)
		s.AddAgentToBike(agent)
		changedAgent = agentID
		break
	}

	agentToCheck := s.GetAgentMap()[changedAgent]
	if agentToCheck.GetBike() != bike {
		t.Error("agent's bike is not as expected")
	}
	fmt.Printf("\nSet biker bike passed \n")
}
