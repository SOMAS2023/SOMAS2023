package frameworks

import (
	voting "SOMAS2023/internal/common/voting"
)

// VoteDictator: This determines how our agent decides its preferences on which agent should be dictator of the bike.
// Each agent is assigned a score based on the average trust level our agent has had for them in previous rounds.
// The scores are then normalised such that they add up to 1 (as required by infrastructure), giving our agent's preference vote.

type VoteOnDictatorHandler struct {
	IDecisionFramework[VoteOnAgentsInput, voting.IdVoteMap]
}

func NewVoteOnDictatorHandler() *VoteOnDictatorHandler {
	return &VoteOnDictatorHandler{}
}

func (voteHandler *VoteOnDictatorHandler) GetDecision(inputs VoteOnAgentsInput) voting.IdVoteMap {
	agentScoreMap := make(voting.IdVoteMap)
	totalScore := 0.0
	// Assign a score to each agent based on our average trust for them in previous iterations.
	// Use this score to determine our agent's preference.
	for _, agent_id := range inputs.AgentCandidates {
		agentConnection, exists := inputs.CurrentSocialNetwork[agent_id]
		var agentScore float64
		if !exists {
			agentScore = 0.5
		} else {
			agentScore = agentConnection.GetAverageTrustLevels()
		}
		totalScore += agentScore
		agentScoreMap[agent_id] = agentScore
	}
	// Return a vote map where the sum of the votes is 1, as expected by the environment.
	vote := NormaliseVote(agentScoreMap, totalScore)

	return vote
}
