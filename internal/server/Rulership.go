package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

func (s *Server) RunRulerAction(bike objects.IMegaBike, governance utils.Governance) uuid.UUID {
	// vote for dictator
	agents := bike.GetAgents()
	votes := make([]voting.IdVoteMap, len(agents))
	for i, agent := range agents {
		switch governance {
		case utils.Dictatorship:
			votes[i] = agent.VoteDictator()
		case utils.Leadership:
			votes[i] = agent.VoteLeader()
		}
	}

	IVotes := make([]voting.IVoter, len(votes))
	for i, vote := range votes {
		IVotes[i] = vote
	}

	ruler := voting.WinnerFromDist(IVotes)
	// communicate dictator
	bike.SetRuler(ruler)
	// get dictators direction choice
	rulerAgent := s.GetAgentMap()[ruler]
	var direction uuid.UUID
	switch governance {
	case utils.Dictatorship:
		direction = rulerAgent.DictateDirection()
	case utils.Leadership:
		direction = rulerAgent.LeadDirection()
	}
	return direction
}
