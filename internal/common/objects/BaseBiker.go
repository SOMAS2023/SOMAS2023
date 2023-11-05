package basebiker

import (
	utils "SOMAS2023/internal/common/utils"

	"math/rand"

	baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
)

// placeholder
type ServerStatus struct {
	// []lootboxes
	// []bikes
	// []bikers
}

// agent with defualt strategy for MVP:
type IBaseBiker interface {
	baseAgent.IAgent[IBaseBiker]
	DecideAction(server_status ServerStatus) int
	DecideForce(server_status ServerStatus) utils.Forces                      // defines the vector you pass to the bike: [pedal, brake, turning]
	ChangeBike(server_status ServerStatus) uuid.UUID                          // action never performed in MVP, might call PickPike() in future implementations
	UpdateColour(tot_colours utils.Colour)                                    // called if a box of the desired colour has been looted
	UpdateAgent(energy_gained float64, energy_lost float64, point_gained int) // called by server
}

// Assumptions:
// - the server will update the energy level of the agent at the end of the round (both by subtracting the spent energy
// by adding the loot energy if relevant)
// - server calls for agent to update colour
// - server gives loot based on alloc decision (MVP)
// - server gives points based on colour of loot
// - when biker dies the server deletes the instance (can have it set a status field "alive" to false instead if we
// want to keep a record of dead bikers)
// - server keeps track of round number
// - asssume server assigns initial bikes to ppl

// What we need to know:
// - How is loot location given to us? (ie ServerStatus)
// - How does the physics engine work

type BaseBiker struct {
	*baseAgent.BaseAgent[IBaseBiker]
	sought_colour utils.Colour // the colour of the lootbox that the agent is currently seeking
	on_bike       bool
	energy_level  float64 // float between 0 and 1
	points        int
	bike_id       uuid.UUID
	alive         bool
}

// returns 0 if biker decides to pedal and 1 if it decides to change bike
// based on this the server will call either DecideForce or ChangeBike
func (bb *BaseBiker) DecideAction(server_status ServerStatus) int {
	return 0
}

// once we know what ServerStatus looks like we can pass what we need (ie maybe just lootboxes and info on our bike)
func (bb *BaseBiker) DecideForce(server_status ServerStatus) utils.Forces {
	// the way this is determined depends on how the physics engine works and on what exactly the server passes us
	forces := utils.Forces{
		Pedal:   3.5,
		Brake:   1.2,
		Turning: 2.8,
	}
	return forces
}

// decide which bike to go to
func (bb *BaseBiker) ChangeBike(server_status ServerStatus) uuid.UUID {
	return uuid.New()
}

func (bb *BaseBiker) UpdateColour(tot_colours utils.Colour) {
	bb.sought_colour = utils.Colour(rand.Intn(int(tot_colours)))
}

func (bb *BaseBiker) UpdateAgent(energy_gained float64, energy_lost float64, points_gained int) {
	bb.energy_level += (energy_gained - energy_lost)
	bb.points += points_gained
	bb.alive = bb.energy_level > 0
}

func (bb *BaseBiker) GetLifeStatus() bool {
	return bb.alive
}

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIBaseBiker(tot_colours utils.Colour, bike_id uuid.UUID) IBaseBiker {
	return &BaseBiker{
		BaseAgent:     baseAgent.NewBaseAgent[IBaseBiker](),
		sought_colour: utils.Colour(rand.Intn(int(tot_colours))),
		on_bike:       true,
		energy_level:  1.0,
		points:        0,
		bike_id:       bike_id,
		alive:         true,
	}
}

// this function will be used by GetTeamAgent to get the ref to the BaseBiker
func GetBaseBiker(tot_colours utils.Colour, bike_id uuid.UUID) *BaseBiker {
	return &BaseBiker{
		BaseAgent:     baseAgent.NewBaseAgent[IBaseBiker](),
		sought_colour: utils.Colour(rand.Intn(int(tot_colours))),
		on_bike:       true,
		energy_level:  1.0,
		points:        0,
		bike_id:       bike_id,
		alive:         true,
	}
}
