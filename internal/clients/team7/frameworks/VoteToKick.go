package frameworks

import (
	"github.com/google/uuid"
)

func VoteToKickWrapper(voteInputs VoteInputs) Vote {
	var vote map[uuid.UUID]interface{}
	var threshold ScoreType
	var agent_score ScoreType
	threshold = 0.5 // TODO: This could come from voteParameters in VoteInputs.

	for _, agent_id := range voteInputs.Candidates {
		agent_score = VoteToKickScore(agent_id)
		vote_decision := VoteToKickYesNo(agent_score, threshold)
		vote[agent_id] = vote_decision
	}

	return Vote{result: vote}
}

// Assign a score to express approval/disapproval of a proposal.
func VoteToKickScore(agent_id uuid.UUID) ScoreType {
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
