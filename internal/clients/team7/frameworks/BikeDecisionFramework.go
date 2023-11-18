package frameworks

import (
	objects "SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"fmt"
)

type BikeDecisionType int

const (
	OnBike BikeDecisionType = iota
	OffBike
)

type BikeDecisionInputs struct {
	decisionType    BikeDecisionType
	currentLocation utils.Coordinates
	availableBikes  []objects.MegaBike
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
