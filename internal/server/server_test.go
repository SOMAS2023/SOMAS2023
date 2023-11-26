package server_test

import (
	"SOMAS2023/internal/server"
	"testing"
)

func TestInitServer(t *testing.T) {

	it := 3
	s := server.Initialize(it)

	if len(s.GetAgentMap()) != server.BikerAgentCount {
		t.Error("Agents not properly instantiated")
	}

	if len(s.GetMegaBikes()) != server.MegaBikeCount {
		t.Error("mega bikes not properly instantiated")
	}

	//if len(s.GetAudi()) != 1 {
	//	t.Error("Mega bikes not properly instantiated")
	//}

	s.RunGameLoop()
	s.Start()
}
