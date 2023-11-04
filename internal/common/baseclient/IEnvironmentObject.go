package ienvironmentobject

/*

The IEnvironementObject is an interface class that all objects (including agents) must implement.

*/

import (
	"github.com/google/uuid"
)

type IEnvironementObject struct {
	// Returns if object is alive or dead
	IsAlive() bool

	// returns the unique ID of an object
	GetID() uuid.UUID

	// returns the current coordinates of Agent
	GetPosition() [2]float64

	// Gets the current force, which will be called by the server after the agents have called UpdateForces()
	// during the UpdateAgentInternalState() method call.
	// The forces are[pedal, brake, turning]
	// The pedal and brake forces is a float from 0.0 to 1.0, the turning force is a float from -1.0 (90° left) to 1.0 (90° right)
	GetForces() [3]float64

	// After making a decision, the agent must update their force each round.
	// The forces are forces["pedalForce"], forces["brakeForce"] and forces["turningForce"]
	// The pedal and brake forces is a float from 0.0 to 1.0, the turning force is a float from -1.0 (90° left) to 1.0 (90° right)
	UpdateForces()

}