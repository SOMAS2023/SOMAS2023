package modules

import (
	"SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

// UtilsModule - Module for handling various Utils
type UtilsModule struct{}

// NewUtilsModule - Constructor for UtilsModule
func NewUtilsModule() *UtilsModule {
	return &UtilsModule{}
}

func (um *UtilsModule) ProjectForce(actual, expected utils.Forces) float64 {
	actualVec := GetForceVector(actual)
	expectVec := GetForceVector(expected)
	return actualVec.CosineSimilarity(*expectVec) * actual.Pedal
}

// Get the forces to the target coordinated
func (um *UtilsModule) GetForcesToTarget(agentPosition, targetPosition utils.Coordinates) utils.Forces {

	deltaX := targetPosition.X - agentPosition.X
	deltaY := targetPosition.Y - agentPosition.Y
	angle := math.Atan2(deltaY, deltaX)
	normalisedAngle := angle / math.Pi
	turningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: normalisedAngle,
	}
	return utils.Forces{
		Pedal:   utils.BikerMaxForce,
		Brake:   0.0,
		Turning: turningDecision,
	}
}

// Called by Events to obtain Event Value for update Institution
// Assume what they broadcast is the truth
func (um *UtilsModule) RuleAdherenceValue(agentID uuid.UUID, expectedAction, actualAction utils.Forces) float64 {
	actualVec := GetForceVector(actualAction)
	expectVec := GetForceVector(expectedAction)
	return actualVec.CosineSimilarity(*expectVec) * actualAction.Pedal
}

// GetForcesToTargetWithDirectionOffset calculates the forces to be applied on an agent to steer towards a target position,
// taking into account a specified degree of angular offset.
func (um *UtilsModule) GetForcesToTargetWithDirectionOffset(force, degree float64, currPos, targetPos utils.Coordinates) utils.Forces {
	deltaX := targetPos.X - currPos.X
	deltaY := targetPos.Y - currPos.Y
	angle := math.Atan2(deltaY, deltaX)
	normalisedAngle := angle/math.Pi + math.Remainder(degree, 2)

	if normalisedAngle < -1 {
		normalisedAngle = normalisedAngle + 2
	} else if normalisedAngle > 1 {
		normalisedAngle = normalisedAngle - 2
	}
	turningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: normalisedAngle,
	}
	return utils.Forces{
		Pedal:   force,
		Brake:   0.0,
		Turning: turningDecision,
	}
}
