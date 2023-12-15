package frameworks

// VoteToAccept: Determines how our agent votes on accepting a new biker onto its bike.
// If the average trust level our agent has had for the biker is less than a threshold, we reject the biker.
// Otherwise, we accept the biker onto our bike.
type VoteToAcceptAgentHandler struct {
	IDecisionFramework[VoteOnAgentsInput, MapIdBool]
}

func NewVoteToAcceptAgentHandler() *VoteToAcceptAgentHandler {
	return &VoteToAcceptAgentHandler{}
}

func (voteHandler *VoteToAcceptAgentHandler) GetDecision(inputs VoteOnAgentsInput) MapIdBool {
	vote := make(MapIdBool)
	threshold := ScoreType(0.4)

	for _, agent_id := range inputs.AgentCandidates {
		// Assign a score to each agent based on our average trust for them in previous iterations.
		// Use this score to determine whether to accept agent onto bike.
		agentConnection, exists := inputs.CurrentSocialNetwork[agent_id]
		var agentScore float64
		if !exists {
			agentScore = 0.5
		} else {
			agentScore = agentConnection.GetAverageTrustLevels()
		}
		vote[agent_id] = ScoreType(agentScore) < threshold
	}

	return vote
}
