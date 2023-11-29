// AccessBaseBiker is an example of how to access BaseBiker fields and methods
// func (a *AgentTwo) AccessBaseBiker() {
//     // Accessing a field of BaseBiker
//     a.BaseBiker.SomeField = "some value"

//     // Calling a method of BaseBiker
//     a.BaseBiker.SomeMethod()
// }

// TODO: Reputation evaluation

package team2

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"fmt"

	// "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

type IBaseBiker interface {
	objects.IBaseBiker
}

type AgentTwo struct {
	// BaseBiker represents a basic biker agent.
	*objects.BaseBiker
	// CalculateSocialCapitalOtherAgent: (trustworthiness - cosine distance, social networks - friends, institutions - num of rounds on a bike)
	SocialCapital      map[uuid.UUID]float64 // Social Captial of other agents
	Trust              map[uuid.UUID]float64 // Trust of other agents
	Institution        map[uuid.UUID]float64 // Institution of other agents
	Network            map[uuid.UUID]float64 // Network of other agents
	GameIterations     int32                 // Keep track of game iterations // TODO: WHAT IS THIS?
	forgivenessCounter int32                 // Keep track of how many rounds we have been forgiving an agent
	gameState          objects.IGameState    // updated by the server at every round
	megaBikeId         uuid.UUID
	bikeCounter        map[uuid.UUID]int32
	actions            []Action
	soughtColour       utils.Colour // the colour of the lootbox that the agent is currently seeking
	onBike             bool
	energyLevel        float64 // float between 0 and 1
	points             int
	forces             utils.Forces
	allocationParams   objects.ResourceAllocationParams
}

func NewBaseTeam2Biker(agentId uuid.UUID) *AgentTwo {
	color := utils.GenerateRandomColour()
	baseBiker := objects.GetBaseBiker(color, agentId)
	return &AgentTwo{
		BaseBiker:          baseBiker,
		SocialCapital:      make(map[uuid.UUID]float64),
		Trust:              make(map[uuid.UUID]float64),
		Institution:        make(map[uuid.UUID]float64),
		Network:            make(map[uuid.UUID]float64),
		GameIterations:     0,
		forgivenessCounter: 0,
		gameState:          nil,
		megaBikeId:         uuid.UUID{},
		bikeCounter:        make(map[uuid.UUID]int32),
		actions:            make([]Action, 0),
		soughtColour:       color,
		onBike:             false,
		energyLevel:        1.0,
		points:             0,
		forces:             utils.Forces{},
		allocationParams:   objects.ResourceAllocationParams{},
	}
}

const (
	TrustWeight       = 1.0
	InstitutionWeight = 0.0
	NetworkWeight     = 0.0
)

type Action struct {
	AgentID         uuid.UUID
	Action          string
	Force           utils.Forces
	GameLoop        int32
	lootBoxlocation Vector //utils.Coordinates
}

// TODO: function CalculateSocialCapital
func (a *AgentTwo) CalculateSocialCapital() {
	// Implement this method
	// Hardcode the weightings for now: Trust 1, Institution 0, Network 0
	// Calculate social capital of all agents
	// Calculate trustworthiness of all agents
	// Calculate social networks of all agents
	// Calculate institutions of all agents
	// Iterate over each agent
	for agentID, _ := range a.Trust {
		trustworthiness := a.Trust[agentID]
		institution := a.Institution[agentID]
		network := a.Network[agentID] // Assuming these values are already calculated

		a.SocialCapital[agentID] = TrustWeight*trustworthiness + InstitutionWeight*institution + NetworkWeight*network
	}
}

type Vector struct {
	X float64
	Y float64
}

func forcesToVectorConversion(force utils.Forces) Vector {
	xCoordinate := force.Pedal * float64(math.Cos(float64(math.Pi*force.Turning.SteeringForce)))
	yCoordinate := force.Pedal * float64(math.Sin(float64(math.Pi*force.Turning.SteeringForce)))

	newVector := Vector{X: xCoordinate, Y: yCoordinate}
	return newVector
}

func dotProduct(v1, v2 Vector) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func magnitude(v Vector) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func cosineSimilarity(v1, v2 Vector) float64 {
	return dotProduct(v1, v2) / (magnitude(v1) * magnitude(v2))
}

