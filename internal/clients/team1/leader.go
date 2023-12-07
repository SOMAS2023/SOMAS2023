// FUNCTIONS FOR WHEN LEADER

package team1

import (
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

// --------------------LEADER FUNCTIONS------------------
func (bb *Biker1) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	// decides the weights of other peoples votes
	// Leadership democracy
	// takes in proposed action as a parameter
	// only run for the leader after everyone's proposeDirection is run
	// assigns vector of weights to everyone's proposals, 0.5 is neutral

	//consider adding weights for agents with low points
	fellowBikers := bb.GetFellowBikers()
	weights := map[uuid.UUID]float64{}

	for _, agent := range fellowBikers {
		op, ok := bb.opinions[agent.GetID()]
		if !ok {
			weights[agent.GetID()] = 0.5
		} else {
			weights[agent.GetID()] = op.opinion
		}
	}
	return weights
}

//--------------------END OF LEADER FUNCTIONS------------------
