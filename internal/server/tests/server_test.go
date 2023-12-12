package server_test

import (
	"SOMAS2023/internal/server"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestInitialize(t *testing.T) {

	it := 3
	s := server.Initialize(it)

	if len(s.GetAgentMap()) != server.BikerAgentCount {
		t.Error("Agents not properly instantiated")
	}

	if len(s.GetMegaBikes()) != server.MegaBikeCount {
		t.Error("mega bikes not properly instantiated")
	}

	if len(s.GetLootBoxes()) != server.LootBoxCount {
		t.Error("Mega bikes not properly instantiated")
	}

	if s.GetAudi().GetID() == uuid.Nil {
		t.Error("audi not properly instantiated")
	}

	fmt.Printf("\nInitialize passed \n")
}

func TestRunGame(t *testing.T) {
	server.Initialize(1).Start()
}
