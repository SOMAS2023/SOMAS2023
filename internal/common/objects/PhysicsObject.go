package objects

/*

The IPhysicsObject is an interface class that all moving objects (Biker and Awdi) must implement.
The Bikers/Agents will not need to implement this interface

*/

import (
	utils "SOMAS2023/internal/common/utils"

	"math"

	"github.com/google/uuid"
)

type IPhysicsObject interface {
	// returns the unique ID of the object
	GetID() uuid.UUID
	// returns the current coordinates of the object
	GetPosition() utils.Coordinates
	GetVelocity() float64
	GetOrientation() float64
	GetForce() float64
	GetPhysicalState() utils.PhysicalState

	// Server must set these variables since it updates the gamestate
	SetPhysicalState(state utils.PhysicalState)

	// This method will update the force of the PhysicsObject based on the current GameState.
	// I.e. for MegaBike, force will be cacluated from the bikers
	// For the awdi, force will be calculated from the target MegaBike
	UpdateForce()
	// Similar to UpdateForce, this will update the desired orientation for the PhysicsObject,
	// based on the current GameState
	UpdateOrientation()
	CheckForCollision(otherObject IPhysicsObject) bool
}

type PhysicsObject struct {
	id           uuid.UUID
	coordinates  utils.Coordinates
	mass         float64
	acceleration float64
	velocity     float64
	orientation  float64
	force        float64
}

// returns the unique ID of the object
func (po *PhysicsObject) GetID() uuid.UUID {
	return po.id
}

// returns the current coordinates of the object
func (po *PhysicsObject) GetPosition() utils.Coordinates {
	return po.coordinates
}

func (po *PhysicsObject) GetVelocity() float64 {
	return po.velocity
}

func (po *PhysicsObject) GetOrientation() float64 {
	return po.orientation
}

func (po *PhysicsObject) GetForce() float64 {
	return po.force
}

func (po *PhysicsObject) GetPhysicalState() utils.PhysicalState {
	return utils.PhysicalState{
		Position:     po.coordinates,
		Acceleration: po.acceleration,
		Velocity:     po.velocity,
		Mass:         po.mass,
	}
}

func (po *PhysicsObject) SetPhysicalState(state utils.PhysicalState) {
	po.mass = state.Mass
	po.coordinates = state.Position
	po.acceleration = state.Acceleration
	po.velocity = state.Velocity
}

// this will be used to check if a MegaBike has looted a LootBok or if the Awdi has collided with a MegaBike
func (po *PhysicsObject) CheckForCollision(otherObject IPhysicsObject) bool {
	otherPos := otherObject.GetPosition()
	distance := math.Sqrt(math.Pow(otherPos.X-po.coordinates.X, 2) + math.Pow(otherPos.Y-po.coordinates.Y, 2))
	if distance < utils.CollisionThreshold {
		return true
	} else {
		return false
	}
}

func (po *PhysicsObject) UpdateForce() {}

func (po *PhysicsObject) UpdateOrientation() {}

func GetPhysicsObject(mass float64) *PhysicsObject {
	return &PhysicsObject{
		id:           uuid.New(),
		coordinates:  utils.GenerateRandomCoordinates(),
		mass:         mass,
		acceleration: 0.0,
		velocity:     0.0,
		orientation:  0.0,
	}
}
