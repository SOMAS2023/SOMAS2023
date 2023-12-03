package team2

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"fmt"
	"math"
	"sort"

	"github.com/google/uuid"
)

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
				// fmt.Println("DecideAction action: ", action)
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

// find the number of agents on the bike
func (a *AgentTwo) GetAgentNum(bikeID uuid.UUID) float64 {
	var count = 0.0
	for _, agent := range a.gameState.GetMegaBikes()[bikeID].GetAgents() {
		// agentID := agent.GetID()
		_ = agent.GetID()
		count++
	}
	return count
}

// TODO: Once the MVP is complete, we can start thinking about this and then feed it into DecideForce
func (a *AgentTwo) CalcExpectedGainForLootbox(lootboxID uuid.UUID) float64 {
	// Implement this method
	// Calculate gain of going for a given lootbox(box colour and distance to it), to decide the action (e.g. pedal, brake, turn) to take

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

	// assumes loot is equally divided among all agents on the bike
	energyFromLootbox := a.gameState.GetLootBoxes()[lootboxID].GetTotalResources() / a.GetAgentNum(a.GetBike())
	energyToPedal := lootboxForces.Pedal * utils.MovingDepletion //Moving depleteion is a constant 1 atm. TODO: Change this to a variable and check how Pedal relates to energy depletion

	expectedGain := energyFromLootbox - energyToPedal

	return expectedGain
}

