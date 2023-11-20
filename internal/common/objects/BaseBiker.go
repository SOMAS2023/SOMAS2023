package objects

import (
	utils "SOMAS2023/internal/common/utils"
	"math"

	"math/rand"

	baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
)

// this struct holds the allocation parameters that we want the allocation protocol to take into account
// These can change based on how we want the allocation to happend, for now they are taken from
// the lecture slides, but more/less could be taken into account.
type ResourceAllocationParams struct {
	resourceNeed          float64 // 0-1, how much energy the agent needs, could be set to 1 - energyLevel
	resourceDemand        float64 // 0-1, how much energy the agent wants, might differ from resourceNeed
	resourceProvision     float64 // 0-1, how much energy the agent has given to reach a goal (could be either the sum of pedaling forces since last lootbox, or the latest pedalling force, or something else
	resourceAppropriation float64 // 0-1, the proportion of what the server allocates that the agent actually gets, for MVP, set to 1
}

// agent with defualt strategy for MVP:
type IBaseBiker interface {
	baseAgent.IAgent[IBaseBiker]
	// Based on this, the server will call either DecideForce or ChangeBike
	DecideAction() BikerAction // determines what action the agent is going to take this round. Based on this, the server will call either DecideForce or ChangeBike
	DecideForce()              // defines the vector you pass to the bike: [pedal, brake, turning]
	GetForces() utils.Forces
	ChangeBike() uuid.UUID                 // called when biker wants to change bike, it will choose which bike to try and join based on agent-specific strategies
	SetBike(uuid.UUID)                     // tells the biker which bike it is on
	UpdateColour(totColours utils.Colour)  // called if a box of the desired colour has been looted
	UpdatePoints(pointGained int)          // called by server
	GetEnergyLevel() float64               // returns the energy level of the agent
	UpdateEnergyLevel(energyLevel float64) // increase the energy level of the agent by the allocated lootbox share or decrease by expended energy
	GetResourceAllocationParams() ResourceAllocationParams
	SetAllocationParameters()
	GetColour() utils.Colour              // returns the colour of the lootbox that the agent is currently seeking
	GetLocation() utils.Coordinates       // gets the agent's location
	UpdateGameState(gameState IGameState) // sets the gameState field at the beginning of each round
}

type BikerAction int

const (
	Pedal BikerAction = iota
	ChangeBike
)

type BaseBiker struct {
	*baseAgent.BaseAgent[IBaseBiker]              // BaseBiker inherits functions from BaseAgent such as GetID(), GetAllMessages() and UpdateAgentInternalState()
	soughtColour                     utils.Colour // the colour of the lootbox that the agent is currently seeking
	onBike                           bool
	energyLevel                      float64 // float between 0 and 1
	points                           int
	alive                            bool
	forces                           utils.Forces
	megaBikeId                       uuid.UUID  // if they are not on a bike it will be 0
	gameState                        IGameState // updated by the server at every round
	allocationParams                 ResourceAllocationParams
}

func (bb *BaseBiker) GetEnergyLevel() float64 {
	return bb.energyLevel
}

// the function will be called by the server to:
// - reduce the energy level based on the force spent pedalling (energyLevel will be neg.ve)
// - increase the energy level after a lootbox has been looted (energyLevel will be pos.ve)
func (bb *BaseBiker) UpdateEnergyLevel(energyLevel float64) {
	bb.energyLevel += energyLevel
	bb.alive = bb.energyLevel > 0
}

func (bb *BaseBiker) GetColour() utils.Colour {
	return bb.soughtColour
}

// this function will be called everytime a lootbox has to be distributed
// these will be defined either based on team strategy and/or according to centralised rules
// for example: it might be decided that the provision must be the average pedalling force provided
// since the last lootbox, an agent might decide as part of their strategy to demand less than they
// need when their energy is above a certain treshold etc etc
func (bb *BaseBiker) SetAllocationParameters() {
	allocParams := ResourceAllocationParams{
		resourceNeed:          1 - bb.energyLevel,
		resourceDemand:        1 - bb.energyLevel,
		resourceProvision:     0,
		resourceAppropriation: 1,
	}
	bb.allocationParams = allocParams
}

// the biker itself doesn't technically have a location (as it's on the map only when it's on a bike)
// in fact this function is only called when the biker needs to make a decision about the pedaling forces
func (bb *BaseBiker) GetLocation() utils.Coordinates {
	megaBikes := bb.gameState.GetMegaBikes()
	return megaBikes[bb.megaBikeId].GetPosition()
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
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestDest = loot.GetPosition()
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
func (bb *BaseBiker) DecideForce() {

	// NEAREST BOX STRATEGY (MVP)
	currLocation := bb.GetLocation()
	nearestLoot := bb.NearestLoot()
	deltaX := nearestLoot.X - currLocation.X
	deltaY := nearestLoot.Y - currLocation.Y
	angle := math.Atan2(deltaX, deltaY)
	angleInDegrees := angle * math.Pi / 180

	nearestBoxForces := utils.Forces{
		Pedal:   utils.BikerMaxForce,
		Brake:   0.0,
		Turning: angleInDegrees,
	}
	bb.forces = nearestBoxForces
}

// decide which bike to go to
// for now it just returns a random uuid
func (bb *BaseBiker) ChangeBike() uuid.UUID {
	return uuid.New()
}

func (bb *BaseBiker) SetBike(bikeId uuid.UUID) {
	bb.megaBikeId = bikeId
}

// this is called when a lootbox of the desidered colour has been looted in order to update the sought colour
func (bb *BaseBiker) UpdateColour(totColours utils.Colour) {
	bb.soughtColour = utils.Colour(rand.Intn(int(totColours)))
}

// update the points at the end of a round
func (bb *BaseBiker) UpdatePoints(pointsGained int) {
	bb.points += pointsGained
}

func (bb *BaseBiker) GetLifeStatus() bool {
	return bb.alive
}

func (bb *BaseBiker) GetForces() utils.Forces {
	return bb.forces
}

func (bb *BaseBiker) UpdateGameState(gameState IGameState) {
	bb.gameState = gameState
}

func (bb *BaseBiker) GetResourceAllocationParams() ResourceAllocationParams {
	return bb.allocationParams
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

/* Added setters below */
func (bb *BaseBiker) SetSoughtColour(colour utils.Colour) {
	bb.soughtColour = colour
}

func (bb *BaseBiker) SetOnBike(onBike bool) {
	bb.onBike = onBike
}

func (bb *BaseBiker) SetEnergyLevel(energyLevel float64) {
	bb.energyLevel = energyLevel
}

func (bb *BaseBiker) SetPoints(points int) {
	bb.points = points
}

func (bb *BaseBiker) SetAlive(alive bool) {
	bb.alive = alive
}

func (bb *BaseBiker) SetForces(forces utils.Forces) {
	bb.forces = forces
}

func (bb *BaseBiker) SetMegaBikeId(megaBikeId uuid.UUID) {
	bb.megaBikeId = megaBikeId
}

func (bb *BaseBiker) SetGameState(gameState IGameState) {
	bb.gameState = gameState
}

func (bb *BaseBiker) SetAllocationParams(allocationParams ResourceAllocationParams) {
	bb.allocationParams = allocationParams
}

func (bb *BaseBiker) GetGameState() IGameState {
	return bb.gameState
}

func (bb *BaseBiker) GetMegaBikeId() uuid.UUID {
	return bb.megaBikeId
}