const (
	forgivenessFactor = 0.5
)

func (a *AgentTwo) updateTrustworthiness(agentID uuid.UUID, actualAction, expectedAction Vector) {
	// Calculates the cosine Similarity of actual and expected vectors. One issue is that it does not consider magnitude, only direction
	// TODO: Take magnitude into account
	similarity := cosineSimilarity(actualAction, expectedAction)

	// CosineSimilarity output ranges from -1 to 1. Need to scale it back to 0-1
	normalisedTrustworthiness := (similarity + 1) / 2

	// Moving average
	// a.Trust[agentID] = (forgivenessFactor*a.Trust[agentID]*float64(a.GameIterations) + normalisedTrustworthiness) / (float64(a.GameIterations) + 1)

	// Bad action but with high trustworthiness in prev rounds, we feel remorse and we forgive them
	if a.Trust[agentID] > normalisedTrustworthiness && a.forgivenessCounter <= 3 { // If they were trustworthy in prev rounds, we feel remorse and we forgive them
		a.forgivenessCounter++
		a.Trust[agentID] = (a.Trust[agentID]*float64(a.GameIterations) + (normalisedTrustworthiness + forgivenessFactor*(normalisedTrustworthiness-a.Trust[agentID]))) / (float64(a.GameIterations) + 1)
	} else if a.forgivenessCounter > 3 {
		// More than 3 rounds of BETRAYAL, we don't forgive them anymore...
		a.Trust[agentID] = (a.Trust[agentID]*float64(a.GameIterations) + normalisedTrustworthiness) / (float64(a.GameIterations) + 1)
	} else {
		// Good action with high trustworthiness
		a.forgivenessCounter = 0
		a.Trust[agentID] = (a.Trust[agentID]*float64(a.GameIterations) + normalisedTrustworthiness) / (float64(a.GameIterations) + 1)
	}

	fmt.Println("Trust: ", a.Trust)

}

// func (a *AgentTwo) updateInstitution(agentID uuid.UUID) float64 {

// 	// return 0.5 // This is just a placeholder value
// }

// func (a *AgentTwo) updateNetwork(agentID uuid.UUID) float64 {
// 	// return 0.5 // This is just a placeholder value
// }

// func (a *AgentTwo) calculateTrustworthiness(agentID uuid.UUID) float64 {

// 	return 0.5 // This is just a placeholder value
// }

// func (a *AgentTwo) calculateInstitution(agentID uuid.UUID) float64 {

// 	// return 0.5 // This is just a placeholder value
// }

// func (a *AgentTwo) calculateNetwork(agentID uuid.UUID) float64 {
// 	// return 0.5 // This is just a placeholder value
// }

func (a *AgentTwo) ChooseOptimalBike() uuid.UUID {
	// Implement this method
	// Calculate utility of all bikes for our own survival (remember previous actions (has space, got lootbox, direction) of all bikes so you can choose a bike to move to to max our survival chances) -> check our reputation (trustworthiness, social networks, institutions)

	// - We change the bike if an agent sees more than N agents below a social capital threshold.
	var N int32 = 3
	SocialCapitalThreshold := 0.5
	// - N and the Social Capital Threshold could be varied.

	currentBikeID := a.GetBike()
	a.gameState = a.GetGameState()

	for bikeID, bike := range a.gameState.GetMegaBikes() {
		for _, agent := range bike.GetAgents() {
			if a.SocialCapital[agent.GetID()] > SocialCapitalThreshold {
				a.bikeCounter[bikeID]++
			}
		}
	}

	if a.bikeCounter[currentBikeID] > N {
		// Stay on bike
		return currentBikeID
	} else {
		// find max bike counter in a.bikeCounter map
		// change bike to that bike
		var maxValue int32 = 0
		maxBikeID := uuid.UUID{}
		for bikeID, counter := range a.bikeCounter {
			if maxValue < counter {
				maxValue = counter
				maxBikeID = bikeID
			}
		}
		return maxBikeID
	}
}

