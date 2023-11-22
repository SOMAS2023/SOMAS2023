package frameworks

func VoteOnProposalsWrapper(voteInputs VoteInputs) Vote {
	var vote Map
	vote_decision := true
	vote["decision"] = vote_decision
	return Vote{result: vote}
}

// TODO: Add functions for voting on which loot box to go to.
