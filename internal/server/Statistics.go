package server

import (
	"github.com/google/uuid"
	"math"
)

type GameStatistics struct {
	AgentLifetime       map[uuid.UUID]int     `json:"agent_lifetime"`
	AgentEnergyAverage  map[uuid.UUID]float64 `json:"agent_energy_average"`
	AgentEnergyVariance map[uuid.UUID]float64 `json:"agent_energy_variance"`
	AgentPointsAverage  map[uuid.UUID]float64 `json:"agent_points_average"`
	AgentPointsVariance map[uuid.UUID]float64 `json:"agent_points_variance"`
}

func CalculateStatistics(gameStates []GameStateDump) GameStatistics {
	getAgentEnergy := func(agent *AgentDump) float64 { return agent.EnergyLevel }
	getAgentPoints := func(agent *AgentDump) float64 { return float64(agent.Points) }
	return GameStatistics{
		AgentLifetime:       agentLifetime(gameStates),
		AgentEnergyAverage:  agentAverage(gameStates, getAgentEnergy),
		AgentEnergyVariance: agentVariance(gameStates, getAgentEnergy),
		AgentPointsAverage:  agentAverage(gameStates, getAgentPoints),
		AgentPointsVariance: agentVariance(gameStates, getAgentPoints),
	}
}

func agentLifetime(gameStates []GameStateDump) map[uuid.UUID]int {
	result := make(map[uuid.UUID]int)
	for i, gameState := range gameStates {
		for id := range gameState.Agents {
			result[id] = i + 1
		}
	}
	return result
}

func agentAverage(gameStates []GameStateDump, agentProperty func(agentDump *AgentDump) float64) map[uuid.UUID]float64 {
	agentLifetime := agentLifetime(gameStates)

	result := make(map[uuid.UUID]float64)
	// result[id] := Σx
	for _, gameState := range gameStates {
		for id, agent := range gameState.Agents {
			result[id] += agentProperty(&agent)
		}
	}
	// result[id] := Σx/n == E(x)
	for id := range result {
		result[id] /= float64(agentLifetime[id])
	}
	return result
}

func agentVariance(gameStates []GameStateDump, agentProperty func(agentDump *AgentDump) float64) map[uuid.UUID]float64 {
	agentLifetime := agentLifetime(gameStates)
	agentAverage := agentAverage(gameStates, agentProperty)

	result := make(map[uuid.UUID]float64)
	// result[id] := Σ(x^2)
	for _, gameState := range gameStates {
		for id, agent := range gameState.Agents {
			result[id] += math.Pow(agentProperty(&agent), 2)
		}
	}
	for id := range result {
		// result[id] := Σ(x^2)/n == E(x^2)
		result[id] /= float64(agentLifetime[id])
		// result[id] := E(x^2) - E(x)^2 == Var(x)
		result[id] -= math.Pow(agentAverage[id], 2)
	}
	return result
}
