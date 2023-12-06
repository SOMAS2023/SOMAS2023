package team5Agent

import (
	"SOMAS2023/internal/common/utils"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

// if t5.prev == nul{
// 	make map
// 	make diff map
// 	t5.prevEnergy = currMap
// }
// else{
// 	current - agent.prev > 0
// 	then update prev
// }

func (t5 *team5Agent) Something(BikeId uuid.UUID) map[uuid.UUID]float64 {
	if t5.prevEnergy == nil {
		t5.prevEnergy = make(map[uuid.UUID]float64)
	}

	bike := t5.GetGameState().GetMegaBikes()[BikeId]
	agentsOnBike := bike.GetAgents()
	energyChange := make(map[uuid.UUID]float64)

	for i, agent := range agentsOnBike {
		fmt.Printf("Agent at index %d: %v\n", i, agent)
		previousEnergy := t5.prevEnergy[agent.GetID()]

		if previousEnergy <= agent.GetEnergyLevel() {
			energyChange[agent.GetID()] = agent.GetEnergyLevel() - previousEnergy
		}

		t5.prevEnergy[agent.GetID()] = agent.GetEnergyLevel()
	}

	return energyChange

}

func (t5 *team5Agent) VoteForKickout() map[uuid.UUID]int {
	agentsOnBike := t5.GetFellowBikers()
	numberOfAgents := float64(len(agentsOnBike))

	internalRanking := make(map[uuid.UUID]float64)
	ranking := make(map[uuid.UUID]int)
	threshold := 0.25 // need to tune

	a := 1.0
	b := 1.0
	c := 0.6

	scaleFactor := 1.0

	if t5.state == 0 {
		scaleFactor = 0.2
	} else if t5.state == 1 {
		scaleFactor = 0.5
	} else if t5.state == 2 {
		scaleFactor = 0.6
	} else if t5.state == 3 {
		scaleFactor = 0.8
	}

	for _, agentB := range agentsOnBike {

		keyId := agentB.GetID()
		if keyId == t5.GetID() {
			continue
		}

		reputation := t5.QueryReputation(keyId)

		energyLevel := agentB.GetEnergyLevel()
		utility := (a * energyLevel) + (b * reputation) + (c * numberOfAgents)
		utilityNorm := utility / (2.0 + (c * utils.BikersOnBike))
		utilityNorm = utilityNorm * scaleFactor

		internalRanking[keyId] = utilityNorm

		if utilityNorm > threshold {
			ranking[keyId] = 0

		} else {
			ranking[keyId] = 1
		}

	}

	return ranking
}

func (t5 *team5Agent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {

	agentMap := t5.GetGameState().GetAgents()
	result := make(map[uuid.UUID]bool)
	pendingAgentUtility := make(map[uuid.UUID]float64)
	threshold := 0.5
	maxBikers := utils.BikersOnBike
	currentBikers := maxBikers - len(t5.GetFellowBikers())

	a := 1.0
	b := 1.0
	c := 1.0
	// energyMax := 1.0
	targetColor := t5.GetColour()

	scaleFactor := 1.0

	if t5.state == 0 {
		scaleFactor = 0.2
	} else if t5.state == 1 {
		scaleFactor = 0.5
	} else if t5.state == 2 {
		scaleFactor = 0.6
	} else if t5.state == 3 {
		scaleFactor = 0.8
	}

	for _, agentID := range pendingAgents {
		agentState := agentMap[agentID]

		key := agentState.GetID()
		reputation := t5.QueryReputation(key)
		energyLevel := agentState.GetEnergyLevel()
		pendingAgentColor := agentState.GetColour()

		isColorSame := 0.0

		if targetColor == pendingAgentColor {
			isColorSame = 1.0
		}

		// color has to be a 0/1 and replaced with
		utility := (a * energyLevel) + (b * reputation) + (c * isColorSame)
		utilityNorm := utility / 3.0
		utilityNorm = utilityNorm * scaleFactor

		pendingAgentUtility[agentID] = utilityNorm

	}

	type kv struct {
		Key   uuid.UUID
		Value float64
	}

	var ss []kv
	for k, v := range pendingAgentUtility {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value // Sorting in descending order
	})

	for i, pair := range ss {
		result[pair.Key] = i < currentBikers && pair.Value >= threshold
	}

	return result
}
