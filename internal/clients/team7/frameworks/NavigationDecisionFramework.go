package frameworks

import (
	utils "SOMAS2023/internal/common/utils"
	"fmt"
	"math"
)

type NavigationInputs struct {
	desiredLootbox  utils.Coordinates
	currentLocation utils.Coordinates
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
	dx := ndf.inputs.desiredLootbox.X - ndf.inputs.currentLocation.X
	dy := ndf.inputs.desiredLootbox.Y - ndf.inputs.currentLocation.Y

	angle_radians := math.Atan2(dy, dx)

	// Normalize angle to be between -1 and 1
	turning_force := angle_radians / math.Pi
	pedalling_force := float64(1)
	braking_force := float64(0)

	forces := utils.Forces{Pedal: pedalling_force, Brake: braking_force, Turning: turning_force}

	fmt.Println("NavigationDecisionFramework: GetDecision called")
	fmt.Println("NavigationDecisionFramework: Current location: ", ndf.inputs.currentLocation)
	fmt.Println("NavigationDecisionFramework: Desired lootbox: ", ndf.inputs.desiredLootbox)
	fmt.Println("NavigationDecisionFramework: Forces: ", forces)

	return forces
}

func NewNavigationDecisionFramework() *NavigationDecisionFramework {
	return &NavigationDecisionFramework{
		inputs: &NavigationInputs{},
	}
}
