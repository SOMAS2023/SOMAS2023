package frameworks

import (
	utils "SOMAS2023/internal/common/utils"
	"fmt"
	"math"
)

type NavigationInputs struct {
	DesiredLootbox  utils.Coordinates
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
	// Get distances between current location and desired lootbox
	dx := ndf.inputs.DesiredLootbox.X - ndf.inputs.CurrentLocation.X
	dy := ndf.inputs.DesiredLootbox.Y - ndf.inputs.CurrentLocation.Y

	angle_radians := math.Atan2(dy, dx)

	// Normalize angle to be between -1 and 1
	turning_force := angle_radians / math.Pi
	turning_input := utils.TurningDecision{SteerBike: true, SteeringForce: turning_force}
	pedalling_force := float64(1)
	braking_force := float64(0)

	forces := utils.Forces{Pedal: pedalling_force, Brake: braking_force, Turning: turning_input}

	fmt.Println("NavigationDecisionFramework: GetDecision called")
	fmt.Println("NavigationDecisionFramework: Current location: ", ndf.inputs.CurrentLocation)
	fmt.Println("NavigationDecisionFramework: Desired lootbox: ", ndf.inputs.DesiredLootbox)
	fmt.Println("NavigationDecisionFramework: Forces: ", forces)

	return forces
}

func NewNavigationDecisionFramework() *NavigationDecisionFramework {
	return &NavigationDecisionFramework{
		inputs: &NavigationInputs{},
	}
}
