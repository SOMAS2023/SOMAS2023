package objects

import (
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type IMegaBike interface {
	IPhysicsObject
	AddAgent(biker IBaseBiker)
	RemoveAgent(bikerId uuid.UUID)
	GetAgents() []IBaseBiker
	UpdateMass()
	KickOutAgent() map[uuid.UUID]int
	GetGovernance() utils.Governance
	GetRuler() uuid.UUID
	SetGovernance(governance utils.Governance)
	SetRuler(ruler uuid.UUID)
}

// MegaBike will have the following forces
type MegaBike struct {
	*PhysicsObject
	agents         []IBaseBiker
	kickedOutCount int
	governance     utils.Governance
	ruler          uuid.UUID
}

// GetMegaBike is a constructor for MegaBike that initializes it with a new UUID and default position.
func GetMegaBike() *MegaBike {
	return &MegaBike{
		PhysicsObject: GetPhysicsObject(utils.MassBike),
		governance:    utils.Democracy,
		ruler:         uuid.Nil,
	}
}

// adds
func (mb *MegaBike) AddAgent(biker IBaseBiker) {
	mb.agents = append(mb.agents, biker)
}

// Remove agent from bike, given its ID
func (mb *MegaBike) RemoveAgent(bikerId uuid.UUID) {
	// Create a new slice to store the updated agents
	var updatedAgents []IBaseBiker

	// Iterate through the agents and copy them to the updatedAgents slice
	for _, agent := range mb.agents {
		if agent.GetID() != bikerId {
			updatedAgents = append(updatedAgents, agent)
		}
	}

	// Replace the mb.agents slice with the updatedAgents slice
	mb.agents = updatedAgents
}

func (mb *MegaBike) GetAgents() []IBaseBiker {
	return mb.agents
}

// Calculate the mass of the bike with all it's agents
func (mb *MegaBike) UpdateMass() {
	mass := utils.MassBike
	mass += float64(len(mb.agents))
	mb.mass = mass
}

// Calculates and returns the total force of the Megabike based on the Biker's force
func (mb *MegaBike) UpdateForce() {
	if len(mb.agents) == 0 {
		mb.force = 0.0
	}
	totalPedal := 0.0
	totalBrake := 0.0
	for _, agent := range mb.agents {
		force := agent.GetForces()

		if force.Pedal != 0 {
			totalPedal += float64(force.Pedal)
		} else {
			totalBrake += float64(force.Brake)
		}
	}
	mb.force = (float64(totalPedal) - float64(totalBrake))
}

// Calculates the final orientation of the Megabike, between -1 and 1 (-180° to 180°), given the Biker's Turning forces
func (mb *MegaBike) UpdateOrientation() {
	totalTurning := 0.0
	numOfSteeringAgents := 0
	for _, agent := range mb.agents {
		// If agents do not want to steer, they must set their TurningDecision.SteerBike to false and their steering
		// will not have an impact on the direction of the bike.
		turningDecision := agent.GetForces().Turning
		if turningDecision.SteerBike {
			// Only dictators can steer if Governance is set to dictatorship
			if (mb.governance != utils.Dictatorship) || (agent.GetID() == mb.ruler) {
				numOfSteeringAgents += 1
				totalTurning += float64(turningDecision.SteeringForce)
			}
		}
	}
	// Do not update orientation if no biker want to steer
	if numOfSteeringAgents > 0 {
		averageTurning := totalTurning / float64(numOfSteeringAgents)
		mb.orientation += (averageTurning)
	}
	// ensure the orientation wraps around if it exceeds the range 1.0 or -1.0

	if mb.orientation > 1.0 {
		mb.orientation -= 2
	} else if mb.orientation < -1.0 {
		mb.orientation += 2
	}
}

// get the count of kicked out agents
func (mb *MegaBike) GetKickedOutCount() int {
	return mb.kickedOutCount
}

func (mb *MegaBike) KickOutAgent() map[uuid.UUID]int {
	voteCount := make(map[uuid.UUID]int)
	// Count votes for each agent
	for _, agent := range mb.agents {
		agentVotes := agent.VoteForKickout() // Assuming this now returns map[uuid.UUID]int
		for agentID, votes := range agentVotes {
			voteCount[agentID] += votes
		}
	}

	// Find all agents with votes > half the number of agents
	agentsToKickOut := make(map[uuid.UUID]int)
	for agentID, votes := range voteCount {
		if votes > len(mb.agents)/2 {
			agentsToKickOut[agentID] = votes
		}
	}

	mb.kickedOutCount += len(agentsToKickOut)

	return agentsToKickOut
}

func (mb *MegaBike) GetGovernance() utils.Governance {
	return mb.governance
}

func (mb *MegaBike) GetRuler() uuid.UUID {
	return mb.ruler
}

func (mb *MegaBike) SetGovernance(governance utils.Governance) {
	mb.governance = governance
}

func (mb *MegaBike) SetRuler(ruler uuid.UUID) {
	mb.ruler = ruler
}
