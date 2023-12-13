package team6

import (
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"slices"

	"github.com/google/uuid"
)

func (bb *Team6Biker) DecideGovernance() utils.Governance {
	var bikeList []uuid.UUID
	for _, bike := range bb.GetGameState().GetMegaBikes() {
		bikeList = append(bikeList, bike.GetID())
	}
	if !slices.Contains(bikeList, bb.GetBike()) {
		return utils.Dictatorship
	}
	// choose the majority governance if energy level is too low to change bike -- undecided:
	fellowBikers := bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetAgents()
	var sameColourCount int
	sameColourCount = 0
	// fmt.Println(fellowBikers)
	for _, agent := range fellowBikers {
		// fmt.Println(bb.GetColour(), agent.GetColour())
		if bb.GetColour() == agent.GetColour() {
			sameColourCount = sameColourCount + 1
		}
	}
	// fmt.Println(sameColourCount)
	if sameColourCount > (len(fellowBikers) / 2) {
		return utils.Dictatorship
	} else if sameColourCount > (len(fellowBikers) / 3) {
		return utils.Leadership
	} else {
		return utils.Democracy
	}
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
