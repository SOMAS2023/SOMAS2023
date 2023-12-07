package frameworks

import (
	utils "SOMAS2023/internal/common/utils"
	"math"
)

type NavigationInputs struct {
	IsDestination   bool
	Destination     utils.Coordinates
	CurrentLocation utils.Coordinates
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
	ndf.inputs = &inputs

	turningForce := ndf.GetTurnAngle(inputs)
	turningInput := utils.TurningDecision{SteerBike: true, SteeringForce: turningForce}
	pedallingForce := float64(1)
	brakingForce := float64(0)

	forces := utils.Forces{Pedal: pedallingForce, Brake: brakingForce, Turning: turningInput}

	//** fmt.Println("NavigationDecisionFramework: GetDecision called")
	//** fmt.Println("NavigationDecisionFramework: Current location: ", ndf.inputs.CurrentLocation)
	//** fmt.Println("NavigationDecisionFramework: Desired lootbox: ", ndf.inputs.Destination)
	//** fmt.Println("NavigationDecisionFramework: Forces: ", forces)

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
