package team5Agent

import (
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
	"fmt"
	"math"
)

//for testing use any box in targetLootBoxID

//i added a comment for the printstatement in get colour remember to remove and in the original decideforce fn in base biker remember to remove those
func (t5 *team5Agent) DecideForce(targetLootBoxID uuid.UUID) {
	fmt.Println("testing 1")

	currLocation := t5.GetLocation()
	fmt.Println("Current Location: ", currLocation)

	nearestLoot := t5.ProposeDirection()
	fmt.Println("Nearest Loot ID: ", nearestLoot)

	currentLootBoxes := t5.GetGameState().GetLootBoxes()
	fmt.Println("Number of Loot Boxes: ", len(currentLootBoxes))

	if len(currentLootBoxes) > 0 {
		targetPos := currentLootBoxes[nearestLoot].GetPosition()
		fmt.Println("Target Position: ", targetPos)

		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		fmt.Println("Delta X: ", deltaX, "Delta Y: ", deltaY)

		angle := math.Atan2(deltaX, deltaY)
		normalisedAngle := angle / math.Pi
		fmt.Println("Angle: ", angle, "Normalized Angle: ", normalisedAngle)

		orientation := t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetOrientation()
		fmt.Println("Bike Orientation: ", orientation)

		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - orientation,
		}
		fmt.Println("Turning Decision: ", turningDecision)

		nearestBoxForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		fmt.Println("Nearest Box Forces: ", nearestBoxForces)

		t5.SetForces(nearestBoxForces)
	} else {
		idleForces := utils.Forces{
			Pedal:   0.0,
			Brake:   0.0,
			Turning: utils.TurningDecision{SteerBike: false},
		}
		fmt.Println("Idle Forces: ", idleForces)

		t5.SetForces(idleForces)
	}
}

	


// so this bassically adjusts the force depending on the energy of the agent
func (t5 *team5Agent) calculatePedalForceBasedOnEnergy() float64 {
	ownEnergyLevel := t5.GetEnergyLevel()
	// ask the guys what number i want to put the own energy level and if it
	// should be adjusted based on energy level to sva eenergy
	if ownEnergyLevel < 0.5 {
		return ownEnergyLevel * utils.BikerMaxForce
	}
	return utils.BikerMaxForce
}

// here I can add implementation of stategy like:
//___________________________________________________________________________________________________________________________

// func (t5 *Team5Biker) calculateAverageEnergyOfBikeagents() float64 {
//     bikeagents := t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetAgents()
//     var totalEnergy float64
//     var count float64

//     for _, agents := range bikeagents {
//         if agents.GetID() != t5.GetID() { // Exclude self
//             totalEnergy += agents.GetEnergyLevel()
//             count++
//         }
//     }

//     if count == 0 {
//         return 1 // If no other mates, return full energy or none and just bike hop maybe? see what others think
//     }
//     return totalEnergy / count
// }
//___________________________________________________________________________________________________________________________

// can add this to decide force

// add a function depends


// 2.)speed of other bikes
// and 
// position of other bikes and how fast to peddle depending on that

// 1.)so the lootbox is the direction but we may need to turn more if the bike doesnt turn enough.



//have a meeting with others discuss what other fns i can implement nd what helps others
// runs no issues