package voting

import (
	"SOMAS2023/internal/common/utils"
	"errors"
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
func GetAcceptanceRanking(rankings map[uuid.UUID]map[uuid.UUID]bool, weights map[uuid.UUID]float64) []uuid.UUID {
	// sum the number of acceptance rankings for all the agents
	cumulativeRank := make(map[uuid.UUID]float64)
	quorum := float64(len(rankings)) / 2.0
	for voter, ranking := range rankings {
		for agent, outcome := range ranking {
			val, ok := cumulativeRank[agent]
			if outcome && ok {
				cumulativeRank[agent] = val + weights[voter]
			} else if outcome {
				cumulativeRank[agent] = 1.0
			}
		}
	}
	passedUnsorted := make(map[uuid.UUID]float64)
	for agent, val := range cumulativeRank {
		if val > quorum {
			passedUnsorted[agent] = val
		}
	}

	// sort according to ranking
	unsortedAcceptedList := make([]uuid.UUID, len(passedUnsorted))
	i := 0
	for key := range passedUnsorted {
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
func CumulativeDist(voters map[uuid.UUID]IVoter, weights map[uuid.UUID]float64) map[uuid.UUID]float64 {
	if len(voters) == 0 {
		panic("no votes provided")
	}
	// Vote checks for each voter
	aggregateVotes := make(map[uuid.UUID]float64)

	// initialise votes to 0.0
	for voter := range voters {
		aggregateVotes[voter] = 0.0
	}

	for agentID, voter := range voters {
		voteSum := SumOfValues(voter)
		votes := voter.GetVotes()
		weight := weights[agentID]
		for id, vote := range votes {
			aggregateVotes[id] += weight * vote / voteSum
		}
	}

	normalizeFactor := 0.0
	for _, vote := range aggregateVotes {
		normalizeFactor += vote
	}
	if normalizeFactor == 0.0 {
		panic("all votes summed to zero")
	}
	// normalising step for all voters involved
	for agentId, vote := range aggregateVotes {
		aggregateVotes[agentId] = vote / normalizeFactor
	}
	return aggregateVotes
}

// return the votesMap
func GetVotesMap(voters map[uuid.UUID]IVoter) map[uuid.UUID]map[uuid.UUID]float64 {
	if len(voters) == 0 {
		panic("no votes provided")
	}
	// Vote checks for each voter
	VotesOfAgents := make(map[uuid.UUID]map[uuid.UUID]float64)
	for agentID, voter := range voters {
		voteSum := SumOfValues(voter)
		votes := voter.GetVotes()
		for id := range votes {
			if id == uuid.Nil {
				panic("agent voted for a nil uuid")
			}
			votes[id] /= voteSum
		}
		VotesOfAgents[agentID] = votes
	}

	return VotesOfAgents
}

// returns the winner accoring to chosen voting strategy (assumes all the maps contain a voting between 0-1
// for each option, and that all the votings sum to 1)
func WinnerFromDist(voters map[uuid.UUID]IVoter, voteWeight map[uuid.UUID]float64) uuid.UUID {
	VotesOfAgents := GetVotesMap(voters)
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
	// TODO call group 8 voting function
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

// Need to check if the input param is expecting a vote that is just one governance type
func TallyFoundingVotes(voters map[uuid.UUID]utils.Governance) (map[utils.Governance]int, error) {
	// check if length of votes is greater than one
	if len(voters) == 0 {
		return nil, errors.New("no votes provided")
	}

	// Summing up the votes for each governance type
	aggregateFoundingTotals := make(map[utils.Governance]int)

	// Get the governance type for each agent
	for _, vote := range voters {
		// Add to the tally for each governance type
		if val, ok := aggregateFoundingTotals[vote]; ok {
			aggregateFoundingTotals[vote] = val + 1
		} else {
			aggregateFoundingTotals[vote] = 1
		}
	}

	return aggregateFoundingTotals, nil
}
