package frameworks

func VoteToKickWrapper(voteInputs VoteInputs) Vote {
	var vote Map
	var threshold ScoreType
	threshold = 0.5 // TODO: This could come from voteParameters in VoteInputs.

	score := VoteToKickScore(voteInputs)
	vote_decision := VoteToKickYesNo(score, threshold)
	vote["decision"] = vote_decision
	return Vote{result: vote}
}

// Assign a score to express approval/disapproval of a proposal.
func VoteToKickScore(voteInputs VoteInputs) ScoreType {
	var score ScoreType
	score = 0.8 //TODO: Simple implementation for now. Will depend on factors such as opinion of agent.
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
