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

type Forces struct {
	Pedal   float64
	Brake   float64
	Turning float64
}

type Coordinates struct {
	X float64
	Y float64
}
