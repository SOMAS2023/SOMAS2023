package objects

import (
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"math"

	"math/rand"

	baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
)

// this struct holds the allocation parameters that we want the allocation protocol to take into account
// These can change based on how we want the allocation to happend, for now they are taken from
// the lecture slides, but more/less could be taken into account.
type ResourceAllocationParams struct {
	ResourceNeed          float64 `json:"need"`          // 0-1, how much energy the agent needs, could be set to 1 - energyLevel
	ResourceDemand        float64 `json:"demand"`        // 0-1, how much energy the agent wants, might differ from ResourceNeed
	ResourceProvision     float64 `json:"provision"`     // 0-1, how much energy the agent has given to reach a goal (could be either the sum of pedaling forces since last lootbox, or the latest pedalling force, or something else
	ResourceAppropriation float64 `json:"appropriation"` // 0-1, the proportion of what the server allocates that the agent actually gets, for MVP, set to 1
}

type IBaseBiker interface {
	baseAgent.IAgent[IBaseBiker]

	DecideAction() BikerAction                                      // ** determines what action the agent is going to take this round. (changeBike or Pedal)
	DecideForce(direction uuid.UUID)                                // ** defines the vector you pass to the bike: [pedal, brake, turning]
	DecideJoining(pendinAgents []uuid.UUID) map[uuid.UUID]bool      // ** decide whether to accept or not accept bikers, ranks the ones
	ChangeBike() uuid.UUID                                          // ** called when biker wants to change bike, it will choose which bike to try and join
	ProposeDirection() uuid.UUID                                    // ** returns the id of the desired lootbox based on internal strategy
	FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap // ** stage 3 of direction voting
	DecideAllocation() voting.IdVoteMap                             // ** decide the allocation parameters

	GetForces() utils.Forces                               // returns forces for current round
	GetColour() utils.Colour                               // returns the colour of the lootbox that the agent is currently seeking
	GetLocation() utils.Coordinates                        // gets the agent's location
	GetBike() uuid.UUID                                    // tells the biker which bike it is on
	GetEnergyLevel() float64                               // returns the energy level of the agent
	GetResourceAllocationParams() ResourceAllocationParams // returns set allocation parameters
	GetBikeStatus() bool                                   // returns whether the biker is on a bike or not

	SetBike(uuid.UUID)                     // sets the megaBikeID. this is either the id of the bike that the agent is on or the one that it's trying to join
	SetForces(forces utils.Forces)         // sets the forces (to be updated in DecideForces())
	UpdateColour(totColours utils.Colour)  // called if a box of the desired colour has been looted
	UpdatePoints(pointGained int)          // called by server
	UpdateEnergyLevel(energyLevel float64) // increase the energy level of the agent by the allocated lootbox share or decrease by expended energy
	UpdateGameState(gameState IGameState)  // sets the gameState field at the beginning of each round
	ToggleOnBike()                         // called when removing or adding a biker on a bike
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
}

func (bb *BaseBiker) GetColour() utils.Colour {
	return bb.soughtColour
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *BaseBiker) DecideAllocation() voting.IdVoteMap {
	bikeID := bb.GetBike()
	fellowBikers := bb.gameState.GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	for _, agent := range fellowBikers {
		if agent.GetID() == bb.GetID() {
			distribution[agent.GetID()] = 1.0
		} else {
			distribution[agent.GetID()] = 0.0
		}
	}
	return distribution
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
func (bb *BaseBiker) nearestLoot() uuid.UUID {
	currLocation := bb.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestBox uuid.UUID
	var currDist float64
	for _, loot := range bb.gameState.GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return nearestBox
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

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *BaseBiker) DecideForce(direction uuid.UUID) {

	// NEAREST BOX STRATEGY (MVP)
	currLocation := bb.GetLocation()
	nearestLoot := bb.nearestLoot()
	currentLootBoxes := bb.gameState.GetLootBoxes()

	// Check if there are lootboxes available and move towards closest one
	if len(currentLootBoxes) > 0 {
		targetPos := currentLootBoxes[nearestLoot].GetPosition()

		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaX, deltaY)
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - bb.gameState.GetMegaBikes()[bb.megaBikeId].GetOrientation(),
		}

		nearestBoxForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(nearestBoxForces)
	} else { // otherwise move away from audi
		audiPos := bb.GetGameState().GetAudi().GetPosition()

		deltaX := audiPos.X - currLocation.X
		deltaY := audiPos.Y - currLocation.Y

		// Steer in opposite direction to audi
		angle := math.Atan2(-deltaX, -deltaY)
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle,
		}

		escapeAudiForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		bb.SetForces(escapeAudiForces)
	}
}

// decide which bike to go to
// for now it just returns a random uuid
func (bb *BaseBiker) ChangeBike() uuid.UUID {
	return uuid.New()
}

func (bb *BaseBiker) SetBike(bikeId uuid.UUID) {
	bb.megaBikeId = bikeId
}

func (bb *BaseBiker) GetBike() uuid.UUID {
	return bb.megaBikeId
}

// this is called when a lootbox of the desidered colour has been looted in order to update the sought colour
func (bb *BaseBiker) UpdateColour(totColours utils.Colour) {
	bb.soughtColour = utils.Colour(rand.Intn(int(totColours)))
}

// update the points at the end of a round
func (bb *BaseBiker) UpdatePoints(pointsGained int) {
	bb.points += pointsGained
}

func (bb *BaseBiker) GetForces() utils.Forces {
	return bb.forces
}

func (bb *BaseBiker) SetForces(forces utils.Forces) {
	bb.forces = forces
}

func (bb *BaseBiker) UpdateGameState(gameState IGameState) {
	bb.gameState = gameState
}

func (bb *BaseBiker) GetResourceAllocationParams() ResourceAllocationParams {
	return bb.allocationParams
}

// default implementation returns the id of the nearest lootbox
func (bb *BaseBiker) ProposeDirection() uuid.UUID {
	return bb.nearestLoot()
}

func (bb *BaseBiker) ToggleOnBike() {
	bb.onBike = !bb.onBike
}

func (bb *BaseBiker) GetBikeStatus() bool {
	return bb.onBike
}

func (bb *BaseBiker) GetGameState() IGameState {
	return bb.gameState
}

func (bb *BaseBiker) GetMegaBikeId() uuid.UUID {
	return bb.megaBikeId
}

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *BaseBiker) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *BaseBiker) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	votes := make(voting.LootboxVoteMap)
	totOptions := len(proposals)
	normalDist := 1.0 / float64(totOptions)
	for _, proposal := range proposals {
		votes[proposal] = normalDist
	}
	return votes
}

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIBaseBiker(totColours utils.Colour, bikeId uuid.UUID) IBaseBiker {
	return &BaseBiker{
		BaseAgent:    baseAgent.NewBaseAgent[IBaseBiker](),
		soughtColour: utils.GenerateRandomColour(),
		onBike:       true,
		energyLevel:  1.0,
		points:       0,
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
	}
}
