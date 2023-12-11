package team6

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"testing"

	//"fmt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInitialiseBiker6(t *testing.T) {
	agent := InitialiseBiker6(objects.GetBaseBiker(utils.GenerateRandomColour(), uuid.New()))
	assert.NotNil(t, agent)
	assert.Equal(t, 0, agent.GetPoints())
	assert.Equal(t, 1.0, agent.GetEnergyLevel())
}

// func GetMostCommonColor(agents []objects.IBaseBiker) (utils.Colour, int, int)
func TestGetMostCommonColor(t *testing.T) {
	agent1 := objects.GetBaseBiker(utils.Red, uuid.New())
	agent2 := InitialiseBiker6(objects.GetBaseBiker(utils.Red, uuid.New()))
	agent3 := InitialiseBiker6(objects.GetBaseBiker(utils.Red, uuid.New()))
	agent4 := InitialiseBiker6(objects.GetBaseBiker(utils.Blue, uuid.New()))
	agent5 := InitialiseBiker6(objects.GetBaseBiker(utils.Yellow, uuid.New()))
	agent6 := InitialiseBiker6(objects.GetBaseBiker(utils.Blue, uuid.New()))
	agent7 := InitialiseBiker6(objects.GetBaseBiker(utils.Blue, uuid.New()))
	agents := []objects.IBaseBiker{agent1, agent2, agent3, agent4, agent5, agent6, agent7}

	color, _, _ := GetMostCommonColor(agents)

	assert.Equal(t, color, utils.Red)
}
