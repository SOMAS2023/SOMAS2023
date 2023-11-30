package server_test

import (
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/server"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestRulerElectionDictator(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	gs := s.NewGameStateDump()
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}
	// pass gamestate
	var ruler uuid.UUID
	for _, bike := range s.GetMegaBikes() {
		agents := bike.GetAgents()
		if len(agents) != 0 {
			ruler = s.RulerElection(agents, utils.Dictatorship)
			if ruler == uuid.Nil {
				t.Error("no ruler elected")
			}
		}
	}
	// the actual logic of get winner from dist will be tested elsewhere
	fmt.Printf("\nRuler election passed \n")

}

func TestRulerElectionLeader(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	gs := s.NewGameStateDump()
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}
	// pass gamestate
	var ruler uuid.UUID
	for _, bike := range s.GetMegaBikes() {
		agents := bike.GetAgents()
		if len(agents) != 0 {
			ruler = s.RulerElection(agents, utils.Leadership)
			if ruler == uuid.Nil {
				t.Error("no ruler elected")
			}
		}
	}
	// the actual logic of get winner from dist will be tested elsewhere
	fmt.Printf("\nRuler election leader passed \n")

}

func TestRunRulerActionDictator(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	gs := s.NewGameStateDump()
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}

	for _, bike := range s.GetMegaBikes() {
		agents := bike.GetAgents()
		if len(agents) != 0 {
			// make them vote for the dictator (assume that function works properly)
			// get the dictator id (or check what it should be given the MVP strategy, this must be deterministic though)
			ruler := s.RulerElection(agents, utils.Dictatorship)
			direction := s.RunRulerAction(bike, utils.Dictatorship)
			// set the force of the dictator
			// check that the function works for it

			if bike.GetRuler() != ruler {
				t.Error("error in setting bike's ruler")
			}

			// check that the direction is one of the loots (for now)

			_, exists := s.GetLootBoxes()[direction]
			if !exists {
				t.Error("dictator returned wrong direction")
			}
		}
	}
	fmt.Printf("\nRuler action passed \n")

}

func TestRunRulerActionLeader(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	gs := s.NewGameStateDump()
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}

	for _, bike := range s.GetMegaBikes() {
		agents := bike.GetAgents()
		if len(agents) != 0 {
			// make them vote for the dictator (assume that function works properly)
			// get the dictator id (or check what it should be given the MVP strategy, this must be deterministic though)
			ruler := s.RulerElection(agents, utils.Leadership)
			direction := s.RunRulerAction(bike, utils.Leadership)
			// set the force of the dictator
			// check that the function works for it

			if bike.GetRuler() != ruler {
				t.Error("error in setting bike's ruler")
			}

			// check that the direction is one of the loots (for now)

			_, exists := s.GetLootBoxes()[direction]
			if !exists {
				t.Error("leader returned wrong direction")
			}
		}
	}
	fmt.Printf("\nRuler action  leader passed \n")
}

func TestRunDemocratingAction(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	gs := s.NewGameStateDump()
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}

	for _, bike := range s.GetMegaBikes() {
		agents := bike.GetAgents()
		if len(agents) != 0 {

			direction := s.RunDemocraticAction(bike)

			_, exists := s.GetLootBoxes()[direction]
			if !exists {
				t.Error("returned wrong direction")
			}
		}
	}
	fmt.Printf("\nDemocratic action passed \n")
}
