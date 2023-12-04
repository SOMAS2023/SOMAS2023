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

func (c Colour) String() string {
	switch c {
	case Red:
		return "red"
	case Green:
		return "green"
	case Blue:
		return "blue"
	case Yellow:
		return "yellow"
	case Orange:
		return "orange"
	case Purple:
		return "purple"
	case Pink:
		return "pink"
	case Brown:
		return "brown"
	case Gray:
		return "gray"
	case White:
		return "white"
	default:
		return "unknown"
	}
}

type TurningDecision struct {
	SteerBike     bool    `json:"steer_bike"`
	SteeringForce float64 `json:"steering_force"`
}

type Forces struct {
	Pedal float64 `json:"pedal"` // Pedal is a force from 0-1 where 1 is 100% power
	Brake float64 `json:"brake"`

	// If agents do not want to steer, they must set their TurningDecision.SteerBike to false and their steering will not have an impact on the direction of the bike.
	// TurningDecision.SteeringForce is a force from -1 to 1 which maps to -180° to 180°.
	Turning TurningDecision `json:"turning"` // Brake is a force from 0-1 opposing the direction of travel (bike cannot go backwards)
}

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type PhysicalState struct {
	Position     Coordinates `json:"position"`
	Acceleration float64     `json:"acceleration"`
	Velocity     float64     `json:"velocity"`
	Mass         float64     `json:"mass"`
}

type Governance int

const (
	Democracy Governance = iota
	Leadership
	Dictatorship
	Invalid
)
