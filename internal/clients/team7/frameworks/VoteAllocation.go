package frameworks

import (
	voting "SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

type VoteOnAllocationInput struct {
	AgentCandidates []uuid.UUID
	MyPersonality   *Personality
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

	// Whether we share depends on how agreeable we are
	// Low agreeableness => no sharing => we give ourselves 100% of the share
	// Mid agreeableness => some sharing => give ourselves certain share and divide rest among others
	// High agreeableness => all sharing => we give everyone equal shares of the vote
	agreeableness := inputs.MyPersonality.Agreeableness
	candidates := inputs.AgentCandidates
	numcandidates := len(candidates)
	othershares := agreeableness / float64(numcandidates)
	myshare := 1 - (othershares * (float64(numcandidates) - 1))

	for _, agentId := range inputs.AgentCandidates {
		if agentId == inputs.MyId {
			vote[agentId] = myshare
		} else {
			vote[agentId] = othershares
		}
	}

	return vote
}
