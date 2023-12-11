package server

import (
	"SOMAS2023/internal/clients/team1"
	"SOMAS2023/internal/clients/team3"
	"SOMAS2023/internal/clients/team4"
	team5Agent "SOMAS2023/internal/clients/team5"
	"SOMAS2023/internal/clients/team8"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

type AgentInitFunction func(baseBiker *objects.BaseBiker) objects.IBaseBiker

var AgentInitFunctions = []AgentInitFunction{
	team1.GetBiker1, // Team 1
	//team2.GetBiker,      // Team 2
	team3.GetT3Agent,    // Team 3
	team4.GetBiker4,     // Team 4
	team5Agent.GetBiker, // Team 5
	team8.GetIBaseBiker, // Team 8
}

func GetAgentGenerators() []baseserver.AgentGeneratorCountPair[objects.IBaseBiker] {
	bikersPerTeam := BikerAgentCount / (len(AgentInitFunctions) + 1)
	extraBaseBikers := BikerAgentCount % (len(AgentInitFunctions) + 1)
	agentGenerators := []baseserver.AgentGeneratorCountPair[objects.IBaseBiker]{
		// Spawn base bikers
		baseserver.MakeAgentGeneratorCountPair(BikerAgentGenerator(nil), bikersPerTeam+extraBaseBikers),
	}
	for _, initFunction := range AgentInitFunctions {
		agentGenerators = append(agentGenerators, baseserver.MakeAgentGeneratorCountPair(BikerAgentGenerator(initFunction), bikersPerTeam))
	}
	return agentGenerators
}

func BikerAgentGenerator(initFunc func(baseBiker *objects.BaseBiker) objects.IBaseBiker) func() objects.IBaseBiker {
	return func() objects.IBaseBiker {
		baseBiker := objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New())
		if initFunc == nil {
			return baseBiker
		} else {
			return initFunc(baseBiker)
		}
	}
}

func (s *Server) spawnLootBox() {
	lootBox := objects.GetLootBox()
	s.lootBoxes[lootBox.GetID()] = lootBox
}

func (s *Server) replenishLootBoxes() {
	count := LootBoxCount - len(s.lootBoxes)
	for i := 0; i < count; i++ {
		s.spawnLootBox()
	}
}

func (s *Server) spawnMegaBike() {
	megaBike := objects.GetMegaBike()
	s.megaBikes[megaBike.GetID()] = megaBike
}

func (s *Server) replenishMegaBikes() {
	neededBikes := MegaBikeCount - len(s.megaBikes)
	for i := 0; i < neededBikes; i++ {
		s.spawnMegaBike()
	}
}
