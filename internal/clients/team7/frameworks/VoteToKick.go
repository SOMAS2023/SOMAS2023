package frameworks

import (
	"github.com/google/uuid"
)

// VoteToKick: Determines how our agent votes on kicking a biker off its bike.
// If the average trust level our agent has had for the biker is less than a threshold, we vote to kick off the biker.
// Otherwise, we vote to allow the biker to stay on our bike.

type VoteToKickAgentHandler struct {
	IDecisionFramework[VoteOnAgentsInput, map[uuid.UUID]int]
}

func NewVoteToKickAgentHandler() *VoteToKickAgentHandler {
	return &VoteToKickAgentHandler{}
}

func (voteHandler *VoteToKickAgentHandler) GetDecision(inputs VoteOnAgentsInput) map[uuid.UUID]int {
	vote := make(map[uuid.UUID]int)
	threshold := ScoreType(0.1)

	for _, agent_id := range inputs.AgentCandidates {
		// Assign a score to each agent based on our average trust for them in previous iterations.
		// Use this score to determine whether to kick agent off bike.
		agentConnection, exists := inputs.CurrentSocialNetwork[agent_id]
		var averageTrustLevel float64
		if !exists {
			averageTrustLevel = 1
		} else {
			averageTrustLevel = agentConnection.GetAverageTrustLevels()
		}
		agentScore := ScoreType(averageTrustLevel)
		if agentScore < threshold {
			vote[agent_id] = 1
		} else {
			vote[agent_id] = 0
		}
	}

	return vote
}
