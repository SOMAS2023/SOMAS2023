package modules

import (
	"fmt"
	"math"

	"github.com/google/uuid"
)

type SocialCapital struct {
	forgivenessCounter map[uuid.UUID]int
	SocialCapital      map[uuid.UUID]float64
	Reputation         map[uuid.UUID]float64
	Institution        map[uuid.UUID]float64
	SocialNetwork      map[uuid.UUID]float64
}

func (sc *SocialCapital) GetAverage(scComponent map[uuid.UUID]float64) float64 {
	// Prevent divide
	if len(scComponent) == 0 {
		return 0.5
	}
	var sum = 0.0
	for _, value := range scComponent {
		sum += value
	}
	return sum / float64(len(scComponent))
}

func (sc *SocialCapital) GetSum(scComponent map[uuid.UUID]float64) float64 {
	var sum = 0.0
	for _, value := range scComponent {
		sum += value
	}
	return sum
}

func (sc *SocialCapital) GetMinimumSocialCapital() (uuid.UUID, float64) {
	min := math.MaxFloat64
	minAgentId := uuid.Nil
	for agentId, value := range sc.Reputation {
		if sc.SocialCapital[agentId] < min {
			min = value
			minAgentId = agentId
		}
	}
	return minAgentId, min
}

func (sc *SocialCapital) GetMaximumSocialCapital() (uuid.UUID, float64) {
	max := 0.0
	maxAgentId := uuid.Nil
	for agentId, value := range sc.SocialCapital {
		if sc.SocialCapital[agentId] > max {
			max = value
			maxAgentId = agentId
		}
	}
	return maxAgentId, max
}

func (sc *SocialCapital) ClipValues(input float64) float64 {
	value := input
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return value
}

func (sc *SocialCapital) UpdateValue(agentId uuid.UUID, eventValue float64, eventWeight float64, scComponent map[uuid.UUID]float64) {
	_, exists := scComponent[agentId]
	if !exists {
		scComponent[agentId] = sc.GetAverage(scComponent)
	}

	scComponent[agentId] = sc.ClipValues(scComponent[agentId] + eventValue*eventWeight)
}

func (sc *SocialCapital) UpdateReputation(agentId uuid.UUID, eventValue float64, eventWeight float64) {
	sc.UpdateValue(agentId, eventValue, eventWeight, sc.Reputation)
}

func (sc *SocialCapital) UpdateInstitution(agentId uuid.UUID, eventValue float64, eventWeight float64) {
	sc.UpdateValue(agentId, eventValue, eventWeight, sc.Institution)
}

func (sc *SocialCapital) UpdateSocialNetwork(agentId uuid.UUID, eventValue float64, eventWeight float64) {
	sc.UpdateValue(agentId, eventValue, eventWeight, sc.SocialNetwork)
}

// Must be called once every round.
func (sc *SocialCapital) UpdateSocialCapital() {
	fmt.Printf("[UpdateSocialCapital] Social Capital Before: %v\n", sc.SocialCapital)

	for id := range sc.SocialNetwork { // Assumes all maps have the same keys.
		// Add to Forgiveness Counters.
		if _, ok := sc.forgivenessCounter[id]; !ok {
			sc.forgivenessCounter[id] = 0.0
		}

		// Update Forgiveness Counter.
		newSocialCapital := ReputationWeight*sc.Reputation[id] + InstitutionWeight*sc.Institution[id] + NetworkWeight*sc.SocialNetwork[id]

		if sc.SocialCapital[id] < newSocialCapital {
			sc.forgivenessCounter[id] = 0
		}

		if sc.SocialCapital[id] > newSocialCapital && sc.forgivenessCounter[id] <= 3 {
			// Forgive if forgiveness counter is less than 3 and new social capital is less.
			sc.forgivenessCounter[id]++
			sc.SocialCapital[id] = newSocialCapital + forgivenessFactor*(sc.SocialCapital[id]-newSocialCapital)
		} else {
			sc.SocialCapital[id] = newSocialCapital
		}
	}
	fmt.Printf("[UpdateSocialCapital] Social Capital After: %v\n", sc.SocialCapital)
}

func NewSocialCapital() *SocialCapital {
	return &SocialCapital{
		forgivenessCounter: make(map[uuid.UUID]int),
		SocialCapital:      make(map[uuid.UUID]float64),
		Reputation:         make(map[uuid.UUID]float64),
		Institution:        make(map[uuid.UUID]float64),
		SocialNetwork:      make(map[uuid.UUID]float64),
	}
}
