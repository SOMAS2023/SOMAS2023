package server

import (
	team_1 "SOMAS2023/internal/clients/team1"
	team_7 "SOMAS2023/internal/clients/team7/agents"
	team_8 "SOMAS2023/internal/clients/team8"
	team_3 "SOMAS2023/internal/clients/team_3"
	team_5 "SOMAS2023/internal/clients/team_5"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

func GetAgentGenerators() []baseserver.AgentGeneratorCountPair[objects.IBaseBiker] {
	return []baseserver.AgentGeneratorCountPair[objects.IBaseBiker]{
		baseserver.MakeAgentGeneratorCountPair[objects.IBaseBiker](Biker1AgentGenerator, BikerAgentCount),
		baseserver.MakeAgentGeneratorCountPair[objects.IBaseBiker](Biker5AgentGenerator, BikerAgentCount),
		// baseserver.MakeAgentGeneratorCountPair[objects.IBaseBiker](Biker7AgentGenerator, BikerAgentCount),
		baseserver.MakeAgentGeneratorCountPair[objects.IBaseBiker](Biker8AgentGenerator, BikerAgentCount),
		// baseserver.MakeAgentGeneratorCountPair[objects.IBaseBiker](Biker3AgentGenerator, BikerAgentCount),
	}
}

func Biker1AgentGenerator() objects.IBaseBiker {
	return team_1.GetBiker1(utils.GenerateRandomColour(), uuid.New())
}

func Biker3AgentGenerator() objects.IBaseBiker {
	return team_3.NewTeam3Agent(utils.GenerateRandomColour(), uuid.New())
}

func Biker5AgentGenerator() objects.IBaseBiker {
	return team_5.NewTeam5Agent(utils.GenerateRandomColour(), uuid.New())
}

func Biker7AgentGenerator() objects.IBaseBiker {
	return team_7.NewBaseTeamSevenBiker()
}

func Biker8AgentGenerator() objects.IBaseBiker {
	return team_8.GetIBaseBiker(utils.GenerateRandomColour(), uuid.New())
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
