package team2

import (
	"SOMAS2023/internal/common/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseTeam2Biker(t *testing.T) {
	agentId := uuid.New()
	agent := NewBaseTeam2Biker(agentId)
	assert.NotNil(t, agent)
	assert.Equal(t, 0, agent.points)
	assert.Equal(t, 1.0, agent.energyLevel)
}

func TestCalculateSocialCapital(t *testing.T) {
	agent := NewBaseTeam2Biker(uuid.New())
	testAgentID := uuid.New()

	// Set up predefined values for trust, institution, and network
	agent.Trust[testAgentID] = 0.8
	agent.Institution[testAgentID] = 0.3
	agent.Network[testAgentID] = 0.5

	agent.CalculateSocialCapital()

	expectedSocialCapital := TrustWeight*0.8 + InstitutionWeight*0.3 + NetworkWeight*0.5
	assert.Equal(t, expectedSocialCapital, agent.SocialCapital[testAgentID])
}

func TestForcesToVectorConversion(t *testing.T) {
	force := utils.Forces{
		Pedal: 2.0,
		Turning: utils.TurningDecision{
			SteeringForce: 0.25, // 45 degrees since -1, 1 maps to -180, 180
		},
	}

	expectedVector := ForceVector{X: 1.414, Y: 1.414}

	resultVector := forcesToVectorConversion(force)
	// since floating point, need comparison within threshold
	assert.InDelta(t, expectedVector.X, resultVector.X, 0.001)
	assert.InDelta(t, expectedVector.Y, resultVector.Y, 0.001)
}

func TestCosineSimilarity(t *testing.T) {
	v1 := ForceVector{X: 1, Y: 0}
	v2 := ForceVector{X: 0, Y: 1}

	expectedResult := 0.0 // dot product is 0 since vectors are perpendicular
	result := cosineSimilarity(v1, v2)

	assert.Equal(t, expectedResult, result)
}
