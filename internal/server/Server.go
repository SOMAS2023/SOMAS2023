package server

import (
	"SOMAS2023/internal/common/objects"
	"encoding/json"
	"fmt"
	"os"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

const LootBoxCount = BikerAgentCount * 2
const MegaBikeCount = BikerAgentCount / 2

type IBaseBikerServer interface {
	baseserver.IServer[objects.IBaseBiker]
}

type Server struct {
	baseserver.BaseServer[objects.IBaseBiker]
	lootBoxes map[uuid.UUID]objects.ILootBox
	megaBikes map[uuid.UUID]objects.IMegaBike
	// megaBikeRiders is a mapping from Agent ID -> ID of the bike that they are riding
	// helps with efficiently managing ridership status
	megaBikeRiders map[uuid.UUID]uuid.UUID
	audi           objects.IAudi
}

func Initialize(iterations int) IBaseBikerServer {
	server := &Server{
		BaseServer:     *baseserver.CreateServer[objects.IBaseBiker](GetAgentGenerators(), iterations),
		lootBoxes:      make(map[uuid.UUID]objects.ILootBox),
		megaBikes:      make(map[uuid.UUID]objects.IMegaBike),
		megaBikeRiders: make(map[uuid.UUID]uuid.UUID),
		audi:           objects.GetIAudi(),
	}
	server.replenishLootBoxes()
	server.replenishMegaBikes()

	// Randomly allocate bikers to bikes
	for _, biker := range server.GetAgentMap() {
		server.SetBikerBike(biker, server.GetRandomBikeId())
	}

	return server
}

func (s *Server) RemoveAgent(agent objects.IBaseBiker) {
	id := agent.GetID()
	s.BaseServer.RemoveAgent(agent)
	if bikeId, ok := s.megaBikeRiders[id]; ok {
		s.megaBikes[bikeId].RemoveAgent(id)
		delete(s.megaBikeRiders, id)
	}
}

func (s *Server) outputResults(gameStates []GameStateDump) {
	statisticsJson, err := json.MarshalIndent(CalculateStatistics(gameStates), "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println("Statistics:\n" + string(statisticsJson))

	file, err := os.Create("game_dump.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(gameStates); err != nil {
		panic(err)
	}
}
