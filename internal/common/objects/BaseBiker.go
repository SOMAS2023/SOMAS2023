package objects

import (
	utils "SOMAS2023/internal/common/utils"
	"math"

	"math/rand"

	baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
)

// agent with defualt strategy for MVP:
type IBaseBiker interface {
	baseAgent.IAgent[IBaseBiker]
	// Based on this, the server will call either DecideForce or ChangeBike
	DecideAction() BikerAction                                             // determines what action the agent is going to take this round. Based on this, the server will call either DecideForce or ChangeBike
	DecideForce() utils.Forces                                             // defines the vector you pass to the bike: [pedal, brake, turning]
	ChangeBike() uuid.UUID                                                 // called when biker wants to change bike, it will choose which bike to try and join based on agent-specific strategies
	UpdateColour(totColours utils.Colour)                                  // called if a box of the desired colour has been looted
	UpdateAgent(energyGained float64, energyLost float64, pointGained int) // called by server
	GetLocation() utils.Coordinates                                        // gets the agent's location
	UpdateGameState(gameState IGameState)                                  // sets the gameState field at the beginning of each round
}

type BikerAction int

const (
	Pedal BikerAction = iota
	ChangeBike
)

// Assumptions:
// - the server will update the energy level of the agent at the end of the round (both by subtracting the spent energy
// by adding the loot energy if relevant)
// - server calls for agent to update colour
// - server gives loot based on alloc decision (MVP)
// - server gives points based on colour of loot
// - when biker dies the server deletes the instance (can have it set a status field "alive" to false instead if we
// want to keep a record of dead bikers)
// - server keeps track of round number
// - assume server assigns initial bikes to ppl

type BaseBiker struct {
	*baseAgent.BaseAgent[IBaseBiker]              // BaseBiker inherits functions from BaseAgent such as GetID(), GetAllMessages() and UpdateAgentInternalState()
	soughtColour                     utils.Colour // the colour of the lootbox that the agent is currently seeking
	onBike                           bool         // whether an agent is on a bike or not. The game allows for agents to not be on a bike during a round, in that case they will technically be "outside" of the map in a sort of idle state
	energyLevel                      float64      // float between 0 and 1
	points                           int          // tracks the number of lootboxes of the right colour have been looted
	alive                            bool         // an agent can die if: 1. they run out of energy (energyLevel = 0) 2. they get hit by the Audi
	megaBikeId                       uuid.UUID    // if they are not on a bike it will be 0
	gameState                        IGameState   // updated by the server at every round
}

// the biker itself doesn't technically have a location (as it's on the map only when it's on a bike)
// in fact this function is only called when the biker needs to make a decision about the pedaling forces
func (bb *BaseBiker) GetLocation() utils.Coordinates {
	megaBikes := bb.gameState.GetMegaBikes()
	return megaBikes[bb.megaBikeId].coordinates
}

// returns the nearest lootbox with respect to the agent's bike current position
// in the MVP this is used to determine the pedalling forces as all agent will be
// aiming to get to the closest lootbox by default
func (bb *BaseBiker) NearestLoot() utils.Coordinates {
	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestDest utils.Coordinates
	var currDist float64
	for _, loot := range bb.gameState.GetLootBoxes() {
		x, y := loot.coordinates.X, loot.coordinates.Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestDest = loot.coordinates
			shortestDist = currDist
		}
	}
	return nearestDest
}

// in the MVP the biker's action defaults to pedaling (as it won't be able to change bikes)
// in future implementations this function will be overridden by the agent's specific strategy
// which will be used to determine whether to pedalor try to change bike
func (bb *BaseBiker) DecideAction() BikerAction {
	return Pedal
}

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox
func (bb *BaseBiker) DecideForce() utils.Forces {

	// NEAREST BOX STRATEGY (MVP)
	currLocation := bb.GetLocation()
	nearestLoot := bb.NearestLoot()
	deltaX := nearestLoot.X - currLocation.X
	deltaY := nearestLoot.Y - currLocation.Y
	angle := math.Atan2(deltaX, deltaY)
	angleInDegrees := angle * math.Pi / 180

	nearestBoxForces := utils.Forces{
		Pedal:   1.0,
		Brake:   0.0,
		Turning: angleInDegrees,
	}

	return nearestBoxForces
}

// decide which bike to go to
// for now it just returns a random uuid
func (bb *BaseBiker) ChangeBike() uuid.UUID {
	return uuid.New()
}

// this is called when a lootbox of the desidered colour has been looted in order to update the sought colour
func (bb *BaseBiker) UpdateColour(totColours utils.Colour) {
	bb.soughtColour = utils.Colour(rand.Intn(int(totColours)))
}

// update the energy levels and points at the end of a round
func (bb *BaseBiker) UpdateAgent(energyGained float64, energyLost float64, pointsGained int) {
	bb.energyLevel += (energyGained - energyLost)
	bb.points += pointsGained
	bb.alive = bb.energyLevel > 0
}

func (bb *BaseBiker) GetLifeStatus() bool {
	return bb.alive
}

func (bb *BaseBiker) UpdateGameState(gameState IGameState) {
	bb.gameState = gameState
}

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIBaseBiker(totColours utils.Colour, bikeId uuid.UUID) IBaseBiker {
	return &BaseBiker{
		BaseAgent:    baseAgent.NewBaseAgent[IBaseBiker](),
		soughtColour: utils.GenerateRandomColour(),
		onBike:       true,
		energyLevel:  1.0,
		points:       0,
		alive:        true,
	}
}

// this function will be used by GetTeamAgent to get the ref to the BaseBiker
func GetBaseBiker(totColours utils.Colour, bikeId uuid.UUID) *BaseBiker {
	return &BaseBiker{
		BaseAgent:    baseAgent.NewBaseAgent[IBaseBiker](),
		soughtColour: utils.GenerateRandomColour(),
		onBike:       true,
		energyLevel:  1.0,
		points:       0,
		alive:        true,
	}
}
