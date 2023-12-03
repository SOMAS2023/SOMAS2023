package objects

import (
	"testing"
	"github.com/google/uuid"
	"SOMAS2023/internal/common/utils"
)

type MockBiker struct {
	id uuid.UUID
	forces utils.Forces
}

func NewMockBiker() *MockBiker {
	return &MockBiker{
		id: uuid.New(),
		forces: utils.Forces{},
	}
}

func (mb *MockBiker) GetID() uuid.UUID {
	return mb.id
}

func (mb *MockBiker) GetForces() utils.Forces {
	return mb.forces
}


func TestAddRemoveAgent(t *testing.T) {
	mb := GetMegaBike()
	biker := NewMockBiker()
	mb.AddAgent(biker)
	if len(mb.GetAgents()) != 1 {
		t.Errorf("AddAgent failed, expected 1 agent, got %d", len(mb.GetAgents()))
	}

	mb.RemoveAgent(biker.GetID())
	if len(mb.GetAgents()) != 0 {
		t.Errorf("RemoveAgent failed, expected 0 agents, got %d", len(mb.GetAgents()))
	}
}