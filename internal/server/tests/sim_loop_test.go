package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/server"
	"slices"
	"testing"

	"github.com/google/uuid"
)

type MockBiker struct {
	*objects.BaseBiker
	ID             uuid.UUID
	VoteMap        map[uuid.UUID]int
	kickedOutCount int
	governance     utils.Governance
	ruler          uuid.UUID
	point          int
	BikeID         uuid.UUID
}

type MegaBike struct {
	agents         []Biker
	kickedOutCount int
}

func NewMockBiker() *MockBiker {
	baseBiker := objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New())

	return &MockBiker{
		BaseBiker: baseBiker,
		ID:        uuid.New(),
		VoteMap:   make(map[uuid.UUID]int),
	}
}

type Biker interface {
	VoteForKickout() map[uuid.UUID]int
}

func (mb *MockBiker) GetID() uuid.UUID {
	return mb.ID
}

func (mb *MockBiker) DecideGovernance() utils.Governance {
	return mb.governance
}

/*
func createMockBikers(s server.IBaseBikerServer, count int) []*MockBiker {
	var mockBikers []*MockBiker
	for i := 0; i < count; i++ {
		mockBiker := NewMockBiker()
		mockBiker.governance = utils.Democracy

		bikeID := uuid.New()
		mockBiker.BikeID = bikeID

		s.AddAgentToBike(mockBiker)

		if i%2 != 0 {
			s.RemoveAgent(mockBiker)
		}

		mockBiker.UpdateEnergyLevel(0.5)
		mockBiker.point += 10

		mockBikers = append(mockBikers, mockBiker)
	}
	return mockBikers
} */

func TestResetGameState(t *testing.T) {
	it := 2
	s := server.Initialize(it)

	mockBikers := make([]*MockBiker, 4)
	for i := range mockBikers {
		mockBiker := NewMockBiker()
		if i < 2 {
			mockBiker.governance = utils.Democracy
		} else {
			mockBiker.governance = utils.Dictatorship
		}
		s.AddAgent(mockBiker)
		mockBikers[i] = mockBiker
	}

	s.UpdateGameStates()
	s.FoundingInstitutions()

	s.UpdateGameStates()

	s.ResetGameState()

	gsNew := s.NewGameStateDump(0)
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gsNew)
	}

	for _, mockBiker := range mockBikers {
		if mockBiker.GetBikeStatus() {
			t.Errorf("Expected no bikers on bikes after reset, but found biker with ID %v on a bike", mockBiker.GetID())
		}
	}

	if utils.ResetPointsEveryRound {
		for _, agent := range s.GetAgentMap() {
			if agent.GetPoints() != 0 {
				t.Errorf("Expected agent points to be 0, got %d", agent.GetPoints())
			}
		}
	}

	if len(s.GetDeadAgents()) != 0 {
		t.Error("Expected dead agents map to have been emptied")
	}

}

func TestFoundingInstitutions(t *testing.T) {
	OnlySpawnBaseBikers(t)

	it := 2
	s := server.Initialize(it)

	mockBikers := make([]*MockBiker, 4)
	for i := range mockBikers {
		mockBiker := NewMockBiker()
		if i < 2 {
			mockBiker.governance = utils.Democracy
		} else {
			mockBiker.governance = utils.Dictatorship
		}
		s.AddAgent(mockBiker)
		mockBikers[i] = mockBiker
	}
	gs := s.NewGameStateDump(0)
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}
	s.FoundingInstitutions()

	/* 	for _, agent := range s.GetAgentMap() {
		bikeID := agent.GetBike()
		bike := s.GetMegaBikes()[bikeID]
		if bike != nil && bike.GetGovernance() != agent.DecideGovernance() {
			t.Errorf("Agent %v is on bike with governance %v, want %v",
				agent.GetID(), bike.GetGovernance(), agent.DecideGovernance())
		}
	} */

	for _, biker := range mockBikers {
		actualBike := s.GetMegaBikes()[biker.GetBike()]
		if actualBike == nil {
			t.Errorf("Biker %v has not been assigned to any bike", biker.GetID())
			continue
		}
		if actualBike.GetGovernance() != biker.governance {
			t.Errorf("Biker %v is on a bike with incorrect governance", biker.GetID())
		}
	}

	for _, agent := range s.GetAgentMap() {
		bikeID := agent.GetBike()
		if bikeID == uuid.Nil {
			t.Errorf("Agent %v has not been assigned to any bike", agent.GetID())
		}
		if bike, ok := s.GetMegaBikes()[bikeID]; ok {
			if bike.GetGovernance() != agent.DecideGovernance() {
				t.Errorf("Agent %v is on bike with governance %v, want %v",
					agent.GetID(), bike.GetGovernance(), agent.DecideGovernance())
			}
			agents := bike.GetAgents()
			if !slices.Contains(agents, agent) {
				t.Errorf("Agent %v is not on the bike they were assigned to", agent.GetID())
			}
		} else {
			t.Errorf("Agent %v has not been assigned to any bike", agent.GetID())
		}
	}
}
