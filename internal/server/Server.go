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

const LootBoxCount = BikerAgentCount * 2.5 // 2.5 lootboxes available per Agent
const MegaBikeCount = 11                   // Megabikes should have 8 riders
const BikerAgentCount = 56                 // 56 agents in total

type IBaseBikerServer interface {
	baseserver.IServer[objects.IBaseBiker]
	GetMegaBikes() map[uuid.UUID]objects.IMegaBike                                                               // returns all megabikes present on map
	GetLootBoxes() map[uuid.UUID]objects.ILootBox                                                                // returns all looboxes present on map
	GetAwdi() objects.IAwdi                                                                                      // returns the awdi interface
	GetJoiningRequests([]uuid.UUID) map[uuid.UUID][]uuid.UUID                                                    // returns a map from bike id to the id of all agents trying to joing that bike
	GetRandomBikeId() uuid.UUID                                                                                  // gets the id of any random bike in the map
	RulerElection(agents []objects.IBaseBiker, governance utils.Governance) uuid.UUID                            // runs the ruler election
	RunRulerAction(bike objects.IMegaBike) uuid.UUID                                                             // gets the direction from the dictator
	RunDemocraticAction(bike objects.IMegaBike, weights map[uuid.UUID]float64) uuid.UUID                         // gets the direction in voting-based governances
	NewGameStateDump(iteration int) GameStateDump                                                                // creates a new game state dump
	GetLeavingDecisions(gameState objects.IGameState) []uuid.UUID                                                // gets the list of agents that want to leave their bike
	HandleKickoutProcess() []uuid.UUID                                                                           // handles the kickout process
	ProcessJoiningRequests(inLimbo []uuid.UUID)                                                                  // processes the joining requests
	RunActionProcess()                                                                                           // runs the action (direction choice + pedalling) process for each bike
	AwdiCollisionCheck()                                                                                         // checks for collisions between awdi and bikes
	AddAgentToBike(agent objects.IBaseBiker)                                                                     // adds an agent to a bike (which also has some side effects on some server data structures)
	FoundingInstitutions()                                                                                       // runs the founding institutions process
	GetWinningDirection(finalVotes map[uuid.UUID]voting.LootboxVoteMap, weights map[uuid.UUID]float64) uuid.UUID // gets the winning direction according to the selected voting process
	LootboxCheckAndDistributions()                                                                               // checks for collision between bike and lootbox and runs the distribution process
	ResetGameState()                                                                                             // resets game state (at the beginning of a new round)
	GetDeadAgents() map[uuid.UUID]objects.IBaseBiker                                                             // returns the map of dead agents
	UpdateGameStates()                                                                                           // updates the game state object of all agents
}

type Server struct {
	baseserver.BaseServer[objects.IBaseBiker]
	lootBoxes map[uuid.UUID]objects.ILootBox
	megaBikes map[uuid.UUID]objects.IMegaBike
	// megaBikeRiders is a mapping from Agent ID -> ID of the bike that they are riding
	// helps with efficiently managing ridership status
	megaBikeRiders  map[uuid.UUID]uuid.UUID // maps riders to their bike
	awdi            objects.IAwdi
	deadAgents      map[uuid.UUID]objects.IBaseBiker // map of dead agents (used for respawning at the end of a round )
	foundingChoices map[uuid.UUID]utils.Governance
}

func Initialize(iterations int) IBaseBikerServer {
	server := &Server{
		BaseServer:     *baseserver.CreateServer[objects.IBaseBiker](GetAgentGenerators(), iterations),
		lootBoxes:      make(map[uuid.UUID]objects.ILootBox),
		megaBikes:      make(map[uuid.UUID]objects.IMegaBike),
		megaBikeRiders: make(map[uuid.UUID]uuid.UUID),
		deadAgents:     make(map[uuid.UUID]objects.IBaseBiker),
		awdi:           objects.GetIAwdi(),
	}
	server.replenishLootBoxes()
	server.replenishMegaBikes()

	return server
}

// when an agent dies it needs to be removed from its bike, the riders map and the agents map + it's added to the dead agents map
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

// ensures that adding agents to a bike is atomic (ie no agent is added to a bike while still resulting as on another bike)
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
	targetBike := agent.ChangeBike()
	if _, ok := s.megaBikes[targetBike]; !ok {
		panic("agent requested a bike that doesn't exist")
	}
	agent.SetBike(targetBike)

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

	file, err = os.Create("game_dump.json")
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
