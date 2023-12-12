package frameworks

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
		// Agent score depends on our average trust level of the agent.
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

/*
// Assign a score to express approval/disapproval of a proposal.
func (voteHandler *VoteToAcceptAgentHandler) voteToAcceptScore(agent_id interface{}) ScoreType {
	score := ScoreType(0.8) //TODO: Simple implementation for now. Will depend on factors such as opinion of agent and our agent's personality.
	return score
}
*/
