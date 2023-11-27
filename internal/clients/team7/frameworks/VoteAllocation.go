package frameworks

//"github.com/google/uuid"

// func VoteOnAllocationWrapper(fellowBikers []objects.IBaseBiker) Vote {
func VoteOnAllocationWrapper(voteInputs VoteInputs) Vote {
	var vote map[interface{}]interface{} // TODO: Import voting when rebase done and use IdVoteMap type.
	var num_agents float64
	num_agents = float64(len(voteInputs.Candidates))

	for _, agent_id := range voteInputs.Candidates {
		vote[agent_id] = 1 / num_agents
	}

	return Vote{result: vote}
}
