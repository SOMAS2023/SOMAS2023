// VOTING FOR LEADER/DICTATOR AND GOVERNANCE TYPE

package team1

import (
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
)
// -------------------GOVERMENT CHOICE FUNCTIONS--------------------------

// Not implemented on Server yet so this is just a placeholder
func (bb *Biker1) DecideGovernance() utils.Governance {
	if bb.DecideDictatorship() {
		return 2
	} else if bb.DecideLeadership() {
		return 1
	} else {
		// Democracy
		return 0
	}
}

// Might be unnecesary as this is the default goverment choice for us
func (bb *Biker1) DecideDemocracy() bool {
	founding_agents := bb.GetAllAgents()
	totalOpinion := 0.0
	reputation := bb.ourReputation()
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
	reputation := bb.ourReputation()
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
	reputation := bb.ourReputation()
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

	maxOpinion := 0.0
	leaderVote := bb.GetID()
	for _, agent := range fellowBikers {
		votes[agent.GetID()] = 0.0
		avgOp := bb.GetAverageOpinionOfAgent(agent.GetID())
		if agent.GetID() != bb.GetID() {
			val, ok := bb.opinions[agent.GetID()]
			if ok {
				avgOp = (avgOp + val.opinion) / 2
			}
		}
		if avgOp > maxOpinion {
			maxOpinion = avgOp
			leaderVote = agent.GetID()
		}
	}
	votes[leaderVote] = 1.0
	return votes
}

func (bb *Biker1) VoteDictator() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	fellowBikers := bb.GetFellowBikers()

	maxOpinion := 0.0
	leaderVote := bb.GetID()
	for _, agent := range fellowBikers {
		votes[agent.GetID()] = 0.0
		avgOp := bb.GetAverageOpinionOfAgent(agent.GetID())
		if agent.GetID() != bb.GetID() {
			val, ok := bb.opinions[agent.GetID()]
			if ok {
				avgOp = (avgOp + 3*val.opinion) / 4
			}
		}
		if avgOp > maxOpinion {
			maxOpinion = avgOp
			leaderVote = agent.GetID()
		}
	}
	votes[leaderVote] = 1.0
	return votes
}

//--------------------END OF LEADER/DICTATOR VOTING FUNCTIONS------------------




//--------------------END OF GOVERMENT CHOICE FUNCTIONS------------------