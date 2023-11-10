package objects

import (
	utils "SOMAS2023/internal/common/utils"

	"math/rand"

	baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
)

// agent with defualt strategy for MVP:
type IBaseBiker interface {
	baseAgent.IAgent[IBaseBiker]
	// DecideAction determines what action the agent is going to take this round.
	// Based on this, the server will call either DecideForce or ChangeBike
	DecideAction(gameState utils.IGameState) BikerAction
	DecideForce(gameState utils.IGameState) utils.Forces                   // defines the vector you pass to the bike: [pedal, brake, turning]
	ChangeBike(gameState utils.IGameState) uuid.UUID                       // action never performed in MVP, might call PickPike() in future implementations
	UpdateColour(totColours utils.Colour)                                  // called if a box of the desired colour has been looted
	UpdateAgent(energyGained float64, energyLost float64, pointGained int) // called by server
	// Loot distribution functions - Refer to Slide 8 of Lec06
	GetEnergyLevel() float64                        // returns the energy level of the agent
	SetEnergyLevel(energyLevel float64)             // sets the energy level of the agent
	GetResourceNeed() float64                       // returns the resource need of the agent
	SetResourceNeed(need float64)                   // sets the resource need of the agent
	GetResourceDemand() float64                     // returns the resource demand of the agent
	SetResourceDemand(demand float64)               // sets the resource demand of the agent
	GetResourceProvision() float64                  // returns the resource provision of the agent
	SetResourceProvision(provision float64)         // sets the resource provision of the agent
	GetResourceAllocation() float64                 // returns the resource allocation of the agent
	SetResourceAllocation(allocation float64)       // sets the resource allocation of the agent
	GetResourceAppropriation() float64              // returns the resource appropriation of the agent
	SetResourceAppropriation(appropriation float64) // sets the resource appropriation of the agent
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
// - asssume server assigns initial bikes to ppl

// What we need to know:
// - How is loot location given to us? (ie utils.IGameState)
// - How does the physics engine work

type BaseBiker struct {
	*baseAgent.BaseAgent[IBaseBiker]              // BaseBiker inherits functions from BaseAgent such as GetID(), GetAllMessages() and UpdateAgentInternalState()
	soughtColour                     utils.Colour // the colour of the lootbox that the agent is currently seeking
	onBike                           bool
	energyLevel                      float64 // float between 0 and 1
	points                           int
	alive                            bool
	// Loot distribution attributes - Refer to Slide 8 of Lec06
	resourceNeed          float64 // 0-1, how much energy the agent needs, for MVP, set to 1 - energyLevel
	resourceDemand        float64 // 0-1, how much energy the agent wants, for MVP, set to same as resourceNeed
	resourceProvision     float64 // 0-1, how much energy the agent is willing to give, for MVP, set to the energy expended in getting to the lootbox
	resourceAllocation    float64 // 0-1, how much energy the server gives the agent, for MVP, set to the energy expended in getting to the lootbox
	resourceAppropriation float64 // 0-1, the proportion of the resourceAllocation that the agent actually gets, for MVP, set to 1
}

// Loot distribution functions - Refer to Slide 8 of Lec06
func (bb *BaseBiker) GetEnergyLevel() float64 {
	return bb.energyLevel
}

func (bb *BaseBiker) SetEnergyLevel(energyLevel float64) {
	bb.energyLevel += energyLevel
}

func (bb *BaseBiker) GetResourceNeed() float64 {
	return bb.resourceNeed
}

func (bb *BaseBiker) SetResourceNeed(need float64) {
	bb.resourceNeed = need
}

func (bb *BaseBiker) GetResourceDemand() float64 {
	return bb.resourceDemand
}

func (bb *BaseBiker) SetResourceDemand(demand float64) {
	bb.resourceDemand = demand
}

func (bb *BaseBiker) GetResourceProvision() float64 {
	return bb.resourceProvision
}

// TODO: this value must change as the agent moves towards the lootbox
func (bb *BaseBiker) SetResourceProvision(provision float64) {
	bb.resourceProvision = provision
}

func (bb *BaseBiker) GetResourceAllocation() float64 {
	return bb.resourceAllocation
}

func (bb *BaseBiker) SetResourceAllocation(allocation float64) {
	bb.resourceAllocation = allocation
}

func (bb *BaseBiker) GetResourceAppropriation() float64 {
	return bb.resourceAppropriation
}

func (bb *BaseBiker) SetResourceAppropriation(appropriation float64) {
	bb.resourceAppropriation = appropriation
}

func (bb *BaseBiker) ResetLootAttributes() {
	bb.resourceNeed = 1 - bb.energyLevel
	bb.resourceDemand = bb.resourceNeed
	bb.resourceProvision = 0
	bb.resourceAllocation = 0
	bb.resourceAppropriation = 1
}

//

func (bb *BaseBiker) DecideAction(gameState utils.IGameState) BikerAction {
	return Pedal
}

// once we know what utils.IGameState looks like we can pass what we need (ie maybe just lootboxes and info on our bike)
func (bb *BaseBiker) DecideForce(gameState utils.IGameState) utils.Forces {
	// the way this is determined depends on how the physics engine works and on what exactly the server passes us
	forces := utils.Forces{
		Pedal:   3.5,
		Brake:   1.2,
		Turning: 2.8,
	}
	return forces
}

// decide which bike to go to
func (bb *BaseBiker) ChangeBike(gameState utils.IGameState) uuid.UUID {
	return uuid.New()
}

func (bb *BaseBiker) UpdateColour(totColours utils.Colour) {
	bb.soughtColour = utils.Colour(rand.Intn(int(totColours)))
}

func (bb *BaseBiker) UpdateAgent(energyGained float64, energyLost float64, pointsGained int) {
	bb.energyLevel += (energyGained - energyLost)
	bb.points += pointsGained
	bb.alive = bb.energyLevel > 0
}

func (bb *BaseBiker) GetLifeStatus() bool {
	return bb.alive
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