func (a *AgentTwo) GetPreviousAction() {
	// -> get previous action of all bikes and bikers from last 5 gamestates

	// nearestLoot := a.nearestLoot()
	// currentLootBoxes := a.gameState.GetLootBoxes()
	// lootBoxlocation := currentLootBoxes[nearestLoot].GetPosition()
	lootBoxlocation_vector := ForceVector{X: 0.0, Y: 0.0} // need to change this later on (possibly need to alter the updateTrustworthiness function)
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

func (a *AgentTwo) AvoidOwdi(lootboxID uuid.UUID) bool {
	audiPos := a.GetGameState().GetAudi().GetPosition()
	lootboxPos := a.GetGameState().GetLootBoxes()[lootboxID].GetPosition()
	AudiRange := utils.Coordinates{
		X: 5.0,
		Y: 5.0,
	}
	fmt.Println("AudiRange: ", AudiRange)
	fmt.Println("audiPosX-: ", audiPos.X-AudiRange.X)
	fmt.Println("audiPosX+: ", audiPos.X+AudiRange.X)
	fmt.Println("audiPosY-: ", audiPos.Y-AudiRange.Y)
	fmt.Println("audiPosY+: ", audiPos.Y+AudiRange.Y)

	if audiPos.X-AudiRange.X < 0 || audiPos.Y-AudiRange.Y < 0 || audiPos.X+AudiRange.X > utils.GridWidth || audiPos.Y+AudiRange.Y > utils.GridHeight {

		AudiRange = utils.Coordinates{
			X: 0.0,
			Y: 0.0,
		}
	}
	// if we are near an owdi we want to avoid it
	if a.GetLocation().X > audiPos.X-AudiRange.X && a.GetLocation().X < audiPos.X+AudiRange.X && a.GetLocation().Y > audiPos.Y-AudiRange.Y && a.GetLocation().Y < audiPos.Y+AudiRange.Y {
		fmt.Println("avoiding owdi")
		return true
	}

	// if the lootbox is in the owdi's range we avoid it

	if lootboxPos.X > audiPos.X-AudiRange.X && lootboxPos.X < audiPos.X+AudiRange.X && lootboxPos.Y > audiPos.Y-AudiRange.Y && lootboxPos.Y < audiPos.Y+AudiRange.Y {
		fmt.Println("avoiding owdi")

		return true
	}

	return false
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

func (a *AgentTwo) GetOptimalLootbox() uuid.UUID {
	// highestGain := 0.0
	// var lootboxID = uuid.UUID{}
	// m := make(map[uuid.UUID]float64)
	var topLootboxes []struct {
		ID    uuid.UUID
		Gain  float64
		Color utils.Colour
	}
	// agentColor := a.GetColour().String()
	var top3Lootboxes []struct {
		ID    uuid.UUID
		Gain  float64
		Color utils.Colour
	}

	for _, lootbox := range a.gameState.GetLootBoxes() {
		expectedGain := a.CalcExpectedGainForLootbox(lootbox.GetID())

		// 	fmt.Println("lootbox id: ", lootbox.GetID())
		// 	fmt.Println("lootbox reward: ", lootbox.GetTotalResources())
		// 	fmt.Println("gain of lootbox: ", expectedGain)
		// 	fmt.Println("lootbox position: ", lootbox.GetPosition())

		// 	// find the highest gain lootbox
		// 	if a.CalcExpectedGainForLootbox(lootbox.GetID()) > highestGain {
		// 		highestGain = a.CalcExpectedGainForLootbox(lootbox.GetID())
		// 		lootboxID = lootbox.GetID()
		// 	}

		// }

		// return lootboxID
		topLootboxes = append(topLootboxes, struct {
			ID    uuid.UUID
			Gain  float64
			Color utils.Colour
		}{ID: lootbox.GetID(), Gain: expectedGain, Color: lootbox.GetColour()})

		// Sort the lootboxes by gain from the mapping
	}
	sort.Slice(topLootboxes, func(i, j int) bool {
		return topLootboxes[i].Gain > topLootboxes[j].Gain
	})

	top3Lootboxes = topLootboxes
	if len(topLootboxes) > 3 {
		top3Lootboxes = topLootboxes[:3]
	}
	fmt.Println("top3: ", top3Lootboxes)
	for _, top3 := range top3Lootboxes {
		if top3.Color == a.GetColour() {
			return top3.ID
		}
	}

	return top3Lootboxes[0].ID
}

// To overwrite the BaseBiker's DecideForce method in order to record all the previous actions of all bikes (GetForces) and bikers from gamestates
func (a *AgentTwo) DecideForce(direction uuid.UUID) {

	a.votedDirection = direction
	fmt.Println("DecideForce entering")
	fmt.Println("agent energy before: ", a.GetEnergyLevel())
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

	// FIND THE OPTIMAL LOOTBOX AND MOVE TOWARDS IT
	// nearestLoot = a.GetOptimalLootbox()
	nearestLoot = a.GetOptimalLootbox()
	// Check if there are lootboxes available and move towards closest one
	// if len(currentLootBoxes) > 0 {
	if !a.AvoidOwdi(nearestLoot) {
		targetPos := currentLootBoxes[nearestLoot].GetPosition()

		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - a.gameState.GetMegaBikes()[a.GetBike()].GetOrientation(),
		}

		nearestBoxForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		a.SetForces(nearestBoxForces)
	} else { // otherwise move away from audi
		audiPos := a.GetGameState().GetAudi().GetPosition()

		deltaX := audiPos.X - currLocation.X
		deltaY := audiPos.Y - currLocation.Y

		// Steer in opposite direction to audi
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi

		// Steer in opposite direction to audi
		var flipAngle float64
		if normalisedAngle < 0.0 {
			flipAngle = normalisedAngle + 1.0
		} else if normalisedAngle > 0.0 {
			flipAngle = normalisedAngle - 1.0
		}

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: flipAngle - a.gameState.GetMegaBikes()[a.megaBikeId].GetOrientation(),
		}

		escapeAudiForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		a.SetForces(escapeAudiForces)
	}

	a.GameIterations++
	fmt.Println("GameIterations ", a.GameIterations)
	fmt.Println("agent energy after: ", a.GetEnergyLevel())
	// fmt.Println(actions)
}

func (a *AgentTwo) ChangeBike() uuid.UUID {
	// Implement this method
	// Stage 1 called by DecideAction when
	// proposal to change bike to a goal bike
	return uuid.UUID{}
}
