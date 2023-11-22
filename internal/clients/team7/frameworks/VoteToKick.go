package frameworks

func VoteToKickWrapper(voteInputs VoteInputs) Vote {
	var vote Map
	vote_decision := true
	vote["decision"] = vote_decision
	return Vote{result: vote}
}

// TODO: Add functions for voting to kick an agent off the bike.
