package team5Agent

import (
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

//for testing use any box in targetLootBoxID

//i added a comment for the printstatement in get colour remember to remove and in the original decideforce fn in base biker remember to remove those

// i think i found it divide the steering force by the amount of agents on the bike

// so this bassically adjusts the force depending on the energy of the agent

func (t5 *team5Agent) DecideForce(targetLootBoxID uuid.UUID) {
	//fmt.Println("testing 1")

	currLocation := t5.GetLocation()
	orientation := t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetOrientation()
	//fmt.Println("Current Location: ", currLocation)

	// nearestLoot := t5.ProposeDirection()
	//fmt.Println("Nearest Loot ID: ", nearestLoot)

	currentLootBoxes := t5.GetGameState().GetLootBoxes()
	//fmt.Println("Number of Loot Boxes: ", len(currentLootBoxes))

	if len(currentLootBoxes) > 0 {
		targetPos := currentLootBoxes[targetLootBoxID].GetPosition()
		//fmt.Println("Target Position: ", targetPos)

		deltaXB := targetPos.X - currLocation.X
		deltaYB := targetPos.Y - currLocation.Y
		//fmt.Println("Delta X: ", deltaX, "Delta Y: ", deltaY)

		angleToGoal := math.Atan2(deltaYB, deltaXB) / math.Pi

		audiPos := t5.GetGameState().GetAudi().GetPosition()

		deltaXA := audiPos.X - currLocation.X
		deltaYA := audiPos.Y - currLocation.Y

		angleToAudi := math.Atan2(deltaYA, deltaXA) / math.Pi

		distance_to_audi := math.Sqrt((((deltaXA) * (deltaXA)) + (deltaYA * (deltaYA))))

		if distance_to_audi < (2*utils.CollisionThreshold) && math.Abs(angleToAudi-angleToGoal) < 0.5 {
			angleToGoal = angleToAudi - math.Copysign(0.5, angleToAudi-angleToGoal)
		}

		steer := (angleToGoal - orientation)
		if steer < -1 {
			steer = steer + 2
		} else if steer > 1 {
			steer = steer - 2
		}

		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: steer,
		}

		ownEnergyLevel := t5.GetEnergyLevel()
		Biker_pedal := utils.BikerMaxForce
		if len(t5.GetGameState().GetMegaBikes()[t5.GetBike()].GetAgents()) > 3 {
			if t5.state == 0 {
				Biker_pedal = ownEnergyLevel * utils.BikerMaxForce * 0.5
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