// Failsafe: if evergy level is less than oneround in the VOID, don't change bike
// if we have a leader, then we keeop track of how many round each agent was a leader. If we are a leader, we can use this to decide if we want to change bike or not.
// TODO: Create a function to retain history of previous actions of all bikes and bikers from gamestates (Needs conformation about getting access to gamestates)
// TODO: Create a function to calculate expected gain

func (a *AgentTwo) DecideAction() objects.BikerAction {
	fmt.Println("DecideAction entering")
	// lootBoxlocation := Vector{X: 0.0, Y: 0.0} // need to change this later on (possibly need to alter the updateTrustworthiness function)
	//update agent's trustworthiness every round pretty much at the start of each epoch
	a.gameState = a.GetGameState()

	// fmt.Println("DecideAction megabikes: ", a.gameState.GetMegaBikes())
	for _, bike := range a.GetGameState().GetMegaBikes() {
		// fmt.Println("DecideAction bike: ", bike.GetID(), " ", bike.GetAgents())
		for _, agent := range bike.GetAgents() {
			// get the force for the agent with agentID in actions
			agentID := agent.GetID()
			// fmt.Println("DecideAction agentID: ", agentID)
			for _, action := range a.actions {
				fmt.Println("DecideAction action: ", action)
				if action.AgentID == agentID {
					// update trustworthiness
					a.updateTrustworthiness(agentID, forcesToVectorConversion(action.Force), action.lootBoxlocation)
				}
			}
			// a.updateTrustworthiness(agent.GetID(), forcesToVectorConversion(), lootBoxlocation)
		}
	}
	// a.gameState.GetMegaBikes()[a.GetBike()].GetAgents()[0].GetForces()
	// Check energy level, if below threshold, don't change bike
	// energyThreshold := 0.2
	// fmt.Println("OUTSIDE FOR LOOP: ", a.GetEnergyLevel(), energyThreshold, a.ChooseOptimalBike(), a.GetBike())

	// TODO: ChangeBike is broken in GameLoop
	// if (a.GetEnergyLevel() < energyThreshold) || (a.ChooseOptimalBike() == a.GetBike()) {
	// 	return objects.Pedal
	// } else {
	// 	// random for now, changeBike changes to a random uuid for now.
	// 	return objects.ChangeBike
	// }
	return objects.Pedal

	// TODO: When we have access to limbo/void then we can worry about these
	// Utility = expected gain - cost of changing bike(no of rounds in the void * energy level drain)
	// no of rounds in the void = 1 + (distance to lootbox / speed of bike)
}

// TODO: Once the MVP is complete, we can start thinking about this and then feed it into DecideForce
func (a *AgentTwo) CalcExpectedGainForLootbox(lootboxID uuid.UUID) float64 {
	// Implement this method
	// Calculate gain of going for a given lootbox(box colour and distance to it), to decide the action (e.g. pedal, brake, turn) to take

	// a.GetEnergyLevel()
	// energyLost := agent.GetForces().Pedal * utils.MovingDepletion

	//What the server uses to drain energy from us for moving
	// for _, agent := range s.megaBikes[bikeID].GetAgents() {
	// 	agent.DecideForce(direction)
	// 	// deplete energy
	// 	energyLost := agent.GetForces().Pedal * utils.MovingDepletion
	// 	agent.UpdateEnergyLevel(-energyLost)
	// }

	// energy from lootbox - energy from pedalling

	currLocation := a.GetLocation()
	targetPos := a.gameState.GetLootBoxes()[lootboxID].GetPosition()

	deltaX := targetPos.X - currLocation.X
	deltaY := targetPos.Y - currLocation.Y
	angle := math.Atan2(deltaX, deltaY)
	angleInDegrees := angle * math.Pi / 180

	// Default BaseBiker will always
	turningDecision := utils.TurningDecision{
		SteerBike:     true,
		SteeringForce: angleInDegrees,
	}

	lootboxForces := utils.Forces{
		Pedal:   utils.BikerMaxForce,
		Brake:   0.0,
		Turning: turningDecision,
	}

	energyFromLootbox := a.gameState.GetLootBoxes()[lootboxID].GetTotalResources()
	energyToPedal := lootboxForces.Pedal * utils.MovingDepletion //Moving depleteion is a constant 1 atm. TODO: Change this to a variable and check how Pedal relates to energy depletion

	expectedGain := energyFromLootbox - energyToPedal

	return expectedGain
}

