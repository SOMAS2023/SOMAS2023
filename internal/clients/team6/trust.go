package team6

import (
	"github.com/google/uuid"
)

type Trust struct {
	reputations     float64
	goal            float64
	relativeSuccess float64
	trust           float64
}

func (bb *Team6Biker) UpdateAgentgoal(agentID uuid.UUID) {

	newTrust := Trust{
		reputations:     bb.Trust[agentID].reputations,
		goal:            bb.Trust[agentID].goal + 1,
		relativeSuccess: bb.Trust[agentID].relativeSuccess,
		trust:           bb.Trust[agentID].trust,
	}
	bb.Trust[agentID] = newTrust
}

func (bb *Team6Biker) UpdateRelativeSuccess(agentID uuid.UUID) {

	newTrust := Trust{
		reputations:     bb.Trust[agentID].reputations,
		goal:            bb.Trust[agentID].goal,
		relativeSuccess: bb.Trust[agentID].relativeSuccess,
		trust:           bb.Trust[agentID].trust,
	}
	bb.Trust[agentID] = newTrust
}

// Update a certain agent's Trust
func (bb *Team6Biker) UpdateTrust(agent uuid.UUID) {
	_, ok := bb.Trust[agent]
	if !ok {
		//if we have no data on an agent, initialise to neutral
		newTrust := Trust{
			reputations:     1,
			goal:            0,
			relativeSuccess: 0,
			trust:           1,
		}
		bb.Trust[agent] = newTrust
	}

	newTrust := Trust{
		reputations:     bb.Trust[agent].reputations,
		goal:            bb.Trust[agent].goal,
		relativeSuccess: bb.Trust[agent].relativeSuccess,
		trust:           bb.Trust[agent].relativeSuccess * (bb.Trust[agent].goal + bb.Trust[agent].relativeSuccess),
	}

	bb.Trust[agent] = newTrust
}

// Initialise
func (bb *Team6Biker) setTrust() {

	//bb.Trust = make(map[uuid.UUID]Trust)
	for _, agent := range bb.GetGameState().GetAgents() {
		//if (bb.GetID()!=agent.GetID()){
		agentId := agent.GetID()
		bb.UpdateTrust(agentId)
		//}
	}
}
