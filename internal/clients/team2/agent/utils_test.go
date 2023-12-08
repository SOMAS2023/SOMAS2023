package agent

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRuleAdherenceValue_SameDirection(t *testing.T) {
	// Prepare test data
	agentID := uuid.New()
	turningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: 0.25,
	}
	expectedAction := utils.Forces{Pedal: 1, Turning: turningDecision}
	actualAction := utils.Forces{Pedal: 0.2, Turning: turningDecision}

	// Create an instance of AgentTwo
	agent := NewBaseTeam2Biker(objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()))

	// Call the function
	result := agent.Modules.Utils.RuleAdherenceValue(agentID, expectedAction, actualAction)

	// Since there is 0 directional difference, it should return 1.0*0.2 = 0.2
	expectedResult := 0.2

	// Assert the outcome
	assert.Equal(t, expectedResult, result, "The result should match the expected rule adherence value")
}

func TestRuleAdherenceValue_OppositeDirection(t *testing.T) {
	// Prepare test data
	agentID := uuid.New()
	ExpectedTurningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: 0.25,
	}
	ActualTurningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: -0.75,
	}
	expectedAction := utils.Forces{Pedal: 1, Turning: ExpectedTurningDecision}
	actualAction := utils.Forces{Pedal: 0.2, Turning: ActualTurningDecision}

	// Create an instance of AgentTwo
	agent := NewBaseTeam2Biker(objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()))

	// Call the function
	result := agent.Modules.Utils.RuleAdherenceValue(agentID, expectedAction, actualAction)

	// Since there is 180 degree directional difference, it should return -1.0*0.2 (weighting) = -0.2
	expectedResult := -0.2

	// Assert the outcome
	assert.Equal(t, expectedResult, result, "The result should match the expected rule adherence value")
}

func TestRuleAdherenceValue_OrthogonalDirection(t *testing.T) {
	// Prepare test data
	agentID := uuid.New()
	ExpectedTurningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: 0.25,
	}
	ActualTurningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: -0.25,
	}
	expectedAction := utils.Forces{Pedal: 1, Turning: ExpectedTurningDecision}
	actualAction := utils.Forces{Pedal: 0.2, Turning: ActualTurningDecision}

	// Create an instance of AgentTwo
	agent := NewBaseTeam2Biker(objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()))

	// Call the function
	result := agent.Modules.Utils.RuleAdherenceValue(agentID, expectedAction, actualAction)

	// Since there is 90 degree directional difference, it should return 0*0.2 (weighting) = 0
	expectedResult := 0
	threshold := 0.001

	// Assert the outcome
	assert.LessOrEqual(t, result-float64(expectedResult), threshold, "Orthogonal vectors should have a cosine similarity of 0")
}
