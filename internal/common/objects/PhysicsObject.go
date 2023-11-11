package objects

/*

The IPhysicsObject is an interface class that all moving objects (Biker and Audi) must implement.
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
	GetMass() float64
	GetAcceleration() float64
	GetVelocity() float64
	GetOrientation() float64
	Move()
	CheckForCollision(otherObject IPhysicsObject) bool
}

type PhysicsObject struct {
	id           uuid.UUID
	coordinates  utils.Coordinates
	mass         float64
	acceleration float64
	velocity     float64
	orientation  float64
}

// returns the unique ID of the object
func (po *PhysicsObject) GetID() uuid.UUID {
	return po.id
}

// returns the current coordinates of the object
func (po *PhysicsObject) GetPosition() utils.Coordinates {
	return po.coordinates
}

func (po *PhysicsObject) GetMass() float64 {
	return po.mass
}

func (po *PhysicsObject) GetAcceleration() float64 {
	return po.acceleration
}

func (po *PhysicsObject) GetVelocity() float64 {
	return po.velocity
}

func (po *PhysicsObject) GetOrientation() float64 {
	return po.orientation
}

func (po *PhysicsObject) Move() {}

// this will be used to check if a MegaBike has looted a LootBok or if the Audi has collided with a MegaBike
func (po *PhysicsObject) CheckForCollision(otherObject IPhysicsObject) bool {
	otherPos := otherObject.GetPosition()
	distance := math.Sqrt(math.Pow(otherPos.X-po.coordinates.X, 2) + math.Pow(otherPos.Y-po.coordinates.Y, 2))
	if distance < utils.CollisionThreshold {
		return true
	} else {
		return false
	}
}

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
