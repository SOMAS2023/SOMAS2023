// VOTING FOR LEADER/DICTATOR AND GOVERNANCE TYPE

package team1

import (
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
)

// -------------------GOVERMENT CHOICE FUNCTIONS--------------------------

// Not implemented on Server yet so this is just a placeholder
func (bb *Biker1) DecideGovernance() utils.Governance {
	bb.setOpinions()
	if bb.DecideDictatorship() {
		return utils.Dictatorship
	} else if bb.DecideLeadership() {
		return utils.Leadership
	} else {
		// Democracy
		return utils.Democracy
	}
	//return 2
}

// Might be unnecesary as this is the default goverment choice for us
func (bb *Biker1) DecideDemocracy() bool {
	founding_agents := bb.GetAllAgents()
	totalOpinion := 0.0
	reputation := bb.DetermineOurReputation()
	for _, agent := range founding_agents {
		opinion, ok := bb.opinions[agent.GetID()]
		if ok {
			totalOpinion = totalOpinion + opinion.opinion
		}
	}
	normOpinion := totalOpinion / float64(len(founding_agents))
	if (normOpinion > democracyOpinonThreshold) || (reputation > democracyReputationThreshold) {
		return true
	} else {
		return false
	}
}

func (bb *Biker1) DecideLeadership() bool {
	founding_agents := bb.GetAllAgents()
	totalOpinion := 0.0
	reputation := bb.DetermineOurReputation()
	for _, agent := range founding_agents {
		opinion, ok := bb.opinions[agent.GetID()]
		if ok {
			totalOpinion = totalOpinion + opinion.opinion
		}
	}
	normOpinion := totalOpinion / float64(len(founding_agents))
	if (normOpinion > leadershipOpinionThreshold) || (reputation > leadershipReputationThreshold) {
		return true
	} else {
		return false
	}
}

func (bb *Biker1) DecideDictatorship() bool {
	founding_agents := bb.GetAllAgents()
	totalOpinion := 0.0
	reputation := bb.DetermineOurReputation()
	for _, agent := range founding_agents {
		opinion, ok := bb.opinions[agent.GetID()]
		if ok {
			totalOpinion = totalOpinion + opinion.opinion
		}
	}
	normOpinion := totalOpinion / float64(len(founding_agents))
	if (normOpinion > dictatorshipOpinionThreshold) || (reputation > dictatorshipReputationThreshold) {
		return true
	} else {
		return false
	}
}

// ----------------------LEADER/DICTATOR VOTING FUNCTIONS------------------
func (bb *Biker1) VoteLeader() voting.IdVoteMap {

	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()

	for _, agent := range fellowBikers {
		votes[agent.GetID()] = 0.0
		if agent.GetID() != bb.GetID() {
			val, ok := bb.opinions[agent.GetID()]
			if ok {
				votes[agent.GetID()] = val.opinion
			}
		} else {
			votes[agent.GetID()] = 0.7
		}
	}
	votesum := 0.0
	for _, agent := range fellowBikers {
		votesum = votesum + votes[agent.GetID()]
	}
	for _, agent := range fellowBikers {
		votes[agent.GetID()] = votes[agent.GetID()] / votesum
	}
	return votes
}

func (bb *Biker1) VoteDictator() voting.IdVoteMap {

	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetAllAgents()

	for _, agent := range fellowBikers {
		votes[agent.GetID()] = 0.0
		if agent.GetID() != bb.GetID() {
			val, ok := bb.opinions[agent.GetID()]
			if ok {
				votes[agent.GetID()] = val.opinion
			}
		} else {
			votes[agent.GetID()] = 1.0
		}
	}
	votesum := 0.0
	for _, agent := range fellowBikers {
		votesum = votesum + votes[agent.GetID()]
	}
	for _, agent := range fellowBikers {
		votes[agent.GetID()] = votes[agent.GetID()] / votesum
	}
	votesum2 := 0.0
	for _, agent := range fellowBikers {
		votesum2 = votesum2 + votes[agent.GetID()]
	}
	return votes
}

//--------------------END OF LEADER/DICTATOR VOTING FUNCTIONS------------------

//--------------------END OF GOVERMENT CHOICE FUNCTIONS------------------
