package team5Agent

import (
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

func (t5 *team5Agent) calcForce() float64 {
	agentEnergy := t5.GetEnergyLevel()
	if t5.state == 0 {
		return agentEnergy * utils.BikerMaxForce * 0.5
	} else if t5.state == 1 {
		return agentEnergy * utils.BikerMaxForce
	} else if t5.state == 2 {
		return agentEnergy * utils.BikerMaxForce
	}
	return utils.BikerMaxForce
}

func (t5 *team5Agent) DecideForce(targetLootBoxID uuid.UUID) {

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

		awdiPos := t5.GetGameState().GetAwdi().GetPosition()

		deltaXA := awdiPos.X - currLocation.X
		deltaYA := awdiPos.Y - currLocation.Y

		angleToAwdi := math.Atan2(deltaYA, deltaXA) / math.Pi

		distance_to_awdi := math.Sqrt((((deltaXA) * (deltaXA)) + (deltaYA * (deltaYA))))

		if distance_to_awdi < (2*utils.CollisionThreshold) && math.Abs(angleToAwdi-angleToGoal) < 0.5 {
			angleToGoal = angleToAwdi - math.Copysign(0.5, angleToAwdi-angleToGoal)
		}

		steer := (angleToGoal - orientation)
		if steer < -1.0 {
			steer = steer + 2.0
		} else if steer > 1.0 {
			steer = steer - 2.0
		}

		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: steer,
		}

		//ownEnergyLevel := t5.GetEnergyLevel()
		biker_pedal := t5.calcForce()

		forces_movement := utils.Forces{
			Pedal:   biker_pedal,
			Brake:   0.0,
			Turning: turningDecision,
		}
		t5.SetForces(forces_movement)
	}
}
