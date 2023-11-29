package frameworks

//"github.com/google/uuid"
import (
	"fmt"
)

func VoteToKickWrapper(voteInputs VoteInputs) Vote {
	var vote map[interface{}]interface{}
	var threshold ScoreType
	var agent_score ScoreType
	var vote_decision interface{}
	threshold = 0.5 // TODO: This could come from voteParameters in VoteInputs.

	for _, agent_id := range voteInputs.Candidates {
		agent_score = VoteToKickScore(agent_id)
		switch voteInputs.VoteParameters {
		case Proportion:
			vote_decision = agent_score
		case YesNo:
			vote_decision = VoteToKickYesNo(agent_score, threshold)
		default:
			fmt.Println("New decision type!")
			vote_decision = agent_score
		}

		vote[agent_id] = vote_decision
	}

	return Vote{result: vote}
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
