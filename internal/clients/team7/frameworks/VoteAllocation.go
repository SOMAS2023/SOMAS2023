package frameworks

// Greedy implementation. Vote for our agent to get all of resource.
func VoteOnAllocationWrapper(voteInputs VoteInputs) Vote {
	var vote map[interface{}]interface{} // TODO: Import voting when rebase done and use IdVoteMap type.

	for _, agent_id := range voteInputs.Candidates {
		if agent_id == voteInputs.TeamSevenBikerId {
			vote[agent_id] = 1
		} else {
			vote[agent_id] = 0
		}
	}

	return Vote{result: vote}
}
