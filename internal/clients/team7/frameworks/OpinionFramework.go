package frameworks

import (
	"math"

	"github.com/google/uuid"
)

// OpinionFramework: Formulates the agent's opinion on different lootboxes.

type OpinionFrameworkInputs struct {
	AgentOpinion map[uuid.UUID]float64
	Mindset      float64
	OpinionType  OpinionType
}

type OpinionType int

const (
	AgentOpinions OpinionType = iota
	LootboxOpinions
)

type OpinionFramework struct {
	OpinionAgentWeights map[OpinionType]map[uuid.UUID]float64
}

func NewOpinionFramework(of OpinionFrameworkInputs) *OpinionFramework {
	return &OpinionFramework{
		OpinionAgentWeights: make(map[OpinionType]map[uuid.UUID]float64),
	}
}

func (of *OpinionFramework) GetOpinion(inputs OpinionFrameworkInputs) float64 {
	numOpinions := len(inputs.AgentOpinion)

	agentIds := make([]uuid.UUID, 0)
	for agentId := range inputs.AgentOpinion {
		agentIds = append(agentIds, agentId)
	}

	weights := make([]float64, numOpinions)
	currentWeights, hasWeights := of.OpinionAgentWeights[inputs.OpinionType]
	if !hasWeights {
		// No weights on the matter, initialise all to 1
		initWeights := make(map[uuid.UUID]float64)
		for idx, agentId := range agentIds {
			initWeights[agentId] = 1
			weights[idx] = 1
		}
		of.OpinionAgentWeights[inputs.OpinionType] = initWeights
	} else {
		// Weights exist but not necessarily for all of the agents so get the weights but set to 1 if doesn't exist
		for idx, agentId := range agentIds {
			agentWeight, agentHasWeight := currentWeights[agentId]
			if !agentHasWeight {
				of.OpinionAgentWeights[inputs.OpinionType][agentId] = 1
				agentWeight = 1
			}
			weights[idx] = agentWeight
		}
	}

	currentMindset := inputs.Mindset

	opinions := make([]float64, numOpinions)
	for idx, agentId := range agentIds {
		opinions[idx] = inputs.AgentOpinion[agentId]
	}

	affinities := make([]float64, numOpinions)
	for idx := range affinities {
		affinities[idx] = 1.0 - math.Abs(opinions[idx]-currentMindset)/math.Max(currentMindset, 1.0-currentMindset)
	}

	rowSum := 0.0
	for _, val := range weights {
		rowSum += val
	}

	for idx := range weights {
		weights[idx] = weights[idx] + weights[idx]*affinities[idx]
	}

	for idx := range weights {
		weights[idx] /= rowSum
	}

	currentOpinion := 0.0
	for idx := range weights {
		currentOpinion += weights[idx] * opinions[idx]
	}

	// Update the weights we have stored
	for idx, agentId := range agentIds {
		of.OpinionAgentWeights[inputs.OpinionType][agentId] = weights[idx]
	}

	return currentOpinion
}
