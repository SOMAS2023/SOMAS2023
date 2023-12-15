package team6

import (
	utils "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

func (bb *Team6Biker) DecideForce(direction uuid.UUID) {

	// NEAREST BOX STRATEGY (MVP)
	currLocation := bb.GetLocation()
	currentLootBoxes := bb.GetGameState().GetLootBoxes()

	awdiPos := bb.GetGameState().GetAwdi().GetPosition()
	deltaXAwdi := awdiPos.X - currLocation.X
	deltaYAwdi := awdiPos.Y - currLocation.Y
	distAwdi := math.Sqrt(deltaXAwdi*deltaXAwdi + deltaYAwdi*deltaYAwdi)
	targetPos := currentLootBoxes[direction].GetPosition()
	targetDeltaX := targetPos.X - currLocation.X
	targetDeltaY := targetPos.Y - currLocation.Y
	targetNormAngle := math.Atan2(targetDeltaY, targetDeltaX) / math.Pi

	// Check if there are lootboxes available and move towards closest one
	if distAwdi > distAwdiThreshold {

		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi
		var pedalForce float64

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation(),
		}
		if bb.GetEnergyLevel() < 0.2 {
			pedalForce = 0.1
		} else {
			pedalForce = bb.GetEnergyLevel()
		}
		nearestBoxForces := utils.Forces{
			Pedal:   pedalForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(nearestBoxForces)
	} else { // otherwise move away from awdi
		awdiPos := bb.GetGameState().GetAwdi().GetPosition()

		deltaX := awdiPos.X - currLocation.X
		deltaY := awdiPos.Y - currLocation.Y

		// Steer in opposite direction to awdi
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi
		lootAngle := targetNormAngle - normalisedAngle

		var turningAngle float64
		if lootAngle > 0.0 {
			if math.Abs(lootAngle) < 0.5 {
				turningAngle = 0.5 - math.Abs(lootAngle)
			} else if math.Abs(lootAngle) < 1 {
				turningAngle = 0.0
			} else {
				turningAngle = -(math.Abs(lootAngle) - 1.5)
			}

		} else if lootAngle <= 0.0 {
			if math.Abs(lootAngle) < 0.5 {
				turningAngle = -(0.5 - math.Abs(lootAngle))
			} else if math.Abs(lootAngle) < 1 {
				turningAngle = 0.0
			} else {
				turningAngle = math.Abs(lootAngle) - 1.5
			}

		}
		FinalAngle := targetNormAngle + turningAngle

		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: FinalAngle - bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation(),
		}

		escapeAwdiForces := utils.Forces{
			Pedal:   utils.BikerMaxForce * 0.5,
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(escapeAwdiForces)
	}
}
