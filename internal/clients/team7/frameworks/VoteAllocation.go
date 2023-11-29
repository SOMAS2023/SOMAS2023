package frameworks

import (
	voting "SOMAS2023/internal/common/voting"
)

// Greedy implementation. Vote for our agent to get all of resource.
func VoteOnAllocationWrapper(voteInputs VoteInputs) voting.IdVoteMap {
	var vote voting.IdVoteMap // TODO: Import voting when rebase done and use IdVoteMap type.

	for _, agent_id := range voteInputs.Candidates.AgentCandidate {
		if agent_id == voteInputs.TeamSevenBikerId {
			vote[agent_id] = 1
		} else {
			vote[agent_id] = 0
		}
	}

	return vote
}
