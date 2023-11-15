package utils

/*
Environment Parameters
*/
const GridHeight float64 = 500.0
const GridWidth float64 = 500.0
const CollisionThreshold float64 = 7.0

/*
Physics Parameters
*/
const MassBike float64 = 50.0
const MassBiker float64 = 70.0
const MassAudi float64 = 1000.0

const BikerMaxForce float64 = 1.0 // The max force a biker can pedal
const AudiMaxForce float64 = 1.0  // The audi's force is equivalent to that of one biker agent going at maximum speed

const DragCoefficient float64 = 1.0 // Drag coefficient can be optimised in experimentation
