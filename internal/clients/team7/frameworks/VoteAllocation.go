package frameworks

import (
	voting "SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

type VoteOnAllocationInput struct {
	AgentCandidates []uuid.UUID
	MyId            uuid.UUID
}

type VoteOnAllocationHandler struct {
	IDecisionFramework[VoteOnAllocationInput, voting.IdVoteMap]
}

func NewVoteOnAllocationHandler() *VoteOnAllocationHandler {
	return &VoteOnAllocationHandler{}
}

func (voteHandler *VoteOnAllocationHandler) GetDecision(inputs VoteOnAllocationInput) voting.IdVoteMap {
	vote := make(voting.IdVoteMap)

	for _, agentId := range inputs.AgentCandidates {
		if agentId == inputs.MyId {
			vote[agentId] = 1
		} else {
			vote[agentId] = 0
		}
	}

	return vote
}
