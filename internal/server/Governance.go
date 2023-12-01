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
	ruler := s.RulerElection(agents, governance)
	// communicate dictator
	bike.SetRuler(ruler)
	// get dictators direction choice
	rulerAgent := s.GetAgentMap()[ruler]
	var direction uuid.UUID
	switch governance {
	case utils.Dictatorship:
		direction = rulerAgent.DictateDirection()
	case utils.Leadership:
		// direction = rulerAgent.LeadDirection()
		LeaderWeights := make(map[uuid.UUID]float64)
		LeaderWeights = rulerAgent.LeaderAgentWeights() // TODO: implement this in BaseBiker
		direction = s.RunDemocraticAction(bike, governance, LeaderWeights)
	}
	return direction
}

func (s *Server) RulerElection(agents []objects.IBaseBiker, governance utils.Governance) uuid.UUID {
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
	return ruler
}

func (s *Server) RunDemocraticAction(bike objects.IMegaBike, governance utils.Governance, weights map[uuid.UUID]float64) uuid.UUID {
	// map of the proposed lootboxes by bike (for each bike a list of lootbox proposals is made, with one lootbox proposed by each agent on the bike)
	agents := bike.GetAgents()
	proposedDirections := make([]uuid.UUID, len(agents))
	for i, agent := range agents {
		// agents that have decided to stay on the bike (and that haven't been kicked off it)
		// will participate in the voting for the directions
		// ---------------------------VOTING ROUTINE - STEP 1 ---------------------
		if agent.GetBikeStatus() {
			proposedDirections[i] = agent.ProposeDirection()
		}
	}

	// pass the pitched directions of a bike to all agents on that bike and get their final vote
	finalVotes := make([]voting.LootboxVoteMap, len(agents))

	// // Merge these two loops? Done below but kept this just in case I messed up the logic while restructuring
	// if governance == utils.Leadership {
	// 	// Leadership democracy
	// 	for i, agent := range agents {
	// 		for id, weight := range weights {
	// 			if id == agent.GetID() {
	// 				// ---------------------------VOTING ROUTINE - STEP 2 ---------------------
	// 				finalVotes[i] = agent.FinalDirectionVote((proposedDirections))
	// 				// Weight the votes according to the weights of the agents imposed by the leader if the governance is leadership democracy
	// 				for key, value := range finalVotes[i] {
	// 					finalVotes[i][key] = value * weight
	// 				}
	// 				// Assuming that the every agents on the bike has voted, we can penalise them for an energy loss due to voting power
	// 				agent.UpdateEnergyLevel(-utils.LeadershipDemocracyPenalty)
	// 			}
	// 		}
	// 	}
	// } else {
	// 	// Deliberative democracy
	// 	for i, agent := range agents {
	// 		// ---------------------------VOTING ROUTINE - STEP 2 ---------------------
	// 		finalVotes[i] = agent.FinalDirectionVote((proposedDirections))
	// 		// Assuming that the every agents on the bike has voted, we can penalise them for an energy loss due to voting power
	// 		agent.UpdateEnergyLevel(-utils.DeliberativeDemocracyPenalty)
	// 	}
	// }

	for i, agent := range agents {
		// ---------------------------VOTING ROUTINE - STEP 2 ---------------------
		finalVotes[i] = agent.FinalDirectionVote((proposedDirections))

		if governance == utils.Leadership { // Leadership democracy
			for id, weight := range weights {
				if id == agent.GetID() {
					// Weight the votes according to the weights of the agents imposed by the leader if the governance is leadership democracy
					for key, value := range finalVotes[i] {
						finalVotes[i][key] = value * weight
					}
					// Assuming that the every agents on the bike has voted, we can penalise them for an energy loss due to voting power
					agent.UpdateEnergyLevel(-utils.LeadershipDemocracyPenalty)
				}
			}
		} else { // Deliberative democracy
			// Assuming that the every agents on the bike has voted, we can penalise them for an energy loss due to voting power
			agent.UpdateEnergyLevel(-utils.DeliberativeDemocracyPenalty)
		}
	}

	// ---------------------------VOTING ROUTINE - STEP 3 --------------
	direction := s.GetWinningDirection(finalVotes)
	return direction
}
