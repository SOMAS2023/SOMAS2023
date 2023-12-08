package team4

import (
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

func (agent *BaselineAgent) DecideGovernance() utils.Governance {
	// Change behaviour here to return different governance
	return utils.Democracy
}

func (agent *BaselineAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	agent.UpdateDecisionData()
	spare := minFellowBikers - agent.capacity
	decision := make(map[uuid.UUID]bool)

	var scoredAgents []agentScore

	for _, pendingAgent := range pendingAgents {
		reputation := agent.reputation[pendingAgent]
		honesty := agent.honestyMatrix[pendingAgent]
		scoredAgents = append(scoredAgents, agentScore{ID: pendingAgent, Score: ((reputationWeight * reputation) + (honestyWeight * honesty))})
	}
	// Sort the slice based on the combined score
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].Score > scoredAgents[j].Score
	})

	// Make decisions based on the sorted slice
	for i, scoredAgent := range scoredAgents {
		// Example decision making logic
		if i < spare {
			decision[scoredAgent.ID] = true // Accept if there's spare capacity
		} else {
			decision[scoredAgent.ID] = false // Reject if no capacity
		}
	}
	return decision
}
func (agent *BaselineAgent) ChangeBike() uuid.UUID {
	agent.UpdateDecisionData()
	megaBikes := agent.GetGameState().GetMegaBikes()
	optimalBike := agent.currentBike
	weight := float64(-99)
	for _, bike := range megaBikes {
		if bike.GetID() != uuid.Nil && bike.GetID() != agent.currentBike { //get all bikes apart from our agent's bike
			bikeWeight := float64(0)

			for _, biker := range bike.GetAgents() {
				if biker.GetColour() == agent.GetColour() {
					bikeWeight += 1.8
				} else {
					bikeWeight += 1
				}
			}

			if bikeWeight > weight {
				optimalBike = bike.GetID()
			}
		}
	}
	agent.optimalBike = optimalBike
	return optimalBike
}

func (agent *BaselineAgent) VoteForKickout() map[uuid.UUID]int {
	agent.UpdateDecisionData()
	fmt.Println("Vote for Kickout")
	voteResults := make(map[uuid.UUID]int)

	fellowBikers := agent.GetFellowBikers()
	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)

	if e1 != nil || e2 != nil {
		panic("unexpected error!")

	}
	combined := make(map[uuid.UUID]float64)
	worstRank := float64(1)

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		if combined[fellowID] == worstRank && fellowID != uuid.Nil {

			if fellowID != agent.GetID() {
				combined[fellowID] = reputationRank[fellowID] * honestyRank[fellowID]
				if combined[fellowID] < worstRank {
					worstRank = combined[fellowID]
				}
			} else {
				combined[fellowID] = 1.0
			}
		}
	}

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		if agent.capacity > minFellowBikers {
			if combined[fellowID] == worstRank && fellowID != uuid.Nil {
				if agent.reputation[fellowID] < agent.getReputationAverage() || agent.honestyMatrix[fellowID] < agent.getHonestyAverage() {
					if agent.capacity > 4 {
						voteResults[fellowID] = 1
					} else {
						voteResults[fellowID] = 0
					}
				}
			} else {
				voteResults[fellowID] = 0
			}
		} else {
			voteResults[fellowID] = 0
		}
	}
	voteResults[agent.GetID()] = 0
	println("the voting results are:", voteResults)
	return voteResults
}

///////////////////////////////////// LEADER FUNCTIONS ///////////////////////////////////////

// defaults to an equal distribution over all agents for all actions
func (agent *BaselineAgent) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	weights := make(map[uuid.UUID]float64)
	fellows := agent.GetFellowBikers()
	for _, fellow := range fellows {
		if fellow.GetID() != uuid.Nil {
			weights[fellow.GetID()] = 1.0
		} else {
			weights[fellow.GetID()] = 0.0
		}
	}
	return weights
}

func (agent *BaselineAgent) VoteLeader() voting.IdVoteMap {
	agent.UpdateDecisionData()
	votes := make(voting.IdVoteMap)
	fellowBikers := agent.GetFellowBikers()
	totalsum := float64(0)

	var scoredAgents []agentScore

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		reputation := agent.reputation[fellowID]
		honesty := agent.honestyMatrix[fellowID]
		scoredAgents = append(scoredAgents, agentScore{ID: fellowID, Score: ((leaderRepWeight * reputation) + (leaderHonestWeight * honesty))})
	}
	// Sort the slice based on the combined score
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].Score > scoredAgents[j].Score
	})

	for i, scoredAgent := range scoredAgents {
		weight := float64(len(scoredAgents) - i)
		votes[scoredAgent.ID] = weight
		totalsum += weight
	}
	votes[agent.GetID()] = 20.0
	totalsum += 20.0
	//normalize the vote
	for _, scoredAgent := range scoredAgents {
		votes[scoredAgent.ID] = votes[scoredAgent.ID] / totalsum
	}
	return votes
}

/////////////////////////////////// DICATOR FUNCTIONS /////////////////////////////////////

func (agent *BaselineAgent) VoteDictator() voting.IdVoteMap {
	agent.UpdateDecisionData()
	votes := make(voting.IdVoteMap)
	fellowBikers := agent.GetFellowBikers()
	totalsum := float64(0)
	var scoredAgents []agentScore

	for _, fellow := range fellowBikers {
		fellowID := fellow.GetID()
		reputation := agent.reputation[fellowID]
		honesty := agent.honestyMatrix[fellowID]
		scoredAgents = append(scoredAgents, agentScore{ID: fellowID, Score: ((dictatorRepWeight * reputation) + (dictatorHonestWeight * honesty))})
	}
	// Sort the slice based on the combined score
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].Score > scoredAgents[j].Score
	})

	// Make decisions based on the sorted slice
	for i, scoredAgent := range scoredAgents {
		weight := float64(len(scoredAgents) - i)
		votes[scoredAgent.ID] = weight
		totalsum += weight
	}
	votes[agent.GetID()] = 20.0
	totalsum += 20.0
	//normalize the vote
	for _, scoredAgent := range scoredAgents {
		votes[scoredAgent.ID] = votes[scoredAgent.ID] / totalsum
	}
	return votes
}

func (agent *BaselineAgent) DecideKickOut() []uuid.UUID {
	fmt.Println("Decide Kickout")
	kickoutResults := make([]uuid.UUID, 0)
	agent.UpdateDecisionData()

	fellowBikers := agent.GetFellowBikers()
	if agent.capacity > dictatorMinFellowBikers {

		reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
		honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
		if e1 != nil || e2 != nil {
			panic("unexpected error!")
		}
		combined := make(map[uuid.UUID]float64)
		worstRank := float64(1)

		for _, fellow := range fellowBikers {
			fellowID := fellow.GetID()
			if combined[fellowID] == worstRank && fellowID != uuid.Nil {

				if fellowID != agent.GetID() {
					combined[fellowID] = reputationRank[fellowID] * honestyRank[fellowID]
					if combined[fellowID] < worstRank {
						worstRank = combined[fellowID]
					}
				} else {
					combined[fellowID] = 1.0
				}
			}
		}
		for _, fellow := range fellowBikers {
			fellowID := fellow.GetID()
			if fellowID != agent.GetID() {
				if combined[fellowID] == worstRank && fellowID != uuid.Nil {
					if agent.reputation[fellowID] < agent.getReputationAverage() || agent.honestyMatrix[fellowID] < agent.getHonestyAverage() {
						kickoutResults = append(kickoutResults, fellowID)
					}
				}
			}
		}
	}
	return kickoutResults

}
