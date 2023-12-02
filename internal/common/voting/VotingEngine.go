package voting

import (
	"SOMAS2023/internal/common/utils"
	"errors"
	"math"
	"sort"

	"github.com/google/uuid"
)

type GovernanceVote map[utils.Governance]float64

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
func GetAcceptanceRanking(rankings []map[uuid.UUID]bool) []uuid.UUID {
	// sum the number of acceptance rankings for all the agents
	cumulativeRank := make(map[uuid.UUID]int)
	quorum := len(rankings) / 2
	for _, ranking := range rankings {
		for agent, outcome := range ranking {
			val, ok := cumulativeRank[agent]
			if outcome && ok {
				cumulativeRank[agent] = val + 1
			} else if outcome {
				cumulativeRank[agent] = 1
			}
		}
	}
	passedUnsorted := make(map[uuid.UUID]int)
	for agent, val := range cumulativeRank {
		if val > quorum {
			passedUnsorted[agent] = val
		}
	}

	// sort according to ranking
	unsortedAcceptedList := make([]uuid.UUID, len(passedUnsorted))
	i := 0
	for key, _ := range passedUnsorted {
		unsortedAcceptedList[i] = key
		i += 1
	}
	sort.Slice(unsortedAcceptedList, func(i, j int) bool {
		return passedUnsorted[unsortedAcceptedList[i]] > passedUnsorted[unsortedAcceptedList[j]]
	})
	return unsortedAcceptedList
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

// return the votesMap
func GetVotesMap(voters map[uuid.UUID]IVoter) (map[uuid.UUID]map[uuid.UUID]float64, error) {
	if len(voters) == 0 {
		panic("no votes provided")
	}
	// Vote checks for each voter
	VotesOfAgents := make(map[uuid.UUID]map[uuid.UUID]float64)
	for agentID, IVoter := range voters {
		if math.Abs(SumOfValues(IVoter)-1.0) > utils.Epsilon {
			return nil, errors.New("distribution doesn't sum to 1")
		}
		votes := IVoter.GetVotes()
		VotesOfAgents[agentID] = votes
	}

	return VotesOfAgents, nil
}

// returns the winner accoring to chosen voting strategy (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func WinnerFromDist(voters map[uuid.UUID]IVoter, voteWeight map[uuid.UUID]float64) uuid.UUID {
	// TODO handle the error
	VotesOfAgents, _ := GetVotesMap(voters)
	var winner uuid.UUID
	switch utils.VoteAction {
	case utils.PLURALITY:
		winner = Plurality(VotesOfAgents, voteWeight)
	case utils.RUNOFF:
		winner = Runoff(VotesOfAgents, voteWeight)
	case utils.BORDACOUNT:
		winner = BordaCount(VotesOfAgents, voteWeight)
	case utils.INSTANTRUNOFF:
		winner = InstantRunoff(VotesOfAgents, voteWeight)
	case utils.APPROVAL:
		winner = Approval(VotesOfAgents, voteWeight)
	case utils.COPELANDSCORING:
		winner = CopelandScoring(VotesOfAgents, voteWeight)
	}

	return winner
}

func WinnerFromGovernance(voters []GovernanceVote) (utils.Governance, error) {
	// check if length of votes is greater than one
	if len(voters) == 0 {
		return utils.Invalid, errors.New("no votes provided")
	}

	// Summing up the votes for each governance type
	for _, vote := range voters {
		sum := 0.0
		for _, votes := range vote {
			sum += votes
		}
		if sum > 1.0 {
			return utils.Invalid, errors.New("distribution doesn't sum to 1")
		}
	}

	var voteTotals = make(map[utils.Governance]float64)
	var winner utils.Governance
	var highestVotes float64

	// Summing up the votes for each governance type
	for _, vote := range voters {
		for governance, votes := range vote {
			voteTotals[governance] += votes
		}
	}
	// Finding the governance type with the highest votes
	for governance, votes := range voteTotals {
		if votes > highestVotes {
			highestVotes = votes
			winner = governance
		}
	}

	return winner, nil
}
