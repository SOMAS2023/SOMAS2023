package frameworks

type VoteToAcceptAgentHandler struct {
	IDecisionFramework[VoteOnAgentsInput, MapIdBool]
}

func NewVoteToAcceptAgentHandler() *VoteToAcceptAgentHandler {
	return &VoteToAcceptAgentHandler{}
}

func (voteHandler *VoteToAcceptAgentHandler) GetDecision(inputs VoteOnAgentsInput) MapIdBool {
	vote := make(MapIdBool)
	threshold := ScoreType(0.5) // TODO: This could come from voteParameters in VoteInputs.

	for _, agent_id := range inputs.AgentCandidates {
		agent_score := voteHandler.voteToAcceptScore(agent_id)
		vote[agent_id] = agent_score > threshold
	}

	return vote
}

// Assign a score to express approval/disapproval of a proposal.
func (voteHandler *VoteToAcceptAgentHandler) voteToAcceptScore(agent_id interface{}) ScoreType {
	score := ScoreType(0.8) //TODO: Simple implementation for now. Will depend on factors such as opinion of agent and our agent's personality.
	return score
}
