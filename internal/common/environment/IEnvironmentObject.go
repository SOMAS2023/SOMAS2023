package ienvironmentobject

/*

The IEnvironementObject is an interface class that all objects (including agents) must implement.

*/

import (
	"github.com/google/uuid"
)

type IEnvironementObject[T any] struct {
	// returns the unique ID of an object
	GetID() uuid.UUID

	// returns the current coordinates of Agent
	GetPosition() [2]float64
}