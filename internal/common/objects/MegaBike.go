package objects

import (
	//"fmt"
	utils "SOMAS2023/internal/common/utils"
	//baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
)

// MegaBike will have the following forces
type MegaBike struct {
	id           uuid.UUID
	coordinates  utils.Coordinates
	agentsOnBike []uuid.UUID // TODO: add agents to this list when they change/leave bikes
	agentForces  []utils.Forces
	Mass         float64
	acceleration float64
	velocity     float64
	orientation  float64
}

// GetMegaBike is a constructor for MegaBike that initializes it with a new UUID and default position.
func GetMegaBike() *MegaBike {
	return &MegaBike{
		id:          uuid.New(),                        // Generate a new unique identifier
		coordinates: utils.GenerateRandomCoordinates(), // Initialize to randomized position
	}
}

// Getter for agentsOnBike
func (mb *MegaBike) GetAgentsOnBike() []uuid.UUID {
	return mb.agentsOnBike
}

// returns the unique ID of the object
func (mb *MegaBike) GetID() uuid.UUID {
	return mb.id
}

// returns the current coordinates of the object
func (mb *MegaBike) GetPosition() utils.Coordinates {
	return mb.coordinates
}

func (mb *MegaBike) Calculate_Force() float64 {
	force_map := 4.0
	if len(mb.agentForces) == 0 {
		return 0.0
	}
	Total_pedal := 0.0
	Total_brake := 0.0
	Total_mass := utils.MassBike
	for _, agent := range mb.agentForces {
		Total_mass += utils.MassBiker
		if agent.Pedal != 0 {
			Total_pedal += float64(agent.Pedal)
		} else {
			Total_brake += float64(agent.Brake)
		}
	}
	F := force_map * (float64(Total_pedal) - float64(Total_brake))
	return F
}

func (mb *MegaBike) Add_Agent(agentForces utils.Forces) {
	mb.agentForces = append(mb.agentForces, agentForces)
}

func (mb *MegaBike) Calculate_Orientation() {
	if len(mb.agentForces) == 0 {
		return
	}
	Total_turning := 0.0
	for _, agent := range mb.agentForces {
		Total_turning += float64(agent.Turning)
	}
	Average_turning := Total_turning / float64(len(mb.agentForces))
	mb.orientation += (Average_turning)
}
