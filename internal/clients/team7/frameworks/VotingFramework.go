package frameworks

import (
	"github.com/google/uuid"
	//voting "SOMAS2023/internal/common/voting"
)

// This map can hold any type of data as the value
type Map map[uuid.UUID]interface{}

// This map can be used for votes where we return a agent UUIDs mapped to boolean.
type MapIdBool map[uuid.UUID]bool

// Type for scoring different votes
// High value for variables of this type expresses being in favour of vote.
type ScoreType float64

type VoteOnAgentsInput struct {
	AgentCandidates []uuid.UUID
}
