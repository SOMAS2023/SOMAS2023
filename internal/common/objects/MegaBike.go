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
	KickOutAgent(bikerId uuid.UUID)
}

// MegaBike will have the following forces
type MegaBike struct {
	*PhysicsObject
	agents []IBaseBiker
	kickedOutCount    int
	currentLeaderUUID uuid.UUID
}

// GetMegaBike is a constructor for MegaBike that initializes it with a new UUID and default position.
func GetMegaBike() *MegaBike {
	return &MegaBike{
		PhysicsObject: GetPhysicsObject(utils.MassBike),
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
			numOfSteeringAgents += 1
			totalTurning += float64(turningDecision.SteeringForce)
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

// set a current leader
func (mb *MegaBike) SetCurrentLeader(leaderUUID uuid.UUID) {
	mb.currentLeaderUUID = leaderUUID
}

// get the current leader
func (mb *MegaBike) GetCurrentLeader() uuid.UUID {
	return mb.currentLeaderUUID
}

// remove an agent from the bike and increase the kickedOutCount
func (mb *MegaBike) KickOutAgent() uuid.UUID {
	voteCount := make(map[uuid.UUID]int)

	// Count votes for each agent
	for _, agent := range mb.agents {
		votedAgentID := agent.VoteForKickout(mb.agents)
		if votedAgentID != uuid.Nil {
			voteCount[votedAgentID]++
		}
	}

	// Find the agent with the most votes
	var maxVotes int
	var kickedOutAgentId uuid.UUID
	for agentID, votes := range voteCount {
		if votes > maxVotes && votes > len(voteCount)/2{
			maxVotes = votes
			kickedOutAgentId = agentID
		}
	}

	// Kick out the agent with the most votes
	if kickedOutAgentId != uuid.Nil {
		for i, biker := range mb.agents {
			if biker.GetID() == kickedOutAgentId {
				mb.agents = append(mb.agents[:i], mb.agents[i+1:]...)
				mb.kickedOutCount++
				break
			}
		}
	}

	return kickedOutAgentId
}

