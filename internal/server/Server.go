package server

import (
	"SOMAS2023/internal/common/objects"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

const LootBoxCount = BikerAgentCount * 2
const MegaBikeCount = BikerAgentCount / 2

type Server struct {
	baseserver.BaseServer[objects.IBaseBiker]
	lootBoxes map[uuid.UUID]*objects.LootBox
	megaBikes map[uuid.UUID]*objects.MegaBike
	audi      objects.IAudi
}

func Initialize(iterations int) baseserver.IServer[objects.IBaseBiker] {
	server := &Server{
		BaseServer: *baseserver.CreateServer[objects.IBaseBiker](GetAgentGenerators(), iterations),
		lootBoxes:  make(map[uuid.UUID]*objects.LootBox),
		megaBikes:  make(map[uuid.UUID]*objects.MegaBike),
		audi:       objects.GetIAudi(),
	}
	server.replenishLootBoxes()
	server.replenishMegaBikes()
	return server
}

// Needed to add this getter for the collision detection for loot boxes
func (s *Server) GetMegaBikes() map[uuid.UUID]*objects.MegaBike {
	return s.megaBikes
}

func (s *Server) GetLootBoxes() map[uuid.UUID]*objects.LootBox {
	return s.lootBoxes
}
