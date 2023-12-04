package team5Agent

import (
	"SOMAS2023/internal/common/utils"
	"github.com/google/uuid"
	"fmt"
	"math"
)

//for testing use any box in targetLootBoxID

//i added a comment for the printstatement in get colour remember to remove and in the original decideforce fn in base biker remember to remove those

// i think i found it divide the steering force by the amount of agents on the bike

// so this bassically adjusts the force depending on the energy of the agent



func (t5 *team5Agent) DecideForce(targetLootBoxID uuid.UUID) {
	//fmt.Println("testing 1")

	currLocation := t5.GetLocation()
	//fmt.Println("Current Location: ", currLocation)

	nearestLoot := t5.ProposeDirection()
	//fmt.Println("Nearest Loot ID: ", nearestLoot)

	currentLootBoxes := t5.GetGameState().GetLootBoxes()
	//fmt.Println("Number of Loot Boxes: ", len(currentLootBoxes))

	if len(currentLootBoxes) > 0 {
		targetPos := currentLootBoxes[nearestLoot].GetPosition()
		//fmt.Println("Target Position: ", targetPos)

		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		//fmt.Println("Delta X: ", deltaX, "Delta Y: ", deltaY)
		distance_to_loot_box := math.Sqrt((((deltaX)*(deltaX)) + ((deltaY)*(deltaY))))

		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi

		audiPos := t5.GetGameState().GetAudi().GetPosition()

		deltaX2 := audiPos.X - currLocation.X
		deltaY2 := audiPos.Y - currLocation.Y
        
		distance_to_audi := math.Sqrt((((deltaX2)*(deltaX2)) + (deltaY2*(deltaY2))))
		//if the audi is closer than the closest lootbox then run away
		if distance_to_audi < distance_to_loot_box { 
			angle = math.Atan2(-deltaY2, deltaX2)
			normalisedAngle = angle / math.Pi
		}
		if distance_to_audi > distance_to_loot_box { 
			angle = math.Atan2(deltaY, deltaX)
			normalisedAngle = angle / math.Pi
		}
			

		orientation := t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetOrientation()
		//fmt.Println("Bike Orientation: ", orientation) 
		turning_depending_on_agents_on_that_bike := (normalisedAngle - orientation);
		///(float64(len(t5.GetMegaBike())));
		fmt.Println("Angle: ", angle, "Normalized Angle: ", normalisedAngle, " bike orientation ", orientation, "turning_depending_on_agents_on_that_bike ", turning_depending_on_agents_on_that_bike )
		//and i can change this depending on how the enemy agents are turning
		

		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: turning_depending_on_agents_on_that_bike,
		}
		//fmt.Println("Turning Decision: ", turningDecision)

		ownEnergyLevel := t5.GetEnergyLevel()
		// ask the guys what number i want to put the own energy level and if it
		// should be adjusted based on energy level to sva eenergy
		Biker_pedal := utils.BikerMaxForce
		//if our own agents on a bike or just us on a bike we use full force this is only when we are on a bike with other agents or more than 3 agents
		if len(t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetAgents()) > 3 {
			if ownEnergyLevel < 0.2 {
				Biker_pedal = ownEnergyLevel * utils.BikerMaxForce
			}
		}

		Forces_movement := utils.Forces{
			Pedal:   Biker_pedal,
			Brake:   0.0,
			Turning: turningDecision,
		}
		t5.SetForces(Forces_movement)
	}
}

func calculatePedalForceBasedOnEnergy() {
	panic("unimplemented")
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



