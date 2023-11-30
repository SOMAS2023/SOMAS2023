package server

import (
	"SOMAS2023/internal/common/utils"
	"sort"

	"github.com/google/uuid"
)

func (s *Server) RunRoundLoop() {

	s.ResetGameState()
	s.FoundingInstitutions()

	// run this for n iterations
	for {
		s.RunGameLoop()
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
	for id, agent := range s.GetAgentMap() {
		// collect choice from each agent
		choice := agent.ChooseFoundingInstitution()
		s.foundingChoices[id] = choice
	}

	// tally the choices
	// FoundingChoices := make([]voting.IVoter, len(allAllocations))
	// FoundingAllocations is a map of governance method to number of agents that want that governance method
	foundingTotals := voting.tallyFoundingVotes(s.foundingChoices)

	// for each governance method, populate megabikes with the bikers who chose that governance method
	govBikes := make(map[utils.Governance][]uuid.UUID)
	bikesUsed := make([]uuid.UUID, 0)
	for governanceMethod, numBikers := range foundingTotals {
		megaBikesNeeded := math.ceil(float64(numBikers) / float64(utils.BikersOnBike))
		govBikes[governanceMethod] = make([]uuid.UUID, megaBikesNeeded)
		// get bikes for this governance
		for i := 0; i < megaBikesNeeded; i++ {
			foundBike := false
			for !foundBike {
				bike := s.GetRandomBikeId()
				if !slice.Contains(bikesUsed, bike) {
					foundBike = true
					bikesUsed = append(bikesUsed, bike)
					govBikes[governanceMethod][i] = bike
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
