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

func (s *Server) RunSimLoop(iterations int) []GameStateDump {

	s.ResetGameState()
	s.FoundingInstitutions()

	// run this for n iterations
	gameStates := []GameStateDump{s.NewGameStateDump(-1)}
	for i := 0; i < iterations; i++ {
		s.RunRoundLoop()
		gameStates = append(gameStates, s.NewGameStateDump(i))
	}

	return gameStates
}

func (s *Server) ResetGameState() {
	// kick everyone off bikes
	for _, agent := range s.GetAgentMap() {
		if agent.GetBike() != uuid.Nil {
			s.RemoveAgentFromBike(agent)
		}
	}

	// respawn people who died in previous round (conditional)
	if utils.RespawnEveryRound && utils.ReplenishEnergyEveryRound {
		for _, agent := range s.deadAgents {
			s.AddAgent(agent)
		}
	}

	// replenish energy (conditional)
	if utils.ReplenishEnergyEveryRound {
		for _, agent := range s.GetAgentMap() {
			agent.UpdateEnergyLevel(1.0)
		}
	}

	// empty the dead agent map
	clear(s.deadAgents)

	// zero the points (conditional)
	if utils.ResetPointsEveryRound {
		for _, agent := range s.GetAgentMap() {
			agent.ResetPoints()
		}
	}

	s.replenishLootBoxes()
	s.replenishMegaBikes()
}

func (s *Server) FoundingInstitutions() {
	// Say which goverance method you might choose

	// run founding messaging session
	s.UpdateGameStates()
	s.RunMessagingSession()

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
	}

	for agent, governance := range s.foundingChoices {
		// randomly select a biker from the bikers who chose this governance method
		// add that biker to a megabike
		// if there are more bikers for a governance method than there are seats, then evenly distribute them across megabikes

		// select a bike with this governance method which has been assigned the lowest amount of bikers
		bikesAvailable := govBikes[governance]
		sort.Slice(bikesAvailable, func(i, j int) bool {
			// in the order from large to small
			return len(s.GetMegaBikes()[bikesAvailable[i]].GetAgents()) < len(s.GetMegaBikes()[bikesAvailable[j]].GetAgents())
		})

		// get the first one of the sorted bikes
		chosenBike := bikesAvailable[0]
		// add agent to bike
		agentInt := s.GetAgentMap()[agent]
		agentInt.SetBike(chosenBike)
		s.AddAgentToBike(agentInt)
	}

	s.UpdateGameStates()
	// run election process for Leadership and Dictatorship bikes
	for _, bike := range s.GetMegaBikes() {
		gov := bike.GetGovernance()
		agents := bike.GetAgents()
		if (gov == utils.Leadership || gov == utils.Dictatorship) && len(agents) != 0 {
			ruler := s.RulerElection(agents, gov)
			bike.SetRuler(ruler)
		}
	}

	s.UpdateGameStates()

}

func (s *Server) Start() {
	fmt.Printf("Server initialised with %d agents \n\n", len(s.GetAgentMap()))
	gameStates := make([][]GameStateDump, 0, s.GetIterations())
	s.deadAgents = make(map[uuid.UUID]objects.IBaseBiker)
	for i := 0; i < s.GetIterations(); i++ {
		fmt.Printf("Game Loop %d running... \n \n", i)
		fmt.Printf("Main game loop running...\n\n")
		gameStates = append(gameStates, s.RunSimLoop(utils.RoundIterations))
		fmt.Printf("\nMain game loop finished.\n\n")
		fmt.Printf("Messaging session started...\n\n")
		s.RunMessagingSession()
		fmt.Printf("\nMessaging session completed\n\n")
		fmt.Printf("Game Loop %d completed.\n", i)
	}
	s.outputResults(gameStates)
}
