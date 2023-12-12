package frameworks

// This code returns the governance type the agent wants.
// Initial strategy is to always vote for "deliberative democracy".

import (
	utils "SOMAS2023/internal/common/utils"
)

type VoteOnGovernanceHandler struct {
	IDecisionFramework[interface{}, utils.Governance]
}

func NewVoteOnGovernanceHandler() *VoteOnGovernanceHandler {
	return &VoteOnGovernanceHandler{}
}

func (voteHandler *VoteOnGovernanceHandler) GetDecision() utils.Governance {

	return utils.Democracy
}
