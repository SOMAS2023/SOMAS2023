package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"fmt"
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

				for _, agent := range agents {
					// this function allows the agent to decide on its allocation parameters
					// these are the parameters that we want to be considered while carrying out
					// the elected protocol for resource allocation
					agent.SetAllocationParameters()

					// in the MVP  the allocation parameters are ignored and
					// the utility share will simply be 1 / the number of agents on the bike
					utilityShare := 1.0 / float64(totAgents)
					lootShare := utilityShare * lootbox.GetTotalResources()
					// Allocate loot based on the calculated utility share
					fmt.Printf("Agent %s allocated %f loot \n", agent.GetID(), lootShare)
					agent.UpdateEnergyLevel(lootShare)
					if agent.GetEnergyLevel() < 0 {
						s.unaliveAgent(agent)
					}
				}
			}
		}
	}
}

func (s *Server) unaliveAgent(agent objects.IBaseBiker) {
	fmt.Printf("Agent %s got game ended\n", agent.GetID())
	s.RemoveAgent(agent)
	if bikeId, ok := s.megaBikeRiders[agent.GetID()]; ok {
		s.megaBikes[bikeId].RemoveAgent(agent.GetID())
		delete(s.megaBikeRiders, agent.GetID())
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
