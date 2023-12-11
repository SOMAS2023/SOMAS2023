// QUERY REQUIREMENTS

package team1

import (
	"github.com/google/uuid"
)

// ---------------------SOCIAL FUNCTIONS------------------------
// get reputation value of all other agents
func (bb *Biker1) GetReputation() map[uuid.UUID]float64 {
	reputation := map[uuid.UUID]float64{}
	for agent, opinion := range bb.opinions {
		reputation[agent] = opinion.opinion
	}
	return reputation
}

// query for reputation value of specific agent with UUID
func (bb *Biker1) QueryReputation(agent uuid.UUID) float64 {
	val, ok := bb.opinions[agent]
	if ok {
		return val.opinion
	} else {
		return 0.5
	}
}

// set reputation value of specific agent with UUID
func (bb *Biker1) SetReputation(agent uuid.UUID, reputation float64) {
	bb.opinions[agent] = Opinion{
		effort:   bb.opinions[agent].effort,
		trust:    bb.opinions[agent].trust,
		fairness: bb.opinions[agent].fairness,
		opinion:  reputation,
	}
}

//---------------------END OF SOCIAL FUNCTIONS------------------------
