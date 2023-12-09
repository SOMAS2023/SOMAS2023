package server

import (

	// "SOMAS2023/internal/clients/team_3"
	"SOMAS2023/internal/clients/team1"
	"SOMAS2023/internal/clients/team2"
	"SOMAS2023/internal/clients/team4"
	"SOMAS2023/internal/clients/team8"
	team5Agent "SOMAS2023/internal/clients/team_5"

	// "SOMAS2023/internal/clients/team7/agents"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	//team7agents "SOMAS2023/internal/clients/team7/agents"
	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

type AgentInitFunction func(baseBiker *objects.BaseBiker) objects.IBaseBiker

var AgentInitFunctions = []AgentInitFunction{
	//nil,                 // Base Biker
	team1.GetBiker1, // Team 1 works
	team2.GetBiker,  // Team 2 works
	// team_3.GetT3Agent, // Team 3
	team4.GetBiker4,      // Team 4
	team5Agent.GetBiker5, // Team 5 works?
	// agents.GetBiker7,     // Team 7
	team8.GetBiker8, // Team 8 works

}

func GetAgentGenerators() []baseserver.AgentGeneratorCountPair[objects.IBaseBiker] {
	agentGenerators := make([]baseserver.AgentGeneratorCountPair[objects.IBaseBiker], 0, len(AgentInitFunctions))
	for _, initFunction := range AgentInitFunctions {
		agentGenerators = append(agentGenerators, baseserver.MakeAgentGeneratorCountPair(BikerAgentGenerator(initFunction), BikerAgentCount/len(AgentInitFunctions)))
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
