package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"sort"

	"math"
	"slices"

	"github.com/google/uuid"
)

func (s *Server) RunSimLoop(iterations int) {

	s.ResetGameState()
	s.FoundingInstitutions()

	// run this for n iterations
	for i := 0; i < iterations; i++ {
		s.RunRoundLoop()
	}

}

func (s *Server) ResetGameState() {
	// kick everyone off bikes
	for id, _ := range s.GetAgentMap() {
		if bikeId, ok := s.megaBikeRiders[id]; ok {
			s.megaBikes[bikeId].RemoveAgent(id)
			delete(s.megaBikeRiders, id)
		}
	}

	// respawn people who died in previous round (conditional)
	if utils.RespawnAtRound {
		for _, agent := range s.deadAgents {
			s.AddAgent(agent)
		}
	}

	// replenish energy (conditional)
	if utils.ReplenishEnergyAtRound {
		for _, agent := range s.GetAgentMap() {
			agent.UpdateEnergyLevel(1.0)
		}
	}

	// empty the dead agent map
	clear(s.deadAgents)

	// zero the points (conditional)
	if utils.ResetPointsAtRound {
		for _, agent := range s.GetAgentMap() {
			agent.ResetPoints()
		}
	}

	s.replenishLootBoxes()
	s.replenishMegaBikes()
}

func (s *Server) FoundingInstitutions() {
	// Say which goverance method you might choose

	// check which governance method is chosen for each biker
	s.foundingChoices = make(map[uuid.UUID]utils.Governance)
	for id, agent := range s.GetAgentMap() {
		// collect choice from each agent
		choice := agent.DecideGovernance()
		s.foundingChoices[id] = choice
	}

	// tally the choices
	// FoundingChoices := make([]voting.IVoter, len(allAllocations))
	// FoundingAllocations is a map of governance method to number of agents that want that governance method
	foundingTotals, _ := voting.TallyFoundingVotes(s.foundingChoices)

	// for each governance method, populate megabikes with the bikers who chose that governance method
	govBikes := make(map[utils.Governance][]uuid.UUID)
	bikesUsed := make([]uuid.UUID, 0)
	for governanceMethod, numBikers := range foundingTotals {
		megaBikesNeeded := int(math.Ceil(float64(numBikers) / float64(utils.BikersOnBike)))
		govBikes[governanceMethod] = make([]uuid.UUID, megaBikesNeeded)
		// get bikes for this governance
		for i := 0; i < megaBikesNeeded; i++ {
			foundBike := false
			for !foundBike {
				bike := s.GetRandomBikeId()
				if !slices.Contains(bikesUsed, bike) {
					foundBike = true
					bikesUsed = append(bikesUsed, bike)
					govBikes[governanceMethod][i] = bike

					// set the governance
					bikeObj := s.GetMegaBikes()[bike]
					bikeObj.SetGovernance(governanceMethod)
				}
			}
		}

		for agent, governance := range s.foundingChoices {
			// randomly select a biker from the bikers who chose this governance method
			// add that biker to a megabike

			// select a bike with this governance method which has been assigned the lowest amount of bikers
			bikesAvailable := govBikes[governance]
			sort.Slice(bikesAvailable, func(i, j int) bool {
				// in the order from large to small
				return len(s.GetMegaBikes()[bikesAvailable[i]].GetAgents()) < len(s.GetMegaBikes()[bikesAvailable[j]].GetAgents())
			})

			chosenBike := bikesAvailable[0]
			// add agent to bike
			agentInt := s.GetAgentMap()[agent]
			s.GetMegaBikes()[chosenBike].AddAgent(agentInt)
		}
	}

	// if there are more bikers for a governance method than there are seats, then evenly distribute them across megabikes

	// set governance method for each bike so that it stays with the bike during the round

	// bikers comply with governance method on the bike they're on

	// choose leader if required

}

func (s *Server) Start() {
	fmt.Printf("Server initialised with %d agents \n\n", len(s.GetAgentMap()))
	gameStates := make([]GameStateDump, 0, s.GetIterations())
	for i := 0; i < s.GetIterations(); i++ {
		fmt.Printf("Game Loop %d running... \n \n", i)
		fmt.Printf("Main game loop running...\n\n")
		s.deadAgents = make(map[uuid.UUID]objects.IBaseBiker)
		s.RunSimLoop(utils.RoundIterations)
		gameStates = append(gameStates, s.NewGameStateDump())
		fmt.Printf("\nMain game loop finished.\n\n")
		fmt.Printf("Messaging session started...\n\n")
		s.RunMessagingSession()
		fmt.Printf("\nMessaging session completed\n\n")
		fmt.Printf("Game Loop %d completed.\n", i)
	}
	s.outputResults(gameStates)
}