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

	// "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type AgentTwo struct {
	// BaseBiker represents a basic biker agent.
	*objects.BaseBiker
	// CalculateSocialCapitalOtherAgent: (trustworthiness - cosine distance, social networks - friends, institutions - num of rounds on a bike)
	// SocialCapital 		[map[uuid.UUID]float64 // map of other agents' social capital]
	// Trust 				[map[uuid.UUID]float64 // map of other agents' social capital]
	// Institution 			[map[uuid.UUID]float64 // map of other agents' social capital]
	// Network 				[map[uuid.UUID]float64 // map of other agents' social capital]
	// other fields...
}

// TODO: function CalculateSocialCapital
func (a *AgentTwo) CalculateSocialCapital() {
	// Implement this method
	// Hardcode the weightings for now: Trust 1, Institution 0, Network 0
	// Calculate social capital of all agents
	// Calculate trustworthiness of all agents
	// Calculate social networks of all agents
	// Calculate institutions of all agents
}

func (a *AgentTwo) ChangeBikeCalcUtility() {
	// Implement this method
	// Calculate utility of all bikes for our own survival (remember previous actions (has space, got lootbox, direction) of all bikes so you can choose a bike to move to to max our survival chances) -> check our reputation (trustworthiness, social networks, institutions)
}

// Failsafe: if evergy level is less than oneround in the VOID, don't change bike
// if we have a leader, then we keeop track of how many round each agent was a leader. If we are a leader, we can use this to decide if we want to change bike or not.
// TODO: Create a function to retain history of previous actions of all bikes and bikers from gamestates (Needs conformation about getting access to gamestates)
// TODO: Create a function to calculate expected gain
func (a *AgentTwo) CalcExpectedGainForLootbox() {
	// Implement this method
	// Calculate gain of going for a given lootbox(box colour and distance to it), to decide the action (e.g. pedal, brake, turn) to take
}

func (a *AgentTwo) DecideAction() objects.BikerAction {
	// Implement this method
	// Check energy level, if below threshold, don't change bike
	// Calculate expected gain for each bike
	// Utility = expected gain - cost of changing bike(no of rounds in the void * energy level drain)
	// no of rounds in the void = 1 + (distance to lootbox / speed of bike)
	return objects.Pedal
}

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
