package utils

/*
Environment Parameters
*/
const GridHeight float64 = 75.0
const GridWidth float64 = 75.0
const CollisionThreshold float64 = 7.0
const Epsilon float64 = 0.01 // tolerance for FP rounding and checking if == 1.0

/*
Physics Parameters
*/
const MassBike float64 = 50.0
const MassBiker float64 = 70.0
const MassAudi float64 = 1000.0

const BikerMaxForce float64 = 1.0 // The max force a biker can pedal
const AudiMaxForce float64 = 1.0  // The audi's force is equivalent to that of one biker agent going at maximum speed

const DragCoefficient float64 = 1.0 // Drag coefficient can be optimised in experimentation

const MovingDepletion float64 = 0.1 // proportionality of energy loss

const LimboEnergyPenalty float64 = -0.25 // amount of energy lost per round when off a bike

/*
Resources - Points and Energy
*/
const PointsFromSameColouredLootBox = 5.0
