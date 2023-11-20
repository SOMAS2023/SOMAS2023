package voting

import (
	"github.com/google/uuid"
)

// Generic IVoter type to accept different outputs
type IVoter interface {
	GetVotes() map[uuid.UUID]float64
}

// lootboxID:distribution
type LootboxVoteMap map[uuid.UUID]float64

// LootboxVoteMap already has the required structure, so we just add a method to satisfy the IVoter interface.
func (lvm LootboxVoteMap) GetVotes() map[uuid.UUID]float64 {
	return lvm
}

// BikerID:distribution
type IdVoteMap map[uuid.UUID]float64

func (ivm IdVoteMap) GetVotes() map[uuid.UUID]float64 {
	return ivm
}

// this function will take in a list of maps from ids to their corresponding vote (yes/ no in the case of acceptance)
// and retunr a list of ids that can be accepted according to some metric (ie more than half voted yes)
// ranked according to a metric (ie overall number of yes's)
func GetAcceptanceRanking([]map[uuid.UUID]bool) []uuid.UUID {
	// TODO implement
	panic("not implemented")
	// return make([]uuid.UUID, 0)
}

// Returns the normalized vote outcome (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func CumulativeDist(voters []IVoter) map[uuid.UUID]float64 {
	if len(voters) == 0 {
		panic("no votes provided")
	}

	aggregateVotes := make(map[uuid.UUID]float64)
	for _, IVoter := range voters {
		for id, vote := range IVoter.GetVotes() {
			aggregateVotes[id] += vote
		}
	}

	return aggregateVotes
}

// returns the winner accoring to chosen voting strategy (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func WinnerFromDist(voters []IVoter) uuid.UUID {
	aggregateVotes := CumulativeDist(voters)

	var winner uuid.UUID
	var maxVote float64
	for id, vote := range aggregateVotes {
		if vote > maxVote {
			maxVote = vote
			winner = id
		}
	}

	if winner == uuid.Nil {
		panic("no winner found")
	}

	return winner
}
