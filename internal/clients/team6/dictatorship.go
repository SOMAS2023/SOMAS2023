package team6

import (
	voting "SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

func (bb *Team6Biker) DictateDirection() uuid.UUID {
	return bb.ProposeDirection()
}

func (bb *Team6Biker) DecideKickOut() []uuid.UUID {
	kickoutbikers := []uuid.UUID{}
	fellowBikers := bb.GetFellowBikers()
	lowestTrust := 2.0
	lowestTrustID := uuid.UUID{}

	// Based on our calculation of reputation, kick out the biker(s)[uuid] with lowest reputation
	for _, agent := range fellowBikers {
		bikerID := agent.GetID()
		bikerTrust := bb.QueryReputation(bikerID)
		if bikerID != uuid.Nil {
			if bikerTrust < lowestTrust {
				lowestTrust = bikerTrust
				lowestTrustID = bikerID
			}
		}
	}
	kickoutbikers = append(kickoutbikers, lowestTrustID)
	return kickoutbikers
}

func (bb *Team6Biker) DecideDictatorAllocation() voting.IdVoteMap {
	dictatorID := bb.GetID()
	fellowBikers := bb.GetFellowBikers()
	distribution := make(voting.IdVoteMap)
	dictatorEnergy := 0.9

	restenergy := 1 - dictatorEnergy
	averageenergy := restenergy / float64(len(fellowBikers)-1)

	distribution[dictatorID] = dictatorEnergy
	for _, agent := range fellowBikers {
		agentID := agent.GetID()
		if agentID != dictatorID {
			distribution[agentID] = averageenergy
		}

	}
	return distribution
}
