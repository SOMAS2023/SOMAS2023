package team6

import (
	utils "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

func (bb *Team6Biker) DecideForce(direction uuid.UUID) {

	// NEAREST BOX STRATEGY (MVP)
	initialOrientaion := bb.GetGameState().GetMegaBikes()[bb.GetBike()].GetOrientation()
	currLocation := bb.GetLocation()

	audiPos := bb.GetGameState().GetAudi().GetPosition()
	deltaXAudi := audiPos.X - currLocation.X
	deltaYAudi := audiPos.Y - currLocation.Y
	distAudi := math.Sqrt(deltaXAudi*deltaXAudi + deltaYAudi*deltaYAudi)

	// Check if there are lootboxes available and move towards closest one
	if distAudi > distAudiThreshold {
		targetPos := bb.GetGameState().GetLootBoxes()[direction].GetPosition()
		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaX, deltaY)
		normalisedAngle := angle / math.Pi
		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - initialOrientaion,
		}

		nearestBoxForces := utils.Forces{
			Pedal:   bb.GetEnergyLevel(),
			Brake:   0.0,
			Turning: turningDecision,
		}

		bb.SetForces(nearestBoxForces)
	} else { // otherwise move away from audi
		// Steer in opposite direction to audi
		angle := math.Atan2(-deltaXAudi, -deltaYAudi) / math.Pi
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - initialOrientaion,
		}

		escapeAudiForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(escapeAudiForces)
	}
}
