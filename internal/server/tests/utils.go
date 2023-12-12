package server

import (
	"SOMAS2023/internal/server"
	"testing"
)

func OnlySpawnBaseBikers(t *testing.T) {
	oldInitFunctions := server.AgentInitFunctions
	server.AgentInitFunctions = []server.AgentInitFunction{nil}
	t.Cleanup(func() {
		server.AgentInitFunctions = oldInitFunctions
	})
}
