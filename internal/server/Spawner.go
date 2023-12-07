package server

import (
	"SOMAS2023/internal/clients/team1"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

type AgentInitFunction func(baseBiker *objects.BaseBiker) objects.IBaseBiker

func GetAgentGenerators() []baseserver.AgentGeneratorCountPair[objects.IBaseBiker] {
	initFunctions := []AgentInitFunction{
		nil,             // Base Biker
		team1.GetBiker1, // Team 1
	}

	agentGenerators := make([]baseserver.AgentGeneratorCountPair[objects.IBaseBiker], 0, len(initFunctions))
	for _, initFunction := range initFunctions {
		agentGenerators = append(agentGenerators, baseserver.MakeAgentGeneratorCountPair(BikerAgentGenerator(initFunction), BikerAgentCount/len(initFunctions)))
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
