package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

func (s *Server) RunDictatorAction(bike objects.IMegaBike) uuid.UUID {
	// vote for dictator
	agents := bike.GetAgents()
	votes := make([]voting.IdVoteMap, len(agents))
	for i, agent := range agents {
		votes[i] = agent.VoteDictator()
	}

	IVotes := make([]voting.IVoter, len(votes))
	for i, vote := range votes {
		IVotes[i] = vote
	}

	dictator := voting.WinnerFromDist(IVotes)
	// communicate dictator
	bike.SetRuler(dictator)
	// get dictators direction choice
	dictatorAgent := s.GetAgentMap()[dictator]
	return dictatorAgent.DictateDirection()
}
