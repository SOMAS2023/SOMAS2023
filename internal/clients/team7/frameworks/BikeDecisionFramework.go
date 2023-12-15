package frameworks

import (
	objects "SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

// BikeDecisionFramework: Decides whether agent should stay on or leave current bike.
// Depends on whether our bike has moved towards any of our lootbox proposals in previous iterations.

// There will be different times when a decision needs to be made.
// For example, when the biker is on the bike, the biker needs to decide whether to get off the bike or not.
// This enum can be used to specify the type of decision that needs to be made.
type BikeDecisionType int

const (
	StayOrLeaveBike BikeDecisionType = iota
	FindNewBike
)

// All inputs needed for the decisions around staying on a bike or joining a new bike.
// Not all inputs will be needed for all decisions.
type BikeDecisionInputs struct {
	DecisionType                   BikeDecisionType                // Type of decision that needs to be made
	CurrentLocation                utils.Coordinates               // Current location of the biker
	AvailableBikes                 map[uuid.UUID]objects.IMegaBike // Map of available bikes
	PreviousDistancesFromProposals []float64
}

type BikeDecision struct {
	LeaveBike bool
}

type BikeDecisionFramework struct {
	IDecisionFramework[BikeDecisionInputs, BikeDecision]
}

func (bdf *BikeDecisionFramework) GetDecision(inputs BikeDecisionInputs) BikeDecision {
	if inputs.DecisionType == StayOrLeaveBike {
		distances := inputs.PreviousDistancesFromProposals
		// find gradient of distances from previous lootbox proposals
		if len(distances) == 0 {
			return BikeDecision{LeaveBike: false}
		}

		gradient := make([]float64, len(distances)-1)
		for i := 0; i < len(distances)-1; i++ {
			gradient[i] = distances[i+1] - distances[i]
		}

		// Check if last 3 gradient values are positive
		// If so, leave bike
		if len(gradient) >= 3 {
			if gradient[len(gradient)-1] > 0 && gradient[len(gradient)-2] > 0 && gradient[len(gradient)-3] > 0 {
				return BikeDecision{LeaveBike: true}
			}
		} else {
			return BikeDecision{LeaveBike: false}
		}
	}

	return BikeDecision{}
}

func NewBikeDecisionFramework() *BikeDecisionFramework {
	return &BikeDecisionFramework{}
}
