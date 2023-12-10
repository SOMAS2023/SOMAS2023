package server

import (
	"math"

	"github.com/google/uuid"
	"github.com/tealeg/xlsx/v3"
)

type GameStatistics struct {
	PerRound []AgentStatistics `json:"per_round"`
	Average  AgentStatistics   `json:"average"`
}

type AgentStatistics struct {
	AgentLifetime       map[uuid.UUID]float64 `json:"agent_lifetime"`
	AgentEnergyAverage  map[uuid.UUID]float64 `json:"agent_energy_average"`
	AgentEnergyVariance map[uuid.UUID]float64 `json:"agent_energy_variance"`
	AgentPointsAverage  map[uuid.UUID]float64 `json:"agent_points_average"`
	AgentPointsVariance map[uuid.UUID]float64 `json:"agent_points_variance"`
}

type AgentStatisticAccessor func(statistics *AgentStatistics) map[uuid.UUID]float64

var (
	getLifetime       = func(statistics *AgentStatistics) map[uuid.UUID]float64 { return statistics.AgentLifetime }
	getEnergyAverage  = func(statistics *AgentStatistics) map[uuid.UUID]float64 { return statistics.AgentEnergyAverage }
	getEnergyVariance = func(statistics *AgentStatistics) map[uuid.UUID]float64 { return statistics.AgentEnergyVariance }
	getPointsAverage  = func(statistics *AgentStatistics) map[uuid.UUID]float64 { return statistics.AgentPointsAverage }
	getPointsVariance = func(statistics *AgentStatistics) map[uuid.UUID]float64 { return statistics.AgentPointsVariance }
)

func averageStatisticsOverRounds(statisticsPerRound []AgentStatistics, accessor AgentStatisticAccessor) map[uuid.UUID]float64 {
	result := make(map[uuid.UUID]float64)
	roundsAlive := make(map[uuid.UUID]int)
	for _, round := range statisticsPerRound {
		statistics := accessor(&round)
		for agentID, value := range statistics {
			result[agentID] += value
			roundsAlive[agentID]++
		}
	}

	for agentID, n := range roundsAlive {
		result[agentID] /= float64(n)
	}

	return result
}

func CalculateStatistics(gameStates [][]GameStateDump) GameStatistics {
	getAgentEnergy := func(agent *AgentDump) float64 { return agent.EnergyLevel }
	getAgentPoints := func(agent *AgentDump) float64 { return float64(agent.Points) }

	statisticsPerRound := make([]AgentStatistics, 0, len(gameStates))
	for _, round := range gameStates {
		statisticsPerRound = append(statisticsPerRound, AgentStatistics{
			AgentLifetime:       agentLifetime(round),
			AgentEnergyAverage:  agentAverage(round, getAgentEnergy),
			AgentEnergyVariance: agentVariance(round, getAgentEnergy),
			AgentPointsAverage:  agentAverage(round, getAgentPoints),
			AgentPointsVariance: agentVariance(round, getAgentPoints),
		})
	}

	return GameStatistics{
		PerRound: statisticsPerRound,
		Average: AgentStatistics{
			AgentLifetime:       averageStatisticsOverRounds(statisticsPerRound, getLifetime),
			AgentEnergyAverage:  averageStatisticsOverRounds(statisticsPerRound, getEnergyAverage),
			AgentEnergyVariance: averageStatisticsOverRounds(statisticsPerRound, getEnergyVariance),
			AgentPointsAverage:  averageStatisticsOverRounds(statisticsPerRound, getPointsAverage),
			AgentPointsVariance: averageStatisticsOverRounds(statisticsPerRound, getPointsVariance),
		},
	}
}

func agentLifetime(gameStates []GameStateDump) map[uuid.UUID]float64 {
	result := make(map[uuid.UUID]float64)
	for i, gameState := range gameStates {
		for id := range gameState.Agents {
			result[id] = float64(i)
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
		result[id] /= agentLifetime[id] + 1
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
		result[id] /= agentLifetime[id] + 1
		// result[id] := E(x^2) - E(x)^2 == Var(x)
		result[id] -= math.Pow(agentAverage[id]+1, 2)
	}
	return result
}

func (gs *GameStatistics) ToSpreadsheet() *xlsx.File {
	workbook := xlsx.NewFile()

	columnIndexes := make(map[uuid.UUID]int)
	nextIndex := 1
	getColumnIndex := func(agentID uuid.UUID) int {
		if idx, ok := columnIndexes[agentID]; ok {
			return idx
		} else {
			idx = nextIndex
			nextIndex++
			columnIndexes[agentID] = idx
			return idx
		}
	}

	writeSheet := func(sheetName string, accessor AgentStatisticAccessor) {
		sheet, err := workbook.AddSheet(sheetName)
		if err != nil {
			panic(err)
		}

		headerRow := sheet.AddRow()
		headerRow.GetCell(0).SetString("Round")

		for i, round := range gs.PerRound {
			row := sheet.AddRow()
			row.GetCell(0).SetValue(i + 1)
			for id, value := range accessor(&round) {
				columnIndex := getColumnIndex(id)
				headerRow.GetCell(columnIndex).SetString(id.String())
				row.GetCell(columnIndex).SetValue(value)
			}
		}
	}

	writeSheet("Lifetime", getLifetime)
	writeSheet("Energy Average", getEnergyAverage)
	writeSheet("Energy Variance", getEnergyVariance)
	writeSheet("Points Average", getPointsAverage)
	writeSheet("Points Variance", getPointsVariance)

	return workbook
}
