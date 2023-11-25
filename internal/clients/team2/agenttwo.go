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

	// "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

type AgentTwo struct {
	// BaseBiker represents a basic biker agent.
	*objects.BaseBiker
	// CalculateSocialCapitalOtherAgent: (trustworthiness - cosine distance, social networks - friends, institutions - num of rounds on a bike)
	SocialCapital      map[uuid.UUID]float64 // Social Captial of other agents
	Trust              map[uuid.UUID]float64 // Trust of other agents
	Institution        map[uuid.UUID]float64 // Institution of other agents
	Network            map[uuid.UUID]float64 // Network of other agents
	GameIterations     int32                 // Keep track of game iterations
	forgivenessCounter int32                 // Keep track of how many rounds we have been forgiving an agent
	gameState          objects.IGameState    // updated by the server at every round
	// megaBikeId uuid.UUID
	bikeCounter map[uuid.UUID]int32
}

const (
	TrustWeight       = 1.0
	InstitutionWeight = 0.0
	NetworkWeight     = 0.0
)

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
	// lootBoxlocation := Vector{X: 0.0, Y: 0.0} // need to change this later on (possibly need to alter the updateTrustworthiness function)
	//update agent's trustworthiness every round pretty much at the start of each epoch
	for _, bike := range a.gameState.GetMegaBikes() {
		for _, agent := range bike.GetAgents() {
			// get the force for the agent with agentID in actions
			agentID := agent.GetID()
			for _, action := range actions {
				if action.AgentID == agentID {
					// update trustworthiness
					a.updateTrustworthiness(agentID, forcesToVectorConversion(action.Force), action.lootBoxlocation)
				}
			}
			// a.updateTrustworthiness(agent.GetID(), forcesToVectorConversion(), lootBoxlocation)
		}
	}

	// Check energy level, if below threshold, don't change bike
	energyThreshold := 0.2

	if (a.GetEnergyLevel() < energyThreshold) || (a.ChooseOptimalBike() == a.GetBike()) {
		return objects.Pedal
	} else {
		// random for now, changeBike changes to a random uuid for now.
		return objects.ChangeBike
	}

	// TODO: When we have access to limbo/void then we can worry about these
	// Utility = expected gain - cost of changing bike(no of rounds in the void * energy level drain)
	// no of rounds in the void = 1 + (distance to lootbox / speed of bike)
}

// TODO: Once the MVP is complete, we can start thinking about this and then feed it into DecideForce
// func (a *AgentTwo) CalcExpectedGainForLootbox(lootboxID uuid.UUID) float64 {
// 	// Implement this method
// 	// Calculate gain of going for a given lootbox(box colour and distance to it), to decide the action (e.g. pedal, brake, turn) to take

// 	// a.GetEnergyLevel()
// 	// energyLost := agent.GetForces().Pedal * utils.MovingDepletion

// 	//What the server uses to drain energy from us for moving
// 	// for _, agent := range s.megaBikes[bikeID].GetAgents() {
// 	// 	agent.DecideForce(direction)
// 	// 	// deplete energy
// 	// 	energyLost := agent.GetForces().Pedal * utils.MovingDepletion
// 	// 	agent.UpdateEnergyLevel(-energyLost)
// 	// }

// 	// energy from lootbox - energy from pedalling

// 	currLocation := a.GetLocation()
// 	targetPos := a.gameState.GetLootBoxes()[lootboxID].GetPosition()

// 	deltaX := targetPos.X - currLocation.X
// 	deltaY := targetPos.Y - currLocation.Y
// 	angle := math.Atan2(deltaX, deltaY)
// 	angleInDegrees := angle * math.Pi / 180

// 	// Default BaseBiker will always
// 	turningDecision := utils.TurningDecision{
// 		SteerBike:     true,
// 		SteeringForce: angleInDegrees,
// 	}

// 	lootboxForces := utils.Forces{
// 		Pedal:   utils.BikerMaxForce,
// 		Brake:   0.0,
// 		Turning: turningDecision,
// 	}

// 	energyFromLootbox := a.gameState.GetLootBoxes()[lootboxID].GetTotalResources()
// 	energyToPedal := lootboxForces.Pedal * utils.MovingDepletion //Moving depleteion is a constant 1 atm. TODO: Change this to a variable and check how Pedal relates to energy depletion

// 	expectedGain := energyFromLootbox - energyToPedal

// 	return expectedGain
// }

// func (a *AgentTwo) deci

type Action struct {
	AgentID         uuid.UUID
	Action          string
	Force           utils.Forces
	GameLoop        int
	lootBoxlocation Vector
}

var actions []Action
var gameLoopNumber int // TODO: find a function to increment the gameLoopNumber

// To overwrite the BaseBiker's DecideForce method in order to record all the previous actions of all bikes (GetForces) and bikers from gamestates
func (a *AgentTwo) DecideForce() {
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

	lootBoxlocation := Vector{X: 0.0, Y: 0.0} // need to change this later on (possibly need to alter the updateTrustworthiness function)
	//update agent's trustworthiness every round pretty much at the start of each epoch
	for _, bike := range a.gameState.GetMegaBikes() {
		for _, agent := range bike.GetAgents() {
			// Record the action
			action := Action{
				AgentID:         agent.GetID(),
				Action:          "DecideForce",
				Force:           agent.GetForces(),
				GameLoop:        gameLoopNumber, // record the game loop number
				lootBoxlocation: lootBoxlocation,
			}
			// If we have more than 5 actions, remove the oldest one
			if len(actions) >= 5 {
				actions = actions[1:]
			}

			// Append the new action
			actions = append(actions, action)
		}
	}
	// fmt.Println(actions)
}

func (a *AgentTwo) ChangeBike() uuid.UUID {
	// Implement this method
	// Stage 1 called by DecideAction when
	// proposal to change bike to a goal bike
	return uuid.UUID{}
}

// NOTES ------------------------------------------------------------

// 1) Decide on giving agent to access the gameState -> getGameState()

// 2) Those function should only be called by the server, not by the agent

// func (a *AgentTwo) UpdateEnergyLevel(energyLevel float64) { // TODO: TO BE CHECKED WITH TEAM LEADERS!!!!!
// 	// Implement this
// 	// should not be able to call this, server calls this
// }

// func (a *AgentTwo) GetResourceAllocationParams() objects.ResourceAllocationParams { // TODO: TO BE CHECKED WITH TEAM LEADERS!!!!!
// 	// Implement this method
// 	// SERVER CALLS THIS, agent should just ask for a specific demand
// 	// STAGE 4: how we want to proporsion the energy bar distribution
// 	return objects.ResourceAllocationParams{}
// }

// func (a *AgentTwo) SetAllocationParameters(params objects.ResourceAllocationParams) { // TODO: TO BE CHECKED WITH TEAM LEADERS!!!!!
// 	// Implement this method
// 	// should not be able to call this, server calls this
// }

// Founding processes: need to consider this for the wednesday meeting.
