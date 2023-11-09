package server

import (
	"SOMAS2023/internal/common/objects"
	"fmt"
)

func (s *Server) RunGameLoop() {
	// Reset all the forces acting on the bike
	for _, bike := range s.megaBikes {
		bike.ResetAgentForces()
	}

	// Each agent makes a decision
	for id, agent := range s.GetAgentMap() {
		fmt.Printf("Agent %s updating state \n", id)
		agent.UpdateAgentInternalState()
		switch agent.DecideAction(s) {
		case objects.Pedal:
			force := agent.DecideForce(s)
			if bikeId, ok := s.megaBikeRiders[agent.GetID()]; ok {
				s.megaBikes[bikeId].AddAgentForce(force)
			} else {
				panic("agent tried to move when it was not on a bike")
			}
		case objects.ChangeBike:
			newBikeId := agent.ChangeBike(s)
			s.megaBikeRiders[agent.GetID()] = newBikeId
		default:
			panic("agent decided invalid action")
		}
	}

	// Replenish objects
	s.replenishLootBoxes()
	s.replenishMegaBikes()
}

func (s *Server) Start() {
	fmt.Printf("Server initialised with %d agents \n\n", len(s.GetAgentMap()))
	for i := 0; i < s.GetIterations(); i++ {
		fmt.Printf("Game Loop %d running... \n \n", i)
		fmt.Printf("Main game loop running...\n\n")
		s.RunGameLoop()
		fmt.Printf("\nMain game loop finished.\n\n")
		fmt.Printf("Messaging session started...\n\n")
		s.RunMessagingSession()
		fmt.Printf("\nMessaging session completed\n\n")
		fmt.Printf("Game Loop %d completed.\n", i)
	}
}
