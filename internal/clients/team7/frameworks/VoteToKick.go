package frameworks

import (
	"github.com/google/uuid"
)

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
		// Agent score depends on our average trust level of the agent.
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

// Assign a score to express approval/disapproval of a proposal.
/*func (voteHandler *VoteToKickAgentHandler) voteToKickScore(agent_id uuid.UUID) ScoreType {
	score := ScoreType(0.8) //TODO: Simple implementation for now. Will depend on factors such as opinion of agent and our agent's personality.
	return score
}*/
