package utils

/*
Environment Parameters
*/
const GridHeight float64 = 75.0
const GridWidth float64 = 75.0
const CollisionThreshold float64 = 7.0
const Epsilon float64 = 0.01 // tolerance for FP rounding and checking if == 1.0
const BikersOnBike = 8
const ReplenishEnergyEveryRound = true
const ResetPointsEveryRound = true
const RespawnEveryRound = true
const RoundIterations = 100

/*
Server Parameters
*/
const ReplenishLootBoxes bool = true
const ReplenishMegaBikes bool = true

/*
Physics Parameters
*/
const MassBike float64 = 1.0
const MassBiker float64 = 1.0
const MassAudi float64 = 10.0

const BikerMaxForce float64 = 1.0 // The max force a biker can pedal
const AudiMaxForce float64 = 1.0  // The audi's force is equivalent to that of one biker agent going at maximum speed

const DragCoefficient float64 = 0.5 // Drag coefficient can be optimised in experimentation

const MovingDepletion float64 = 0.01 // proportionality of energy loss

const LimboEnergyPenalty float64 = -0.25 // amount of energy lost per round when off a bike

const DeliberativeDemocracyPenalty float64 = 0.05 // amount of energy lost per vote in a deliberative democracy
const LeadershipDemocracyPenalty float64 = 0.025  // amount of energy lost per vote in a leadership democracy

/*
Resources - Points and Energy
*/
const PointsFromSameColouredLootBox = 5.0

/*
Audi Behavior
*/
const AudiTargetsEmptyMegaBike bool = false
const AudiOnlyTargetsStationaryMegaBike bool = true // if false, targeting slowest
const AudiRemovesMegaBike bool = false

/*
Voting Method Choice
*/
type voteMethods int

const (
	PLURALITY voteMethods = iota
	RUNOFF
	BORDACOUNT
	INSTANTRUNOFF
	APPROVAL
	COPELANDSCORING
)

const VoteAction voteMethods = PLURALITY
