package frameworks

import (
	objects "SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"fmt"
)

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
	decisionType    BikeDecisionType   // Type of decision that needs to be made
	currentLocation utils.Coordinates  // Current location of the biker
	availableBikes  []objects.MegaBike // List of available bikes
}

type BikeDecision struct {
	leaveBike bool
}

type BikeDecisionFramework struct {
	IDecisionFramework[BikeDecisionInputs, BikeDecision]
	inputs *BikeDecisionInputs
}

func (bdf *BikeDecisionFramework) GetDecision(inputs BikeDecisionInputs) BikeDecision {
	bdf.inputs = &inputs
	fmt.Println("BikeDecisionFramework: GetDecision called")
	fmt.Println("BikeDecisionFramework: Current location: ", bdf.inputs.currentLocation)
	fmt.Println("BikeDecisionFramework: Available bikes: ", bdf.inputs.availableBikes)
	fmt.Println("BikeDecisionFramework: Decision type: ", bdf.inputs.decisionType)

	return BikeDecision{leaveBike: false}
}

func NewBikeDecisionFramework() *BikeDecisionFramework {
	return &BikeDecisionFramework{
		inputs: &BikeDecisionInputs{},
	}
}