func (a *AgentTwo) GetPreviousAction() {
	// -> get previous action of all bikes and bikers from last 5 gamestates

	// nearestLoot := a.nearestLoot()
	// currentLootBoxes := a.gameState.GetLootBoxes()
	// lootBoxlocation := currentLootBoxes[nearestLoot].GetPosition()
	lootBoxlocation_vector := Vector{X: 0.0, Y: 0.0} // need to change this later on (possibly need to alter the updateTrustworthiness function)
	//update agent's trustworthiness every round pretty much at the start of each epoch
	for _, bike := range a.gameState.GetMegaBikes() {
		for _, agent := range bike.GetAgents() {
			// Record the action
			action := Action{
				AgentID:         agent.GetID(),
				Action:          "DecideForce",
				Force:           agent.GetForces(),
				GameLoop:        a.GameIterations, // record the game loop number
				lootBoxlocation: lootBoxlocation_vector,
			}
			// If we have more than 5 actions, remove the oldest one
			if len(a.actions) >= 5 {
				a.actions = a.actions[1:]
			}

			// Append the new action
			a.actions = append(a.actions, action)
		}
	}
}

// returns the nearest lootbox with respect to the agent's bike current position
// in the MVP this is used to determine the pedalling forces as all agent will be
// aiming to get to the closest lootbox by default
func (a *AgentTwo) nearestLoot() uuid.UUID {
	currLocation := a.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestBox uuid.UUID
	var currDist float64
	for _, loot := range a.gameState.GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return nearestBox
}

// To overwrite the BaseBiker's DecideForce method in order to record all the previous actions of all bikes (GetForces) and bikers from gamestates
func (a *AgentTwo) DecideForce(direction uuid.UUID) {
	fmt.Println("DecideForce entering")
	// Pedal, Brake, Turning
	// GetPreviousAction() -> get previous action of all bikes and bikers from gamestates
	// GetVotedLootbox() -> get voted lootbox from gamestates
	// GetOptimalLootbox() -> get optimal lootbox for ourself from gamestates
	// probabilityOfConformity = selfSocialCapital
	// Generate random number between 0 and 1
	// if random number < probabilityOfConformity, then conform
	// else, don't conform

	// CalculateForceAndSteer(Lootbox) -> calculate force and steer towards lootbox
	// set a.forces.steerbike == True

	a.GetPreviousAction()
	a.gameState = a.GetGameState()

	// NEAREST BOX STRATEGY (MVP)
	currLocation := a.GetLocation()
	nearestLoot := a.nearestLoot()
	currentLootBoxes := a.gameState.GetLootBoxes()
	fmt.Println("DecideForce entering")
	fmt.Println("nearestLoot: ", nearestLoot)
	// fmt.Println("currentLootBoxes: ", currentLootBoxes)
	fmt.Println("currLocation: ", currLocation, " bike: ", a.GetBike(), " energy: ", a.GetEnergyLevel(), " points: ", a.GetPoints())

	// Check if there are lootboxes available and move towards closest one
	if len(currentLootBoxes) > 0 {
		targetPos := currentLootBoxes[nearestLoot].GetPosition()
		fmt.Println("targetPos: ", targetPos)
		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaX, deltaY)
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		fmt.Println(a.gameState.GetMegaBikes()[a.GetBike()].GetOrientation())
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - a.gameState.GetMegaBikes()[a.GetBike()].GetOrientation(),
		}
		// fmt.Println("turningDecision: ", turningDecision)

		nearestBoxForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		a.forces = nearestBoxForces
		a.SetForces(a.forces)
	} else { // otherwise move away from audi
		audiPos := a.GetGameState().GetAudi().GetPosition()

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
		a.forces = escapeAudiForces
	}

	a.GameIterations++
	fmt.Println("GameIterations ", a.GameIterations)
	// fmt.Println(actions)
}

func (a *AgentTwo) ChangeBike() uuid.UUID {
	// Implement this method
	// Stage 1 called by DecideAction when
	// proposal to change bike to a goal bike
	return uuid.UUID{}
}
