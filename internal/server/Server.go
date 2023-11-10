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
	// megaBikeRiders is a mapping from Agent ID -> ID of the bike that they are riding
	// helps with efficiently managing ridership status
	megaBikeRiders map[uuid.UUID]uuid.UUID
}

func Initialize(iterations int) baseserver.IServer[objects.IBaseBiker] {
	server := &Server{
		BaseServer: *baseserver.CreateServer[objects.IBaseBiker](GetAgentGenerators(), iterations),
		lootBoxes:  make(map[uuid.UUID]*objects.LootBox),
		megaBikes:  make(map[uuid.UUID]*objects.MegaBike),
	}
	server.replenishLootBoxes()
	server.replenishMegaBikes()
	return server
}
