package reputation

import (
	// Assuming this package contains the IMegaBike interface
	"SOMAS2023/internal/common/objects"
	"math"

	"github.com/google/uuid"
)

// AgentReputation defines the structure for an agent's reputation
type AgentReputation struct {
	Contribution  float64
	SurvivalScore float64
	//If colours are visible, can be used as another factor (i.e matching colour)
}

func (agentRep *AgentReputation) normaliseRep() float64 { //private
	maxContribution, maxSurvival := 100.0, 100.0
	normContribution := agentRep.Contribution / maxContribution
	normSurvival := agentRep.SurvivalScore / maxSurvival
	finalRep := (normContribution + normSurvival) / 2
	return math.Min(math.Max(finalRep, 0), 1)
}

type ReputationSystem struct {
	agentReputations map[uuid.UUID]float64
	gameState        objects.IGameState
}

// returns rep system pointer with empty map (pass as argument to rest of the functions)
func NewRepSystem(gameState objects.IGameState) *ReputationSystem {
	return &ReputationSystem{
		agentReputations: make(map[uuid.UUID]float64),
		gameState:        gameState,
	}
}

// updates the reputation of an agent
func (repSystem *ReputationSystem) UpdateAgentReputation(agentID uuid.UUID, rep AgentReputation) {
	repSystem.agentReputations[agentID] = rep.normaliseRep()
}

func (repSystem *ReputationSystem) calculateMegaBikeReputation(megaBikeID uuid.UUID) float64 {
	megaBikes := repSystem.gameState.GetMegaBikes() // Get all MegaBikes from the game state (game state not complete rn so this won't work)
	megaBike, exists := megaBikes[megaBikeID]       //exists is true if the megaBikeID is in the map
	if !exists {
		return 0
	}

	agents := megaBike.GetAgents()
	if len(agents) == 0 {
		return 0
	}

	var totalRep float64
	for _, agent := range agents { // '_' is the index
		totalRep += repSystem.agentReputations[agent.GetID()]
	}
	return math.Min(math.Max(totalRep/float64(len(agents)), 0), 1) //restricts to range [0,1]
}

// calculates the average reputation of all agents on a MegaBike
func (repSystem *ReputationSystem) CalculateAllMegaBikeReputations() map[uuid.UUID]float64 {
	megaBikes := repSystem.gameState.GetMegaBikes() // from IGameState
	megaBikeReputations := make(map[uuid.UUID]float64)

	for megaBikeID := range megaBikes {
		megaBikeReputations[megaBikeID] = repSystem.calculateMegaBikeReputation(megaBikeID)
	}
	return megaBikeReputations
}
