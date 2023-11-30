package server_test

import (
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/server"
	"fmt"
	"testing"
)

func TestAudiCollisionProcess(t *testing.T) {
	it := 100
	s := server.Initialize(it)
	nAgentToDelete := 0
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(s)
	}
	for _, anyBike := range s.GetMegaBikes() {
		agentsOnBike := anyBike.GetAgents()
		nAgentToDelete = len(agentsOnBike)
		if agentsOnBike == nil {
			continue
		}
		// find any non-empty bike
		if len(agentsOnBike) > 0 {
			// send audi to it
			s.GetAudi().SetPhysicalState(utils.PhysicalState{Position: anyBike.GetPosition()})
		}
	}
	nAgentsBefore := len(s.GetAgentMap())
	nMegaBikesBefore := len(s.GetMegaBikes())
	s.AudiCollisionCheck()
	nAgentsAfter := len(s.GetAgentMap())
	nMegaBikesAfter := len(s.GetMegaBikes())

	// check if remove agents correctly
	if nAgentsBefore-nAgentsAfter != nAgentToDelete {
		fmt.Printf("Before audi collision, number of agents = %d \n", nAgentsBefore)
		fmt.Printf("After audi collision, number of agents = %d \n", nAgentsAfter)
		fmt.Printf("On bike collide with audi, number of agents = %d \n", nAgentToDelete)
		t.Error("Audi didnt remove agents correctly")
	}

	if utils.AudiRemoveMegaBike {
		// check if remove megaBike correctly
		if nMegaBikesBefore-nMegaBikesAfter != 1 {
			fmt.Printf("Before audi collision, number of megaBikes = %d \n", nMegaBikesBefore)
			fmt.Printf("After audi collision, number of megaBikes = %d \n", nMegaBikesAfter)
			t.Error("Audi didnt remove megaBike correctly")
		}
	}
	fmt.Printf("\nRun action process passed \n")
}
