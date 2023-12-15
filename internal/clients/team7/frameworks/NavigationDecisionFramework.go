package frameworks

import (
	utils "SOMAS2023/internal/common/utils"
	"math"
)

type NavigationInputs struct {
	IsDestination          bool
	Destination            utils.Coordinates
	CurrentLocation        utils.Coordinates
	CurrentEnergy          float64
	ConscientiousnessLevel float64
}

/*
NavigationDecisionFramework: Decides steering, braking and pedalling strategy of agent.

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
	// Low conscientiousness => lazy => Pedal with low effort relative to energy level
	// High conscientiousness => hard-working => Pedal with high effort relative to energy level.
	forceToEnergyRatio := inputs.ConscientiousnessLevel // This value can be adjusted

	// Pedaling force is the current energy multiplied by a constant ratio
	// Set pedalling to zero if current energy is less than 0.2
	var pedallingForce float64
	if inputs.CurrentEnergy < 0.2 {
		pedallingForce = 0
	} else {
		pedallingForce = inputs.CurrentEnergy * forceToEnergyRatio
	}

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
