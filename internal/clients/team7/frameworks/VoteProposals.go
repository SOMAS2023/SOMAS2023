package frameworks

/*
import (
	"fmt"
	"github.com/google/uuid"
)

// Ally needs the infrastructure to be updated for this to work
func VoteOnProposalsWrapper(voteInputs VoteInputs) Vote {
	var vote Vote
	switch voteInputs.VoteParameters {
	case Proportion:
		vote = Proportions(voteInputs)
	case YesNo:
		vote = YesNos(voteInputs)
	default:
		//** fmt.Println("New decision type!")
		vote = Proportions(voteInputs)
	}

	return vote
}

// TODO: Add functions for voting on which loot box to go to.


func Proportions(voteInputs VoteInputs) Vote {
	var votes Map
	var candidates []uuid.UUID
	biker *BaseTeamSevenBiker
	candidates = voteInputs.Candidates
	totOptions := len(candidates)
	normalDist := 1.0 / float64(totOptions)
	for _, proposal := range candidates {
		if(proposal = biker.NearestLoot())
		votes[proposal] = normalDist
	}
	return votes

}
*/
