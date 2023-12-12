package frameworks

import (
	voting "SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

// This map can hold any type of data as the value
type Map map[uuid.UUID]interface{}

// This map can be used for votes where we return a agent UUIDs mapped to boolean.
type MapIdBool map[uuid.UUID]bool

// Type for scoring different votes
// High value for variables of this type expresses being in favour of vote.
type ScoreType float64

type VoteOnAgentsInput struct {
	AgentCandidates      []uuid.UUID
	CurrentSocialNetwork map[uuid.UUID]*SocialConnection
}

type VoteOnLootBoxesInput struct {
	LootBoxCandidates map[uuid.UUID]uuid.UUID
	MyPersonality     *Personality
	MyDesired         uuid.UUID
	MyOpinion         map[uuid.UUID]float64
}

// Expected to return votes which sum to 1 for some voting types.
// This function normalises vote map to sum to 1.
func NormaliseVote(agentScoreMap voting.IdVoteMap, totalScore float64) voting.IdVoteMap {
	normalisedVoteMap := make(voting.IdVoteMap)
	// Find the sum of all the scores.
	for agentId, agentScore := range agentScoreMap {
		normalisedVoteMap[agentId] = agentScore / totalScore
	}

	return normalisedVoteMap
}
