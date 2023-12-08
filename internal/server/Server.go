package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"encoding/json"
	"fmt"
	"os"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

const LootBoxCount = BikerAgentCount * 2
const MegaBikeCount = BikerAgentCount / 2
const BikerAgentCount = 20

type IBaseBikerServer interface {
	baseserver.IServer[objects.IBaseBiker]
	GetMegaBikes() map[uuid.UUID]objects.IMegaBike
	GetLootBoxes() map[uuid.UUID]objects.ILootBox
	GetAudi() objects.IAudi
	GetJoiningRequests([]uuid.UUID) map[uuid.UUID][]uuid.UUID
	GetRandomBikeId() uuid.UUID
	RulerElection(agents []objects.IBaseBiker, governance utils.Governance) uuid.UUID
	RunRulerAction(bike objects.IMegaBike) uuid.UUID
	RunDemocraticAction(bike objects.IMegaBike, weights map[uuid.UUID]float64) uuid.UUID
	NewGameStateDump(iteration int) GameStateDump
	GetLeavingDecisions(gameState objects.IGameState) []uuid.UUID
	HandleKickoutProcess() []uuid.UUID
	ProcessJoiningRequests(inLimbo []uuid.UUID)
	RunActionProcess()
	AudiCollisionCheck()
	AddAgentToBike(agent objects.IBaseBiker)
	FoundingInstitutions()
	GetWinningDirection(finalVotes map[uuid.UUID]voting.LootboxVoteMap, weights map[uuid.UUID]float64) uuid.UUID
	LootboxCheckAndDistributions()
	ResetGameState()
	GetDeadAgents() map[uuid.UUID]objects.IBaseBiker
	UpdateGameStates()
}

type Server struct {
	baseserver.BaseServer[objects.IBaseBiker]
	lootBoxes map[uuid.UUID]objects.ILootBox
	megaBikes map[uuid.UUID]objects.IMegaBike
	// megaBikeRiders is a mapping from Agent ID -> ID of the bike that they are riding
	// helps with efficiently managing ridership status
	megaBikeRiders  map[uuid.UUID]uuid.UUID
	audi            objects.IAudi
	deadAgents      map[uuid.UUID]objects.IBaseBiker
	foundingChoices map[uuid.UUID]utils.Governance
}

func Initialize(iterations int) IBaseBikerServer {
	server := &Server{
		BaseServer:     *baseserver.CreateServer[objects.IBaseBiker](GetAgentGenerators(), iterations),
		lootBoxes:      make(map[uuid.UUID]objects.ILootBox),
		megaBikes:      make(map[uuid.UUID]objects.IMegaBike),
		megaBikeRiders: make(map[uuid.UUID]uuid.UUID),
		deadAgents:     make(map[uuid.UUID]objects.IBaseBiker),
		audi:           objects.GetIAudi(),
	}
	server.replenishLootBoxes()
	server.replenishMegaBikes()

	return server
}

func (s *Server) RemoveAgent(agent objects.IBaseBiker) {
	id := agent.GetID()
	// add agent to dead agent map
	s.deadAgents[id] = agent
	// remove agent from agent map
	s.BaseServer.RemoveAgent(agent)
	if bikeId, ok := s.megaBikeRiders[id]; ok {
		s.megaBikes[bikeId].RemoveAgent(id)
		delete(s.megaBikeRiders, id)
	}
}

func (s *Server) AddAgentToBike(agent objects.IBaseBiker) {
	// Remove the agent from the old bike, if it was on one
	if oldBikeId, ok := s.megaBikeRiders[agent.GetID()]; ok {
		s.megaBikes[oldBikeId].RemoveAgent(agent.GetID())
	}

	// set agent on desired bike
	bikeId := agent.GetBike()
	s.megaBikes[bikeId].AddAgent(agent)
	s.megaBikeRiders[agent.GetID()] = bikeId
	if !agent.GetBikeStatus() {
		agent.ToggleOnBike()
	}
}

func (s *Server) RemoveAgentFromBike(agent objects.IBaseBiker) {
	bike := s.megaBikes[agent.GetBike()]
	bike.RemoveAgent(agent.GetID())
	agent.ToggleOnBike()
	// get new destination for agent
	agent.SetBike(agent.ChangeBike())

	if _, ok := s.megaBikeRiders[agent.GetID()]; ok {
		delete(s.megaBikeRiders, agent.GetID())
	}
}

func (s *Server) GetDeadAgents() map[uuid.UUID]objects.IBaseBiker {
	return s.deadAgents
}

func (s *Server) outputResults(gameStates [][]GameStateDump) {
	statistics := CalculateStatistics(gameStates)

	statisticsJson, err := json.MarshalIndent(statistics.Average, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println("Average Statistics:\n" + string(statisticsJson))

	file, err := os.Create("statistics.xlsx")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := statistics.ToSpreadsheet().Write(file); err != nil {
		panic(err)
	}

	var flattenedGameStates []GameStateDump
	for i := range gameStates {
		flattenedGameStates = append(flattenedGameStates, gameStates[i]...)
	}

	file, err = os.Create("game_dump.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(flattenedGameStates); err != nil {
		panic(err)
	}
}

func (s *Server) UpdateGameStates() {
	gs := s.NewGameStateDump(0)
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}
}

// had to override to address the fact that agents only have access to the game dump
// version of agents, so if the recipients are set to be those it will panic as they
// can't call the handler functions
func (s *Server) RunMessagingSession() {
	agentArray := s.GenerateAgentArrayFromMap()

	for _, agent := range s.GetAgentMap() {
		allMessages := agent.GetAllMessages(agentArray)
		for _, msg := range allMessages {
			recipients := msg.GetRecipients()
			// make recipient list with actual agents
			usableRecipients := make([]objects.IBaseBiker, len(recipients))
			for i, recipient := range recipients {
				usableRecipients[i] = s.GetAgentMap()[recipient.GetID()]
			}
			for _, recip := range usableRecipients {
				if agent.GetID() == recip.GetID() {
					continue
				}
				msg.InvokeMessageHandler(recip)
			}
		}
	}
}
