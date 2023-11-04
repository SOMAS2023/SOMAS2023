package environmentobject

/*

	EnvironmentObject is the implentation of IEnvironmentObject. All objects (including agents)

*/

import "github.com/google/uuid"

// "SOMAS2023/internal/common/baseclient"

type EnvironmentObject struct {
	id       uuid.UUID
	isAlive  bool
	position [2]float64
	forces   [3]float64
}

func (e *EnvironmentObject) IsAlive() bool {
	return e.isAlive
}

func (e *EnvironmentObject) GetID() uuid.UUID {
	return e.id
}

func (e *EnvironmentObject) GetPosition() [2]float64 {
	return e.position
}

func (e *EnvironmentObject) GetForces() [3]float64 {
	return e.forces
}
