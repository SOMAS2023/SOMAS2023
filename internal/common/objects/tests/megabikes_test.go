package objects

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/server"
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

func (mb *MockBiker) VoteForKickout() map[uuid.UUID]int {
	return mb.VoteMap
}

type Biker interface {
	VoteForKickout() map[uuid.UUID]int
}

// Ensure that BaseBiker implements the Biker interface.
//var _ Biker = &objects.BaseBiker{}

func (mb *MockBiker) GetID() uuid.UUID {
	return mb.ID
}

func TestGetMegaBike(t *testing.T) {
	mb := objects.GetMegaBike()

	if mb == nil {
		t.Errorf("GetMegaBike returned nil")
	}

	if mb.GetGovernance() != utils.Democracy {
		t.Errorf("Expected governance to be Democracy, got %v", mb.GetGovernance())
	}

	if mb.GetRuler() != uuid.Nil {
		t.Errorf("Expected ruler to be uuid.Nil, got %v", mb.GetRuler())
	}
}

func TestAddAgent(t *testing.T) {
	mb := objects.GetMegaBike()
	biker := NewMockBiker()

	mb.AddAgent(biker)

	if len(mb.GetAgents()) != 1 {
		t.Errorf("AddAgent failed to add the agent to MegaBike")
	}

	if mb.GetAgents()[0].GetID() != biker.GetID() {
		t.Errorf("The added agent ID does not match the expected MockBiker ID")
	}
}

func TestRemoveAgent(t *testing.T) {
	mb := objects.GetMegaBike()
	biker1 := NewMockBiker()
	biker2 := NewMockBiker()

	mb.AddAgent(biker1)
	mb.AddAgent(biker2)

	mb.RemoveAgent(biker1.GetID())

	agents := mb.GetAgents()

	if len(agents) != 1 {
		t.Errorf("RemoveAgent failed to remove the agent from MegaBike, expected 1 agent, got %d", len(agents))
	}

	if agents[0].GetID() == biker1.GetID() {
		t.Errorf("RemoveAgent did not remove the correct agent")
	}

	if agents[0].GetID() != biker2.GetID() {
		t.Errorf("The remaining agent ID does not match the expected MockBiker ID")
	}
}

func TestUpdateMass(t *testing.T) {
	mb := objects.GetMegaBike()
	initialMass := mb.GetPhysicalState().Mass

	mb.AddAgent(NewMockBiker())
	mb.AddAgent(NewMockBiker())
	mb.UpdateMass()

	updatedMass := mb.GetPhysicalState().Mass

	expectedMass := initialMass + 2

	if updatedMass != expectedMass {
		t.Errorf("UpdateMass did not calculate the correct mass: got %v, want %v", updatedMass, expectedMass)
	}
}

func TestGetSetGovernanceAndRuler(t *testing.T) {
	mb := objects.GetMegaBike()
	originalGovernance := mb.GetGovernance()
	originalRuler := mb.GetRuler()

	newGovernance := utils.Dictatorship
	newRuler := uuid.New()

	mb.SetGovernance(newGovernance)
	mb.SetRuler(newRuler)

	if mb.GetGovernance() != newGovernance {
		t.Errorf("SetGovernance failed, expected %v, got %v", newGovernance, mb.GetGovernance())
	}

	if mb.GetRuler() != newRuler {
		t.Errorf("SetRuler failed, expected %v, got %v", newRuler, mb.GetRuler())
	}

	mb.SetGovernance(originalGovernance)
	mb.SetRuler(originalRuler)
}

func TestKickOutAgent(t *testing.T) {
	it := 3
	s := server.Initialize(it)
	s.FoundingInstitutions()
	gs := s.NewGameStateDump()

	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}

	mb := objects.GetMegaBike()

	//biker1 := NewMockBiker(uuid.New(), map[uuid.UUID]int{ /* votes */ })
	biker1 := NewMockBiker()
	biker2 := NewMockBiker()
	biker3 := NewMockBiker()
	mb.AddAgent(biker1)
	mb.AddAgent(biker2)
	mb.AddAgent(biker3)

	weights := map[uuid.UUID]float64{
		biker1.GetID(): 1.0,
		biker2.GetID(): 1.0,
		biker3.GetID(): 1.0,
	}

	// Voting
	biker1.VoteMap[biker3.GetID()] = 1
	biker2.VoteMap[biker3.GetID()] = 1
	biker3.VoteMap[biker1.GetID()] = 1

	for _, biker := range []Biker{biker1, biker2, biker3} {
		biker.VoteForKickout()
	}

	// Kick out agents based on votes and weights.
	kickedOutAgents := mb.KickOutAgent(weights)

	if len(kickedOutAgents) != 1 {
		t.Fatalf("KickOutAgent kicked out %d agents; want 1", len(kickedOutAgents))
	}

	if kickedOutAgents[0] != biker3.GetID() {
		t.Errorf("KickOutAgent kicked out incorrect agent: got %v, want %v", kickedOutAgents[0], biker3.GetID())
	}

	for _, anyBike := range s.GetMegaBikes() {
		agentsOnBike := anyBike.GetAgents()
		// Skip empty bikes.
		if agentsOnBike == nil {
			continue
		}
		// Check if biker3 is still on a bike.
		for _, agentOnBike := range agentsOnBike {
			if agentOnBike.GetID() == biker3.GetID() {
				t.Errorf("Kicked out agent is still present on a MegaBike")
			}
		}
	}
}
