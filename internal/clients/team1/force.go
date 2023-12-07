// GETS FORCES FOR BIKER

package team1

import (
	utils "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

// // -----------------PEDALLING FORCE FUNCTIONS------------------
func (bb *Biker1) getPedalForce() float64 {
	//can be made more complex
	return utils.BikerMaxForce * bb.GetEnergyLevel()
}

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootboX
// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Biker1) DecideForce(direction uuid.UUID) {

	bb.recentDecided = direction
	bb.recentDecidedColour = bb.GetGameState().GetLootBoxes()[direction].GetColour()
	bb.recentDecidedPosition = bb.GetGameState().GetLootBoxes()[direction].GetPosition()

	bb.prevOnBike = bb.GetBikeStatus()
	lootBoxes := bb.GetGameState().GetLootBoxes()
	currLocation := bb.GetLocation()
	targetPos := lootBoxes[direction].GetPosition()

	// If audi is close, steer away from it
	if bb.DistanceFromAudi(bb.GetBikeInstance()) < audiDistanceThreshold {
		audiPos := bb.GetGameState().GetAudi().GetPosition()
		deltaX := audiPos.X - currLocation.X
		deltaY := audiPos.Y - currLocation.Y
		// Steer in opposite direction to audi (regardless of governance)
		angle := math.Atan2(-deltaY, -deltaX)
		normalisedAngle := angle / math.Pi
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - bb.GetBikeInstance().GetOrientation(),
		}

		escapeAudiForces := utils.Forces{
			Pedal:   bb.getPedalForce(),
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(escapeAudiForces)
	}

	//agent doesn't rebel, just decides to leave next round if dislike vote

	if bb.recentVote != nil {
		result, ok := bb.recentVote[direction]
		if ok && result < votingAlignmentThreshold {
			bb.dislikeVote = true
		} else {
			bb.dislikeVote = false
		}
	}

	deltaX := targetPos.X - currLocation.X
	deltaY := targetPos.Y - currLocation.Y
	angle := math.Atan2(deltaY, deltaX)
	normalisedAngle := angle / math.Pi

	// if the governance is ruler-based and we're not the ruler, don't steer
	var turningDecision utils.TurningDecision
	bike := bb.GetGameState().GetMegaBikes()[bb.GetBike()]
	gov := bike.GetGovernance()
	if gov == utils.Dictatorship || gov == utils.Leadership {
		ruler := bike.GetRuler()
		if ruler != bb.GetID() {
			turningDecision = utils.TurningDecision{
				SteerBike:     false,
				SteeringForce: 0.0,
			}
		} else {
			turningDecision = utils.TurningDecision{
				SteerBike:     true,
				SteeringForce: normalisedAngle - bb.GetBikeInstance().GetOrientation(),
			}
		}
	} else {
		turningDecision = utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - bb.GetBikeInstance().GetOrientation(),
		}
	}

	boxForces := utils.Forces{
		Pedal:   bb.getPedalForce(),
		Brake:   0.0,
		Turning: turningDecision,
	}
	bb.SetForces(boxForces)
}

// -----------------END OF PEDALLING FORCE FUNCTIONS------------------
