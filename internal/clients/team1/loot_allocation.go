package team1

import (
	voting "SOMAS2023/internal/common/voting"
	"fmt"
	"github.com/google/uuid"
	"math"
)

// ---------------LOOT ALLOCATION FUNCTIONS------------------
// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Biker1) DecideAllocation() voting.IdVoteMap {

	fellowBikers := bb.GetFellowBikers()
	if len(fellowBikers) == 1 {
		return voting.IdVoteMap{bb.GetID(): 1}
	}

	sumEnergyNeeds := 0.0
	helpfulAllocation := make(map[uuid.UUID]float64)
	selfishAllocation := make(map[uuid.UUID]float64)

	for _, agent := range fellowBikers {
		energy := agent.GetEnergyLevel()
		energyNeed := 1.0 - energy
		helpfulAllocation[agent.GetID()] = energyNeed
		selfishAllocation[agent.GetID()] = energyNeed
		sumEnergyNeeds = sumEnergyNeeds + energyNeed
	}

	for agentId := range helpfulAllocation {
		helpfulAllocation[agentId] /= sumEnergyNeeds
	}

	sumEnergyNeeds -= (1.0 - bb.GetEnergyLevel()) // remove our energy need from the sum

	for agentId := range selfishAllocation {
		if agentId != bb.GetID() {
			selfishAllocation[agentId] = (selfishAllocation[agentId] / sumEnergyNeeds) * bb.GetEnergyLevel() //NB assuming energy is 0-1
		}
	}

	//3/4) Look in success vector to see relative success of each agent and calculate selfishness score using suc-rel chart (0-1)
	//TI - Around line 350, we have Soma`s pseudocode on agent opinion held in bb.Opinion.opinion, lets assume its normalized between 0-1
	selfishnessScore := make(map[uuid.UUID]float64)
	runningScore := 0.0

	for _, agent := range fellowBikers {
		if agent.GetID() != bb.GetID() {
			score := bb.GetSelfishness(agent)
			id := agent.GetID()
			selfishnessScore[id] = score
			runningScore = runningScore + selfishnessScore[id]
		}
	}

	selfishnessScore[bb.GetID()] = runningScore / float64((len(fellowBikers) - 1))

	//5) Linearly interpolate between selfish and helpful allocations based on selfishness score
	distribution := make(map[uuid.UUID]float64)
	runningDistribution := 0.0
	for _, agent := range fellowBikers {
		id := agent.GetID()
		Adistribution := (selfishnessScore[id] * selfishAllocation[id]) + ((1.0 - selfishnessScore[id]) * helpfulAllocation[id])
		distribution[id] = Adistribution
		runningDistribution = runningDistribution + Adistribution
	}
	for agentId := range distribution {
		distribution[agentId] = distribution[agentId] / runningDistribution // Normalise!
	}
	if math.IsNaN(distribution[bb.GetID()]) {
		fmt.Println("Distribution is NaN")
	}
	return distribution
}

// ---------------END OF LOOT ALLOCATION FUNCTIONS------------------
