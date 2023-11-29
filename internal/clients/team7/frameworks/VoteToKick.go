package frameworks

//"github.com/google/uuid"

func VoteToKickWrapper(voteInputs VoteInputs) MapIdBool {
	var vote MapIdBool
	var threshold ScoreType
	var agent_score ScoreType
	threshold = 0.5 // TODO: This could come from voteParameters in VoteInputs.

	for _, agent_id := range voteInputs.Candidates.AgentCandidate {
		agent_score = VoteToKickScore(agent_id)
		vote[agent_id] = VoteToKickYesNo(agent_score, threshold)

	}

	return vote
}

// Assign a score to express approval/disapproval of a proposal.
func VoteToKickScore(agent_id interface{}) ScoreType {
	var score ScoreType
	score = 0.8 //TODO: Simple implementation for now. Will depend on factors such as opinion of agent and our agent's personality.
	return score
}

// Function to convert vote to bool if neccesary.
// TODO: Could this just be a common function for all votes?
func VoteToKickYesNo(score ScoreType, threshold ScoreType) bool {
	var decision bool
	if score > threshold {
		decision = true
	} else {
		decision = false
	}

	return decision
}
