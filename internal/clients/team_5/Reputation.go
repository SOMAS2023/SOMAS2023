package team5Agent

import (
	// Assuming this package contains the IMegaBike interface
	"SOMAS2023/internal/common/objects"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type ReputationSystem struct {
	agentReputations map[uuid.UUID]float64
	gameState        objects.IGameState
}

// Creates and initialises map of reputations to 0
func NewRepSystem(gameState objects.IGameState) *ReputationSystem {
	agentReputations := make(map[uuid.UUID]float64)
	megaBikes := gameState.GetMegaBikes()
	for _, mb := range megaBikes {
		// Iterate through all agents on each MegaBike
		for _, agent := range mb.GetAgents() {
			// Set initial reputation to 0.5 for each agent
			agentReputations[agent.GetID()] = 0.5
		}
	}

	return &ReputationSystem{
		agentReputations: agentReputations,
		gameState:        gameState,
	}
}

func (repSystem *ReputationSystem) GetAgentReputation(agentID uuid.UUID) (float64, error) {
	rep, exists := repSystem.agentReputations[agentID]
	if !exists {
		return 0, fmt.Errorf("agent with UUID %s not found", agentID)
	}
	return rep, nil
}

// Most important 3 functions:

// Reputation calculation currently just based on energy and force
func (repSystem *ReputationSystem) calculateReputationOfAgent(agentID uuid.UUID) float64 {
	averagePedalForce := repSystem.getAverageForceOfAgents()
	averageEnergy := repSystem.getAverageEnergyOfAgents()

	agentPedalForce := repSystem.getForceOfOneAgent(agentID)
	agentEnergy := repSystem.getEnergyOfOneAgent(agentID)

	forceDeviation := agentPedalForce / averagePedalForce //fraction of agentMetric/averageMetric
	energyDeviation := agentEnergy / averageEnergy

	combinedDeviation := (forceDeviation + energyDeviation) / 2 // keeps it in range [0,1]

	// get current reputation of the agent
	currentRep, exists := repSystem.agentReputations[agentID]
	if !exists {
		currentRep = 0.5 // Default to 0.5 if not found
	}

	weight := 0.2 //maximum change per round
	newRep := currentRep + (combinedDeviation-1)*weight
	return math.Min(math.Max(newRep, 0), 1) //capped at 0 and 1
}

func (repSystem *ReputationSystem) updateReputationOfAllAgents() {
	for agentID := range repSystem.agentReputations {
		newRep := repSystem.calculateReputationOfAgent(agentID)
		repSystem.agentReputations[agentID] = newRep
	}
}

func (repSystem *ReputationSystem) updateGameState(gameState objects.IGameState) {
	repSystem.gameState = gameState
}

//Useful helper functions:

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
	for _, agent := range agents { // _ is index
		totalRep += repSystem.agentReputations[agent.GetID()]
	}
	return math.Min(math.Max(totalRep/float64(len(agents)), 0), 1) //restricts to range [0,1]
}

func (repSystem *ReputationSystem) getAveragePedalSpeedOfMegaBike(megaBikeID uuid.UUID) float64 {
	megaBikes := repSystem.gameState.GetMegaBikes()
	megaBike, exists := megaBikes[megaBikeID]
	if !exists {
		return 0
	}
	agents := megaBike.GetAgents()
	var totalPedalSpeed float64
	for _, agent := range agents {
		totalPedalSpeed += agent.GetForces().Pedal
	}
	return totalPedalSpeed / float64(len(agents))
}

// Functions used in calculating the reputation value:
func (repSystem *ReputationSystem) getAverageEnergyOfAgents() float64 {
	megaBikes := repSystem.gameState.GetMegaBikes()
	var totalEnergy float64
	var totalAgents float64
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			totalEnergy += agent.GetEnergyLevel()
			totalAgents++
		}
	}
	return totalEnergy / totalAgents
}

func (repSystem *ReputationSystem) getAverageForceOfAgents() float64 {
	megaBikes := repSystem.gameState.GetMegaBikes()
	var totalForce float64
	var totalAgents float64
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			forceOfAgent := agent.GetForces().Pedal
			if forceOfAgent > 0 { //only add force if agent is pedalling
				totalForce += forceOfAgent
				totalAgents++
			}
		}
	}
	return totalForce / totalAgents
}

func (repSystem *ReputationSystem) getEnergyOfOneAgent(agentID uuid.UUID) float64 {
	megaBikes := repSystem.gameState.GetMegaBikes()
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			if agent.GetID() == agentID {
				return agent.GetEnergyLevel()
			}
		}
	}
	return 0
}

func (repSystem *ReputationSystem) getForceOfOneAgent(agentID uuid.UUID) float64 {
	megaBikes := repSystem.gameState.GetMegaBikes()
	for _, megaBike := range megaBikes {
		agents := megaBike.GetAgents()
		for _, agent := range agents {
			if agent.GetID() == agentID {
				return agent.GetForces().Pedal
			}
		}
	}
	return 0
}
