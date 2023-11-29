package voting

import (
	"SOMAS2023/internal/common/utils"
	"errors"
	"math"

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

func SumOfValues(voteMap IVoter) float64 {
	sum := 0.0
	for _, value := range voteMap.GetVotes() {
		sum += value
	}
	return sum
}

// Returns the normalized vote outcome (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func CumulativeDist(voters []IVoter) (map[uuid.UUID]float64, error) {
	if len(voters) == 0 {
		panic("no votes provided")
	}
	// Vote checks for each voter
	aggregateVotes := make(map[uuid.UUID]float64)
	for _, IVoter := range voters {
		if math.Abs(SumOfValues(IVoter)-1.0) > utils.Epsilon {
			return nil, errors.New("distribution doesn't sum to 1")
		}
		votes := IVoter.GetVotes()
		for id, vote := range votes {
			aggregateVotes[id] += vote
		}
	}
	// normalising step for all voters involved
	for agentId, vote := range aggregateVotes {
		aggregateVotes[agentId] = vote / float64(len(voters))
	}
	return aggregateVotes, nil
}

// returns the winner accoring to chosen voting strategy (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func WinnerFromDist(voters []IVoter) uuid.UUID {
	// TODO handle the error
	aggregateVotes, _ := CumulativeDist(voters)

	var randomWinner, winner uuid.UUID
	maxVote := 0.0
	for id, vote := range aggregateVotes {
		randomWinner = id
		if vote > maxVote {
			maxVote = vote
			winner = id
		}
	}

	if winner == uuid.Nil {
		winner = randomWinner
	}

	return winner
}
