package environmentparameters

// Grid size
var gridLength float64 = 500.0
var gridWidth float64 = 500.0

// GridLength returns the length of the grid.
func GridLength() float64 {
	return gridLength
}

// GridWidth returns the width of the grid.
func GridWidth() float64 {
	return gridWidth
}
