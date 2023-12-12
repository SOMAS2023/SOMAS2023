package frameworks

import (
	utils "SOMAS2023/internal/common/utils"
	"math"
)

type NavigationInputs struct {
	IsDestination   bool
	Destination     utils.Coordinates
	CurrentLocation utils.Coordinates
	CurrentEnergy   float64
}

/*
This framework can be used for determining the navigation decisions of the biker.

	Input: NavigationInputs
	Output: Forces
*/
type NavigationDecisionFramework struct {
	IDecisionFramework[NavigationInputs, utils.Forces]
	inputs *NavigationInputs
}

func (ndf *NavigationDecisionFramework) GetDecision(inputs NavigationInputs) utils.Forces {
	// Store the inputs in the navigation decision framework for later use
	ndf.inputs = &inputs

	// Calculate the turning force based on the inputs
	turningForce := ndf.GetTurnAngle(inputs)
	// Create a turning decision with the calculated turning force
	turningInput := utils.TurningDecision{SteerBike: true, SteeringForce: turningForce}

	// Constant ratio for pedaling force based on the current energy level
	const PedalingEfficiencyConstant float64 = 0.8 // This value can be adjusted

	// Pedaling force is the current energy multiplied by a constant ratio
	pedallingForce := inputs.CurrentEnergy * PedalingEfficiencyConstant

	// Braking force is set to zero, assuming no need to brake in this context
	brakingForce := float64(0)

	// Combine all the forces into a single structure to be returned
	forces := utils.Forces{Pedal: pedallingForce, Brake: brakingForce, Turning: turningInput}

	// Return the combined forces as the decision
	return forces
}

func (ndf *NavigationDecisionFramework) GetTurnAngle(inputs NavigationInputs) float64 {
	if !inputs.IsDestination {
		return 0
	}
	// Get distances between current location and desired lootbox
	dx := ndf.inputs.Destination.X - ndf.inputs.CurrentLocation.X
	dy := ndf.inputs.Destination.Y - ndf.inputs.CurrentLocation.Y

	angleRadians := math.Atan2(dy, dx)

	// Normalize angle to be between -1 and 1
	turningForce := angleRadians / math.Pi
	return turningForce
}

func NewNavigationDecisionFramework() *NavigationDecisionFramework {
	return &NavigationDecisionFramework{
		inputs: &NavigationInputs{},
	}
}
