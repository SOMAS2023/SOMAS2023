package main

import (
    "math"
    "SOMAS2023/internal/common/utils"
    "SOMAS2023/internal/objects"
    "github.com/google/uuid"
)

// calculates final preferences for loot boxes based on various factors
func CalculateLootBoxPreferences(gameState objects.IGameState, agentID uuid.UUID, cumulativePreferences map[uuid.UUID]float64) map[uuid.UUID]float64 {
    finalPreferences := make(map[uuid.UUID]float64)

    // retrieve agent and loot boxes from game state
    agent := gameState.GetMegaBikes()[agentID]
    lootBoxes := gameState.GetLootBoxes()

    for id, lootBox := range lootBoxes {
        distance := calculateDistance(agent.GetPosition(), lootBox.GetPosition())
        colorPreference := calculateColorPreference(agent.GetColour(), lootBox.GetColour())
        energyPreference := calculateEnergyPreference(agent.GetEnergyLevel(), lootBox.GetTotalResources())
        cumulativePreference := cumulativePreferences[id]

        // combine preferences (weights: 0.4 for distance, 0.3 for color, 0.2 for energy, 0.1 for cumulative)
        // ensure that if cant get first preference, get second preference and so on 
         
        finalPreferences[id] = 0.4*(1/distance) + 0.3*colorPreference + 0.2*energyPreference + 0.1*cumulativePreference
    }

    return finalPreferences
}

// ensure that if cant get first preference, get second preference and so on

// calculates the Euclidean distance between two points
func calculateDistance(a, b utils.Coordinates) float64 {
    return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// calculates preference based on color match
func calculateColorPreference(agentColor, lootBoxColor utils.Colour) float64 {
    if agentColor == lootBoxColor {
        return 1.0
    }
    return 0.0
}

// calculates preference based on energy level and loot resources
// example: higher preference if energy level is low and loot is high
func calculateEnergyPreference(agentEnergy, lootResources float64) float64 {
    return lootResources * (1 - agentEnergy)
}

