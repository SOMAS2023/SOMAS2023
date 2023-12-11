package team4

import (
	"SOMAS2023/internal/common/voting"
)

func (agent *BaselineAgent) DecideAllocation() voting.IdVoteMap {
	//fmt.Println("Decide Allocation")
	agent.UpdateDecisionData()
	distribution := make(voting.IdVoteMap) //make(map[uuid.UUID]float64)
	fellowBikers := agent.GetFellowBikers()
	totalEnergySpent := float64(0)
	totalAllocation := float64(0)

	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		energyLog := agent.energyHistory[fellowID]
		energySpent := energyLog[len(energyLog)-2] - energyLog[len(energyLog)-1]
		totalEnergySpent += energySpent
		distribution[fellow.GetID()] = float64((reputationWeight * reputationRank[fellowID]) + (honestyWeight * honestyRank[fellowID]) + (energySpentWeight * energySpent) + (energyLevelWeight * fellow.GetEnergyLevel()))
		// In the case where the I am the same colour as the lootbox
		if fellowID == agent.GetID() {
			distribution[fellow.GetID()] = float64((reputationWeight * reputationRank[fellowID]) + (honestyWeight * honestyRank[fellowID]) + (energySpentWeight * energySpent * 1.5) + (energyLevelWeight * fellow.GetEnergyLevel()))
			if agent.lootBoxColour == agent.GetColour() {
				distribution[fellow.GetID()] = float64((reputationWeight * reputationRank[fellowID]) + (honestyWeight * honestyRank[fellowID]) + (energySpentWeight * energySpent * 1.5) + (energyLevelWeight * fellow.GetEnergyLevel() * 1.5))
			}
		}
		totalAllocation += distribution[fellow.GetID()]
	}

	//normalize the distribution
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		distribution[fellowID] = distribution[fellowID] / totalAllocation
	}

	return distribution
}

/////////////////////////////////// DICATOR FUNCTIONS /////////////////////////////////////

func (agent *BaselineAgent) DecideDictatorAllocation() voting.IdVoteMap {
	//fmt.Println("Dictate Allocation")
	agent.UpdateDecisionData()
	distribution := make(voting.IdVoteMap)
	fellowBikers := agent.GetFellowBikers()
	totalEnergySpent := float64(0)
	totalAllocation := float64(0)

	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		energyLog := agent.energyHistory[fellowID]
		energySpent := energyLog[len(energyLog)-2] - energyLog[len(energyLog)-1]
		totalEnergySpent += energySpent
		distribution[fellow.GetID()] = float64((reputationWeight * reputationRank[fellowID]) + (honestyWeight * honestyRank[fellowID]) + (energySpentWeight * energySpent) + (energyLevelWeight * fellow.GetEnergyLevel()))
		// In the case where the I am the same colour as the lootbox
		if fellowID == agent.GetID() {
			distribution[fellow.GetID()] = float64((reputationWeight * reputationRank[fellowID]) + (honestyWeight * honestyRank[fellowID]) + (energySpentWeight * energySpent * 1.5) + (energyLevelWeight * fellow.GetEnergyLevel()))
			if agent.lootBoxColour == agent.GetColour() {
				distribution[fellow.GetID()] = float64((reputationWeight * reputationRank[fellowID]) + (honestyWeight * honestyRank[fellowID]) + (energySpentWeight * energySpent * 1.5) + (energyLevelWeight * fellow.GetEnergyLevel() * 1.5))
			}
		}
		totalAllocation += distribution[fellow.GetID()]
	}
	//normalize the distribution
	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		distribution[fellowID] = distribution[fellowID] / totalAllocation
	}

	return distribution
}
