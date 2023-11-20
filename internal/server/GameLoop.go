package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"fmt"

	"github.com/google/uuid"
)

func (s *Server) RunGameLoop() {
	// Each agent makes a decision
	for agentId, agent := range s.GetAgentMap() {
		fmt.Printf("Agent %s updating state \n", agentId)
		agent.UpdateGameState(s)
		agent.UpdateAgentInternalState()
		switch agent.DecideAction() {
		case objects.Pedal:
			agent.DecideForce()
		case objects.ChangeBike:
			s.SetBikerBike(agent, agent.ChangeBike())
		default:
			panic("agent decided invalid action")
		}
	}

	// The Audi makes a decision
	s.audi.UpdateGameState(s)

	// Move the mega bikes
	for _, bike := range s.GetMegaBikes() {
		// Server requests megabikes to update their force and orientation based on agents pedaling
		bike.UpdateForce()
		force := bike.GetForce()
		bike.UpdateOrientation()
		orientation := bike.GetOrientation()

		// Obtains the current state (i.e. velocity, acceleration, position, mass)
		initialState := bike.GetPhysicalState()

		// Generates a new state based on the force and orientation of the bike
		finalState := physics.GenerateNewState(initialState, force, orientation)

		// Sets the new physical state (i.e. updates gamestate)
		bike.SetPhysicalState(finalState)
	}

	// Move the audi
	s.audi.UpdateForce()
	force := s.audi.GetForce()
	s.audi.UpdateOrientation()
	orientation := s.audi.GetOrientation()
	initialState := s.audi.GetPhysicalState()
	finalState := physics.GenerateNewState(initialState, force, orientation)
	s.audi.SetPhysicalState(finalState)

	// Lootbox Distribution
	s.LootboxCheckAndDistributions()

	// Replenish objects
	s.replenishLootBoxes()
	s.replenishMegaBikes()
}

func (s *Server) LootboxCheckAndDistributions() {
	for bikeid, megabike := range s.GetMegaBikes() {
		for lootid, lootbox := range s.GetLootBoxes() {
			if megabike.CheckForCollision(lootbox) {
				// Collision detected
				fmt.Printf("Collision detected between MegaBike %s and LootBox %s \n", bikeid, lootid)
				agents := megabike.GetAgents()
				totAgents := len(agents)

				// Compute resource allocation share.
				accVotes := make(map[uuid.UUID]float64)
				for _, agent := range agents {
					agentVote := agent.GetResourceVote()

					// Normalize votes.
					sum := 0.0
					for _, vote := range agentVote {
						sum += vote
					}
					if sum != 0 {
						for agentID, vote := range agentVote {
							agentVote[agentID] = vote / sum
						}
					}

					// Accumulate votes.
					for agentID, vote := range agentVote {
						if _, exists := accVotes[agentID]; !exists {
							accVotes[agentID] = 0
						}
						accVotes[agentID] += vote
					}
				}

				// Normalize resource allocation share.
				if totAgents != 0 {
					for agentID, vote := range accVotes {
						accVotes[agentID] = vote / float64(totAgents)
					}
				}

				// Distribute loot.
				for _, agent := range agents {
					lootShare, exists := accVotes[agent.GetID()]
					if exists {
						lootShare *= lootbox.GetTotalResources()

						fmt.Printf("Agent %s allocated %f loot \n", agent.GetID(), lootShare)
						agent.UpdateResourceAppropriation(lootShare)
						agent.UpdateEnergyLevel(lootShare)
					}
				}
			}
		}
	}
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
