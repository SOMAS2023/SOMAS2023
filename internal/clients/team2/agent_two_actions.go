package team2

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"fmt"
	"math"

	"github.com/google/uuid"
)

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
					// Needs to be updated so that a.NearLootbox() is replaced with the lootbox location that the agent says that they're going for
					a.updateReputation(agentID, a.GetOptimalLootbox(), a.nearestLoot())
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
