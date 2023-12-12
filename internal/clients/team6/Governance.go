package team6

import (
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
)

func (bb *Team6Biker) DecideGovernance() utils.Governance {
	// choose the majority governance if energy level is too low to change bike -- undecided:

	return utils.Democracy // always select Democracy as goverence
}

func (bb *Team6Biker) VoteLeader() voting.IdVoteMap {

	fellowBikers := bb.GetFellowBikers()
	votes := make(voting.IdVoteMap)

	for _, agent := range fellowBikers {

		bikerID := agent.GetID()
		votes[bikerID] = 0.0

		if bikerID != bb.GetID() {
			trust := bb.QueryReputation(bikerID)
			if agent.GetColour() == bb.GetColour() {
				votes[bikerID] = trust // Identical colour, vote weights = trust
			} else {
				votes[bikerID] = trust * 0.6 // else weighted to 0.6
			}
		} else {
			votes[bb.GetID()] = 1.0
		}
	}
	return votes
}

func (bb *Team6Biker) VoteDictator() voting.IdVoteMap {

	fellowBikers := bb.GetFellowBikers()
	votes := make(voting.IdVoteMap)

	for _, agent := range fellowBikers {

		bikerID := agent.GetID()
		votes[bikerID] = 0.0

		if bikerID != bb.GetID() {
			trust := bb.QueryReputation(bikerID)
			if agent.GetColour() == bb.GetColour() {
				votes[bikerID] = trust * 0.5 // Identical colour, vote weights = trust * 0.5
			} else {
				votes[bikerID] = trust * 0.1 // else weighted to 0.1
			}
		} else {
			votes[bb.GetID()] = 1.0
		}
	}
	return votes
}
