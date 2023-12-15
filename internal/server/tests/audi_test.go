package server_test

import (
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/server"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestAwdiCollisionProcess(t *testing.T) {
	it := 2
	s := server.Initialize(it)
	nAgentToDelete := 0
	for _, anyBike := range s.GetMegaBikes() {
		agentsOnBike := anyBike.GetAgents()
		nAgentToDelete = len(agentsOnBike)
		if agentsOnBike == nil {
			continue
		}
		// find any non-empty bike
		if len(agentsOnBike) > 0 {
			// send awdi to it
			s.GetAwdi().SetPhysicalState(utils.PhysicalState{Position: anyBike.GetPosition()})
		}
	}
	nAgentsBefore := len(s.GetAgentMap())
	nMegaBikesBefore := len(s.GetMegaBikes())
	s.AwdiCollisionCheck()
	nAgentsAfter := len(s.GetAgentMap())
	nMegaBikesAfter := len(s.GetMegaBikes())

	// check if remove agents correctly
	if nAgentsBefore-nAgentsAfter != nAgentToDelete {
		fmt.Printf("Before awdi collision, number of agents = %d \n", nAgentsBefore)
		fmt.Printf("After awdi collision, number of agents = %d \n", nAgentsAfter)
		fmt.Printf("On bike collide with awdi, number of agents = %d \n", nAgentToDelete)
		t.Error("Awdi didnt remove agents correctly")
	}

	if utils.AwdiRemovesMegaBike {
		// check if remove megaBike correctly
		if nMegaBikesBefore-nMegaBikesAfter != 1 {
			fmt.Printf("Before awdi collision, number of megaBikes = %d \n", nMegaBikesBefore)
			fmt.Printf("After awdi collision, number of megaBikes = %d \n", nMegaBikesAfter)
			t.Error("Awdi didnt remove megaBike correctly")
		}
	}
	fmt.Printf("\nRun action process passed \n")
}

func TestAwdiTargeting(t *testing.T) {
	it := 1
	s := server.Initialize(it)
	// required otherwise agents are not initialized to bikes
	gs := s.NewGameStateDump(0)
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}
	s.FoundingInstitutions()
	i := 0
	emptyBikeId := uuid.UUID{}
	slowBikeId := uuid.UUID{}
	for id, bike := range s.GetMegaBikes() {
		if i == 0 {
			// remove agents on it
			for _, agents := range bike.GetAgents() {
				bike.RemoveAgent(agents.GetID())
			}
			// stop the bike
			bike.SetPhysicalState(utils.PhysicalState{Velocity: 0.0})
			fmt.Printf("Megabike{%s} has {%d} agents with velocity {%.2f}\n", id, len(bike.GetAgents()), bike.GetVelocity())
			emptyBikeId = id
		} else if i == 1 {
			// give the bike a slow Velocity
			bike.SetPhysicalState(utils.PhysicalState{Velocity: 1.0})
			fmt.Printf("Megabike{%s} has {%d} agents with velocity {%.2f}\n", id, len(bike.GetAgents()), bike.GetVelocity())
			slowBikeId = id
		} else if i == 2 {
			// give the bike a fast Velocity
			bike.SetPhysicalState(utils.PhysicalState{Velocity: 100.0})
			agentsOnBike := bike.GetAgents()
			if len(agentsOnBike) == 0 {
				emptyBikeId = id
			}
			fmt.Printf("Megabike{%s} has {%d} agents with velocity {%.2f}\n", id, len(bike.GetAgents()), bike.GetVelocity())
		}
		i += 1
	}
	s.GetAwdi().UpdateGameState(gs)
	s.GetAwdi().UpdateForce()
	targetId := s.GetAwdi().GetTargetID()
	fmt.Printf("Awdi is targeting {%s}\n", targetId)
	if utils.AwdiTargetsEmptyMegaBike {
		if targetId != emptyBikeId {
			t.Error("Awdi didnt target empty megaBike!")
		}
	}
	if utils.AwdiOnlyTargetsStationaryMegaBike {
		if targetId == slowBikeId {
			t.Error("Awdi didnt ignore moving megaBike!")
		}
	}
	fmt.Printf("\nRun action process passed \n")
}
