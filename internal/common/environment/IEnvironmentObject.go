package environment

/*

The IEnvironementObject is an interface class that all objects (including agents) must implement.

*/

import (
	utils "SOMAS2023/internal/common/utils"

	"github.com/google/uuid"
)

type IEnvironementObject interface {
	// returns the unique ID of the object
	GetID() uuid.UUID

	// returns the current coordinates of the object
	GetPosition() utils.Coordinates
}
