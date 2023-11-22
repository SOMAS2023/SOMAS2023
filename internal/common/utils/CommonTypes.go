package utils

type Colour int

const (
	Red Colour = iota
	Green
	Blue
	Yellow
	Orange
	Purple
	Pink
	Brown
	Gray
	White
	NumOfColours // add a sentinel for counting the number of colours
)

type TurningDecision struct {
	SteerBike     bool
	SteeringForce float64
}

type Forces struct {
	Pedal float64 // Pedal is a force from 0-1 where 1 is 100% power
	Brake float64 // Brake is a force from 0-1 opposing the direction of travel (bike cannot go backwards)

	// If agents do not want to steer, they must set their TurningDecision.SteerBike to false and their steering will not have an impact on the direction of the bike.
	// TurningDecision.SteeringForce is a force from -1 to 1 which maps to -180° to 180°.
	Turning TurningDecision
}

type Coordinates struct {
	X float64
	Y float64
}

type PhysicalState struct {
	Position     Coordinates
	Acceleration float64
	Velocity     float64
	Mass         float64
}

type Governance int

const (
	Democracy Governance = iota
	Leadership
	Dictatorship
	Invalid
)
