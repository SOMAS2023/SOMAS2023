package frameworks

type VoteToKickAgentHandler struct {
	IDecisionFramework[VoteOnAgentsInput, MapIdBool]
}

func NewVoteToKickAgentHandler() *VoteToKickAgentHandler {
	return &VoteToKickAgentHandler{}
}

func (voteHandler *VoteToKickAgentHandler) GetDecision(inputs VoteOnAgentsInput) MapIdBool {
	vote := make(MapIdBool)
	threshold := ScoreType(0.4)

	for _, agent_id := range inputs.AgentCandidates {
		// Agent score depends on our average trust level of the agent.
		agentConnection, exists := inputs.CurrentSocialNetwork[agent_id]
		var averageTrustLevel float64
		if !exists {
			averageTrustLevel = 0.5
		} else {
			averageTrustLevel = agentConnection.GetAverageTrustLevels()
		}
		agentScore := ScoreType(averageTrustLevel)
		vote[agent_id] = agentScore < threshold
	}

	return vote
}

// Assign a score to express approval/disapproval of a proposal.
/*func (voteHandler *VoteToKickAgentHandler) voteToKickScore(agent_id uuid.UUID) ScoreType {
	score := ScoreType(0.8) //TODO: Simple implementation for now. Will depend on factors such as opinion of agent and our agent's personality.
	return score
}*/
