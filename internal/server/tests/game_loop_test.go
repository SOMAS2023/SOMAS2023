package server_test

import (
	obj "SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/server"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
)

func TestGetLeavingDecisions(t *testing.T) {
	// check that if some biker has on bike set to false they are not on any megabike
	// nor in the megabike riders
	it := 3
	s := server.Initialize(it)
	gs := s.NewGameStateDump()
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}

	s.GetLeavingDecisions(gs)

	for _, agent := range s.GetAgentMap() {
		if !agent.GetBikeStatus() {
			for _, bike := range s.GetMegaBikes() {
				for _, agentOnBike := range bike.GetAgents() {
					if agentOnBike.GetID() == agent.GetID() {
						t.Error("leaving agent is on a bike when it shouldn't be")

					}
				}
			}
		}
	}
	fmt.Printf("\nGet leaving decisions passed \n")
}

func TestHandleKickout(t *testing.T) {
	it := 6
	s := server.Initialize(it)
	gs := s.NewGameStateDump()
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}

	s.HandleKickoutProcess()

	for _, agent := range s.GetAgentMap() {
		if !agent.GetBikeStatus() {
			for _, bike := range s.GetMegaBikes() {
				for _, agentOnBike := range bike.GetAgents() {
					if agentOnBike.GetID() == agent.GetID() {
						t.Error("leaving agent is on a bike when it shouldn't be")

					}
				}
			}
		}
	}
	fmt.Printf("\nHadle kickout passed \n")
}

func TestProcessJoiningRequests(t *testing.T) {
	it := 3
	s := server.Initialize(it)

	// 1: get two bike ids
	targetBikes := make([]uuid.UUID, 2)

	i := 0
	for bikeId, _ := range s.GetMegaBikes() {
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

	// all agents should be accepted as there should be enough room on all bikes (but make it subject to that)
	// check that all of them are now on bikes
	// check that there are no bikers left with on bike = false

	s.ProcessJoiningRequests(make([]uuid.UUID, 0))
	for bikeID, agents := range requests {
		bike := s.GetMegaBikes()[bikeID]
		for _, agent := range agents {
			onBike := false
			for _, agentOnBike := range bike.GetAgents() {
				onBikeId := agentOnBike.GetID()
				if onBikeId == agent {
					onBike = true
					if !agentOnBike.GetBikeStatus() {
						t.Error("biker's status wasn't successfully toggled back")
					}
					break
				}
			}
			if !onBike {
				t.Error("biker wasn't successfully accepted on bike")
			}
		}
	}
	fmt.Printf("\nProcess joining request passed \n")
}

func TestRunActionProcess(t *testing.T) {
	for i := 0; i < 10; i++ {
		it := 1
		s := server.Initialize(it)
		gs := s.NewGameStateDump()

		// Loop through each bike
		for _, bike := range s.GetMegaBikes() {
			// Randomly select a governance strategy for this bike
			governanceTypes := []int{int(utils.Democracy), int(utils.Leadership), int(utils.Dictatorship)}
			governance := utils.Governance(governanceTypes[rand.Intn(len(governanceTypes))])
			bike.SetGovernance(governance)

			// Update the game state for all agents and set the governance of their bike
			for _, agent := range s.GetAgentMap() {
				agent.UpdateGameState(gs)
			}

			// Randomly select a ruler if necessary
			if governance != utils.Democracy {
				agents := bike.GetAgents()
				if len(agents) > 0 {
					randIndex := rand.Intn(len(agents))
					randomAgent := agents[randIndex]
					bike.SetRuler(randomAgent.GetID())
				}
			}
		}

		s.RunActionProcess()
		// check all agents have lost energy (proportionally to how much they have pedalled)
		for _, agent := range s.GetAgentMap() {
			lostEnergy := (utils.MovingDepletion * agent.GetForces().Pedal)

			var agentBike obj.IMegaBike
			for _, bike := range s.GetMegaBikes() {
				if bike.GetID() == agent.GetBike() {
					agentBike = bike
				}
			}

			governance := agentBike.GetGovernance()
			switch governance {
			case utils.Democracy:
				lostEnergy += utils.DeliberativeDemocracyPenalty
			case utils.Leadership:
				lostEnergy += utils.LeadershipDemocracyPenalty
			}
			if agent.GetEnergyLevel() != (1.0 - lostEnergy) {
				t.Error("agents energy hasn't been successfully depleted")
			}
		}
	}
	fmt.Printf("\nRun action process passed \n")
}
