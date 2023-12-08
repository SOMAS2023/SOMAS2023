package agent

import (
	"SOMAS2023/internal/clients/team2/modules"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseTeam2Biker(t *testing.T) {
	agent := NewBaseTeam2Biker(objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()))
	assert.NotNil(t, agent)
	assert.Equal(t, 0, agent.BaseBiker.GetPoints())
	assert.Equal(t, 1.0, agent.BaseBiker.GetEnergyLevel())
}

func TestClippingSocialCapital(t *testing.T) {
	agent := NewBaseTeam2Biker(objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()))
	testAgentID := uuid.New()

	// Set up predefined values for trust, institution, and network
	agent.Modules.SocialCapital.Reputation[testAgentID] = agent.Modules.SocialCapital.ClipValues(1.3)
	agent.Modules.SocialCapital.Institution[testAgentID] = agent.Modules.SocialCapital.ClipValues(-0.3)

	assert.Equal(t, 1.0, agent.Modules.SocialCapital.Reputation[testAgentID])
	assert.Equal(t, 0.0, agent.Modules.SocialCapital.Institution[testAgentID])
}

func TestForcesToVectorConversion(t *testing.T) {
	force := utils.Forces{
		Pedal: 2.0,
		Turning: utils.TurningDecision{
			SteeringForce: 0.25, // 45 degrees since -1, 1 maps to -180, 180
		},
	}

	expectedVector := modules.ForceVector{X: 1.414, Y: 1.414}

	resultVector := modules.GetForceVector(force)
	// since floating point, need comparison within threshold
	assert.InDelta(t, expectedVector.X, resultVector.X, 0.001)
	assert.InDelta(t, expectedVector.Y, resultVector.Y, 0.001)
}

func TestCosineSimilarity(t *testing.T) {
	v1 := modules.ForceVector{X: 1, Y: 0}
	v2 := modules.ForceVector{X: 0, Y: 1}

	expectedResult := 0.0 // dot product is 0 since vectors are perpendicular
	result := v1.CosineSimilarity(v2)

	assert.Equal(t, expectedResult, result)
}
