package frameworks

import (
	voting "SOMAS2023/internal/common/voting"
)

// This file contains code for voting on the bike Dictator.

type VoteOnDictatorHandler struct {
	IDecisionFramework[VoteOnAgentsInput, voting.IdVoteMap]
}

func NewVoteOnDictatorHandler() *VoteOnDictatorHandler {
	return &VoteOnDictatorHandler{}
}

func (voteHandler *VoteOnDictatorHandler) GetDecision(inputs VoteOnAgentsInput) voting.IdVoteMap {
	agentScoreMap := make(voting.IdVoteMap)
	totalScore := 0.0

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

/*
// Assign a score to express approval/disapproval of an agent becoming Dictator.
func (voteHandler *VoteOnDictatorHandler) voteOnDictatorScore(agent_id interface{}) float64 {
	score := 0.8 //TODO: Simple implementation for now. Will depend on factors such as opinion of agent and our agent's personality.
	return score
}
*/
