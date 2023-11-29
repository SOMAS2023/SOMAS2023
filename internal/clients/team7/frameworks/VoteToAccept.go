package frameworks

// Basic implementation. Accept all new agents
func VoteToAcceptWrapper(voteInputs VoteInputs) Vote {
	var vote map[interface{}]interface{}
	vote_decision := true
	vote["decision"] = vote_decision
	return Vote{result: vote}
}

// TODO: Add functions for voting to kick an agent off the bike.
