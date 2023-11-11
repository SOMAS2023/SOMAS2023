package server

import (
	utils "SOMAS2023/internal/common/utils"
	"fmt"

	"github.com/google/uuid"
)

func (s *Server) RunGameLoop() {
	for id, agent := range s.GetAgentMap() {
		fmt.Printf("Agent %s updating state \n", id)
		agent.UpdateAgentInternalState()
	}
	s.replenishLootBoxes()
	s.replenishMegaBikes()
}

func (s *Server) LootboxCheckAndDistributions() {
	EPSILON := 0.1          // Collision detection threshold
	COLOURMULTIPLIER := 1.2 // Multiplier for colour match
	// a > c > b
	a := 1.0
	b := 0.5
	c := 0.3
	/**
	Agents to focus on just meeting needs without hoarding: a to be high relative to b and c.
	Discourage agents from appropriating less than their needs: c high.
	Encourage agents to gather as many resources as possible: b high.
	**/
	for bikeid, megabike := range s.GetMegaBikes() {
		for lootid, lootbox := range s.GetLootBoxes() {
			if utils.CheckCollision(megabike.GetPosition(), lootbox.GetPosition(), EPSILON) {
				// Collision detected
				fmt.Printf("Collision detected between MegaBike %s and LootBox %s \n", bikeid, lootid)
				agents := megabike.GetAgents()
				totalUtility := 0.0
				utilityMap := make(map[uuid.UUID]float64)
				// Get total amount of attributes for each agent on the bike
				totalG := 0.0
				totalQ := 0.0
				totalP := 0.0
				totalR := 0.0
				for _, agent := range agents {
					totalG += agent.GetEnergyLevel()
					totalQ += agent.GetResourceNeed()
					totalP += agent.GetResourceProvision()
					totalR += agent.GetResourceAppropriation()
				}
				// Calculate utility for each agent on the bike
				for _, agent := range agents {
					g := agent.GetEnergyLevel() / totalG
					q := agent.GetResourceNeed() / totalQ
					p := agent.GetResourceProvision() / totalP
					r := agent.GetResourceAppropriation() / totalR
					R := r + (g - p) // Accrued resources

					var u float64 // Utility
					if R >= q {
						u = a*q + b*(R-q)
					} else {
						u = a*R - c*(q-R)
					}
					if agent.GetColour() == lootbox.GetColour() {
						u *= COLOURMULTIPLIER
					}
					utilityMap[agent.GetID()] = u
					totalUtility += u
				}
				// Distribute loot based on utility as a share of the total utility
				for _, agent := range agents {
					utilityShare := utilityMap[agent.GetID()] / totalUtility
					lootShare := utilityShare * lootbox.GetTotalResources()
					// Allocate loot based on the calculated utility share
					fmt.Printf("Agent %s allocated %f loot \n", agent.GetID(), lootShare)
					agent.SetEnergyLevel(lootShare)
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
