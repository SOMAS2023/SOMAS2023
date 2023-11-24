package voting

import (
	"math"
	"sort"

	"github.com/google/uuid"
)

// Auxiliary Structure for Sorting
type kv struct {
	Key   uuid.UUID
	Value float64
}

func Plurality(voteList []map[uuid.UUID]float64) uuid.UUID {
	/*
		Plurality:
			Each voter selects one candidate and the candidate with the most first-placed votes is the winner.
	*/
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	for _, preference := range voteList {
		var maxPreference float64
		var firstLootBoxChoice uuid.UUID
		for lootBox, key := range preference {
			if key > maxPreference {
				firstLootBoxChoice = lootBox
				maxPreference = key
			}
		}
		voteCount[firstLootBoxChoice]++
	}

	// final step: we need to find the winner with highest count number in map.
	var maxVotes float64

	for lootBox, votes := range voteCount {
		if votes > maxVotes {
			maxVotes = votes
			winner = lootBox
		}
	}

	// return the final winner
	return winner
}

func Runoff(voteList []map[uuid.UUID]float64) uuid.UUID {
	/*
		Runoff:
			1st round: 	each voter selects one candidate, and the two candidates with most first-placed votes are identified.
						If either already has a majority, this candidate is declared the winner.
			2nd round: 	each voter selects one candidate, the candidate with most votes now is the winner.
	*/
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	// ----- first round -----
	// find the count number of each lootbox
	for _, preference := range voteList {
		var maxPreference float64
		var firstLootBoxChoice uuid.UUID
		for lootBox, key := range preference {
			if key > maxPreference {
				firstLootBoxChoice = lootBox
				maxPreference = key
			}
		}
		voteCount[firstLootBoxChoice]++
	}

	// find the two candidates with most first-placed votes
	var maxVotes1, maxVotes2 float64
	var winner1, winner2 uuid.UUID
	for lootBox, votes := range voteCount {
		if votes > maxVotes1 {
			winner2 = winner1
			maxVotes2 = maxVotes1
			winner1 = lootBox
			maxVotes1 = votes
		} else if votes > maxVotes2 {
			winner2 = lootBox
			maxVotes2 = votes
		}
	}

	// check if either already has a majority or we need the second round
	if maxVotes1 >= (maxVotes2 * 2) {
		// return the majority lootbox
		return winner1
	} else {
		// ----- second round -----
		voteCount := make(map[uuid.UUID]int)
		for _, preference := range voteList {
			if preference[winner1] > preference[winner2] {
				voteCount[winner1]++
			} else {
				voteCount[winner2]++
			}
		}
		if voteCount[winner1] > voteCount[winner2] {
			winner = winner1
		} else {
			winner = winner2
		}
	}

	return winner
}

func BordaCount(voteList []map[uuid.UUID]float64) uuid.UUID {
	/*
		BordaCount:
			Each voter rank order all the candidates. With n candidates being ranked k scores (n-k)+1 Borda points.
			The candidate with the highest Borda Score is the winner
	*/
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	// initialise the map with all candidates
	for _, preference := range voteList {
		for key := range preference {
			voteCount[key] = 0
		}
	}

	// covert the unodered map into ordered list
	var ss [][]kv
	for _, preference := range voteList {
		var s []kv
		for k, v := range preference {
			// ignore the lootbox if value is 0
			if v != 0 {
				s = append(s, kv{k, v})
			}
		}
		// sort the list using preference value of each lootbox
		sort.Slice(s, func(i, j int) bool {
			// in the order from large to small
			return s[i].Value > s[j].Value
		})
		ss = append(ss, s)
	}

	// calculate the Borda score for each candidates
	for _, sortedList := range ss {
		usedKeys := make(map[uuid.UUID]bool)
		for i, kv := range sortedList {
			score := float64(len(voteCount)) - float64(i) + 1
			voteCount[kv.Key] += score
			usedKeys[kv.Key] = true
		}

		// points shared if not explicity ranked
		remainingKeyNumber := float64(len(voteCount)) - float64(len(sortedList))
		remainingScore := (1 + remainingKeyNumber) * remainingKeyNumber / 2
		for key := range voteCount {
			if !usedKeys[key] {
				voteCount[key] += remainingScore / remainingKeyNumber
			}
		}
	}

	// find the winner with highest score
	var maxScore float64
	for key, value := range voteCount {
		if value > maxScore {
			winner = key
			maxScore = value
		}
	}

	return winner
}

func InstantRunoff(voteList []map[uuid.UUID]float64) uuid.UUID {
	/*
		InstantRunoff:
			Each voter rank orders all candidates, and the candidate with the least number of first-place votes is eliminate.
			This is repeated until only one candidate remains
	*/
	voteCount := make(map[uuid.UUID]float64)
	eliminateVote := make(map[uuid.UUID]bool)
	var winner uuid.UUID

	// initialise the map with all candidates
	for _, preference := range voteList {
		for key := range preference {
			voteCount[key] = 0
		}
	}

	// loop to eliminate the least number of first-place votes
	for len(voteCount) > 1 {
		// reset map with value = 0
		for key := range voteCount {
			voteCount[key] = 0
		}

		// count the number of first-place votes for each lootbox
		for _, preference := range voteList {
			var maxScore float64
			var firstLootBoxChoice uuid.UUID
			for key, value := range preference {
				if (value > maxScore) && !eliminateVote[key] {
					maxScore = value
					firstLootBoxChoice = key
				}
			}
			voteCount[firstLootBoxChoice]++
		}

		// eliminate the lootbox with least votes
		var minVotes float64 = math.MaxFloat64
		var candidateToEliminate uuid.UUID
		for key, value := range voteCount {
			if value < minVotes {
				minVotes = value
				candidateToEliminate = key
			}
		}
		eliminateVote[candidateToEliminate] = true
		delete(voteCount, candidateToEliminate)
	}

	// get the final winner
	for key := range voteCount {
		winner = key
	}

	return winner
}

func Approval(voteList []map[uuid.UUID]float64) uuid.UUID {
	/*
		Approval:
			A ballot represents not a linear rank order of decreasing preference,
			but rather represents the set of candidates who are 'equally acceptable' to the voter
	*/
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	for _, preference := range voteList {
		for key, value := range preference {
			if value > 0 {
				voteCount[key]++
			}
		}
	}

	// find the lootbox with
	var maxVotes float64

	for lootBox, votes := range voteCount {
		if votes > maxVotes {
			maxVotes = votes
			winner = lootBox
		}
	}

	return winner
}

func CopelandScoring(voteList []map[uuid.UUID]float64) uuid.UUID {
	/*
		CopelandScoring:
			Each voter submits a ballot with a linear rank order.
			A win-loss record, the Copeland Score, is calculated for each candidate.
	*/

	// the map to store the winning score for each lootbox
	scores := make(map[uuid.UUID]int)

	// iterate the voting
	for _, vote := range voteList {
		for candidate1, score1 := range vote {
			for candidate2, score2 := range vote {
				// do not compare with itself
				if candidate1 == candidate2 {
					continue
				}

				// update the score of each lootbox
				if score1 > score2 {
					scores[candidate1]++
					scores[candidate2]--
				} else if score1 < score2 {
					scores[candidate1]--
					scores[candidate2]++
				}
			}
		}
	}

	// find the lootbox with the highest score
	var maxScore int
	var maxCandidate uuid.UUID
	for candidate, score := range scores {
		if score > maxScore || maxCandidate == uuid.Nil {
			maxScore = score
			maxCandidate = candidate
		}
	}

	return maxCandidate
}
