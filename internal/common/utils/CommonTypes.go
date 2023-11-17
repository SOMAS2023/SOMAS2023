package utils

import "github.com/google/uuid"

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

type Forces struct {
	Pedal   float64 // Pedal is a force from 0-1 where 1 is 100% power
	Brake   float64 // Brake is a force from 0-1 opposing the direction of travel (bike cannot go backwards)
	Turning float64 // Turning is a force from -1 to 1 which maps to -180° to 180°
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

type INormaliseVoteMap interface {
	IsNormalisedVoteMap()
}

type PositionVoteMap map[Coordinates]float64

func (PositionVoteMap) IsNormalisedVoteMap() {}

type IdVoteMap map[uuid.UUID]float64

func (IdVoteMap) isNormalisedVoteMap() {}
