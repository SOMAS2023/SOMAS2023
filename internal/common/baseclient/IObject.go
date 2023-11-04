package iobject

/*

The IObject is an interface class that all objects (not agents) must inherit

Every object such as Audi, LootBox and Bike must inherit IObject

*/

import (
	"github.com/google/uuid"
)

type IObject[T any] struct {
	// returns the unique ID of an object
	GetID() uuid.UUID

	// returns the current Coordinates of Agent
	GetCoordinates() [2]float64

}