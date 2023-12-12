package team6

import (
	"SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type BikerAction int

const (
	Pedal BikerAction = iota
	ChangeBike
	energyThreshold       = 0.35 // energy level which mind us we need to change bike
	reputationThreshold   = 0.5
	distAudiThreshold     = 10
	energyChangeThreshold = 0.9 // energy level which allows us to change bike

)

type Team6Biker struct {
	*objects.BaseBiker
	Changeflag        bool
	Trust             map[uuid.UUID]Trust
	Goal              bool
	MegabikeTrustList map[uuid.UUID]float64
	CheckPointsget    int
	Getlootboxterm    int
	itercount         int
}

type ITeam6Biker interface {
	objects.IBaseBiker
	DecideAction() objects.BikerAction
	DecideForce(direction uuid.UUID)
	DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool
	ChangeBike() uuid.UUID
	DecideAllocation() voting.IdVoteMap
	nearestLoot() uuid.UUID
	VoteForKickout() map[uuid.UUID]int
	nearestSameColourLoot() uuid.UUID
	GetNearSameColourBikeID() uuid.UUID
	GetNearLootBoxBikeID() uuid.UUID
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause thedistribution to be equal across all of them
func (bb *Team6Biker) DecideAllocation() voting.IdVoteMap {
	bikeID := bb.GetBike()
	fellowBikers := bb.GetGameState().GetMegaBikes()[bikeID].GetAgents()
	distribution := make(voting.IdVoteMap)
	var totalreputaion float64 = 0
	for _, agent := range fellowBikers {
		if agent.GetID() != bb.GetID() {
			totalreputaion += bb.QueryReputation(agent.GetID())
		} else {
			totalreputaion += 1.0
		}
	}

	for _, agent := range fellowBikers {
		if agent.GetID() == bb.GetID() {
			distribution[agent.GetID()] = 1 / totalreputaion
		} else {
			distribution[agent.GetID()] = bb.QueryReputation(agent.GetID()) / totalreputaion
		}
	}

	return distribution
}

// in the MVP the biker's action defaults to pedaling (as it won't be able to change bikes)
// in future implementations this function will be overridden by the agent's specific strategy
// which will be used to determine whether to pedalor try to change bike
func (bb *Team6Biker) DecideAction() objects.BikerAction {
	bb.itercount += 1
	if _, exists := bb.MegabikeTrustList[bb.GetBike()]; exists {
		if bb.GetEnergyLevel() >= 0.995 {
			//bb.Getlootboxterm = bb.Getlootboxterm + 1
			bb.MegabikeTrustList[bb.GetBike()] -= 0.1
		}
	} else {
		bb.MegabikeTrustList[bb.GetBike()] = 1.0
		//bb.Getlootboxterm = 0
		bb.CheckPointsget = bb.GetPoints()
		fmt.Println(bb.CheckPointsget)
	}

	if bb.GetPoints() > bb.CheckPointsget {
		bb.CheckPointsget = bb.GetPoints()
		bb.MegabikeTrustList[bb.GetBike()] += 0.7
	}
	if bb.MegabikeTrustList[bb.GetBike()] < 0.21 {
		fmt.Println(bb.CheckPointsget)
		fmt.Println(bb.GetPoints())
		fmt.Println(bb.MegabikeTrustList[bb.GetBike()])
		fmt.Println(bb.itercount)
		bb.Changeflag = true
	}
	//return objects.Pedal

	if bb.ChangeBike() == bb.GetBike() {
		bb.Changeflag = false
	}
	if bb.Changeflag {
		return objects.ChangeBike
	} else {
		return objects.Pedal
	}
}

func GetMostCommonColor(agents []objects.IBaseBiker) (utils.Colour, int, int) {

	//fmt.Println("Start")
	colorCounts := make(map[utils.Colour]int)

	for _, ags := range agents {
		color := ags.GetColour()
		colorCounts[color]++
		fmt.Println(color, colorCounts[color])
	}

	var mostCommonColor utils.Colour
	maxCount := 0
	mostCommonColorIndex := -1

	for i, ags := range agents {
		color := ags.GetColour()
		count := colorCounts[color]
		if count > maxCount {
			mostCommonColor = color
			maxCount = count
			mostCommonColorIndex = i
		}
	}

	return mostCommonColor, maxCount, mostCommonColorIndex
}

func (bb *Team6Biker) bikeToNearestLoot(Bike objects.IMegaBike) float64 {
	currLocation := Bike.GetPosition()
	shortestDist := math.MaxFloat64

	var currDist float64
	for _, loot := range bb.GetGameState().GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			//nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return shortestDist
}

func (bb *Team6Biker) bikeToSameColorNearestLoot(bike objects.IMegaBike, lootBoxes []objects.ILootBox) float64 {
	currLocation := bike.GetPosition()
	shortestDist := math.MaxFloat64

	var currDist float64
	for _, loot := range lootBoxes {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			//nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return shortestDist
}

// Return the nearest BikeID to the goal color
func (bb *Team6Biker) GetNearLootBoxBikeID() uuid.UUID {

	var nearestBike1 = uuid.Nil
	var currDist float64
	shortest := math.MaxFloat64

	for _, megabike := range bb.GetGameState().GetMegaBikes() {

		currDist = bb.bikeToNearestLoot(megabike)

		//currLocation := megabike.GetPosition()
		_, exists := bb.MegabikeTrustList[megabike.GetID()]
		if currDist < shortest && (bb.MegabikeTrustList[megabike.GetID()] > 0.21 || !exists) && bb.CheckBikeFull(megabike) > 0 {
			nearestBike1 = megabike.GetID()
			shortest = currDist
		}
	}
	if nearestBike1 != uuid.Nil {
		return nearestBike1
	} else {
		return bb.GetBike()
	}

}

// Return the nearest BikeID to the goal color
func (bb *Team6Biker) GetNearSameColourBikeID() uuid.UUID {

	var nearestBike = uuid.Nil
	var currDist float64
	shortest := math.MaxFloat64

	sameColourLootList := []objects.ILootBox{}
	for _, loot := range bb.GetGameState().GetLootBoxes() {
		if loot.GetColour() == bb.GetColour() {
			sameColourLootList = append(sameColourLootList, loot)
		}
	}

	if len(sameColourLootList) == 0 {

		return bb.GetNearLootBoxBikeID()
	}

	for _, megabike := range bb.GetGameState().GetMegaBikes() {
		currDist = bb.bikeToSameColorNearestLoot(megabike, sameColourLootList)

		//currLocation := megabike.GetPosition()
		_, exists := bb.MegabikeTrustList[megabike.GetID()]
		if currDist < shortest && (bb.MegabikeTrustList[megabike.GetID()] > 0.21 || !exists) && bb.CheckBikeFull(megabike) > 0 {
			nearestBike = megabike.GetID()
			shortest = currDist
		}
	}
	if nearestBike != uuid.Nil {
		return nearestBike
	} else {
		return bb.GetBike()
	}
}

func (bb *Team6Biker) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {

		decision[agent] = false
		currentBiker := bb.FindBiker(agent)

		if agent != bb.GetID() {
			if bb.QueryReputation(agent) > reputationThreshold || currentBiker.GetColour() == bb.GetColour() {
				decision[agent] = true
			}
		}
	}
	return decision
}

// decide which bike to go to.
func (bb *Team6Biker) ChangeBike() uuid.UUID {
	bikeID := bb.GetNearLootBoxBikeID()
	if !bb.Goal {
		return bikeID
	} else {
		return bb.GetNearSameColourBikeID()
	}
}

func (bb *Team6Biker) FindBiker(agentID uuid.UUID) objects.IBaseBiker {
	allAgents := bb.GetGameState().GetAgents()
	for _, agent := range allAgents {
		if agentID == agent.GetID() {
			return agent
		}
	}
	fmt.Print("Do not find such agent")
	return bb.BaseBiker
}

func InitialiseBiker6(bb *objects.BaseBiker) objects.IBaseBiker {
	fmt.Printf("Generating Biker for Team 6")
	bb.GroupID = 6
	//bb.soughtColour = utils.Red
	return &Team6Biker{
		BaseBiker:         bb,
		Changeflag:        false,
		Trust:             make(map[uuid.UUID]Trust),
		Goal:              true, // decide same color or nearest loot, true = nearest color, false = nearest loot
		MegabikeTrustList: make(map[uuid.UUID]float64),
		CheckPointsget:    0,
		Getlootboxterm:    0,
		itercount:         0,
	}
}

func (bb *Team6Biker) CheckBikeFull(bike objects.IMegaBike) int {
	return 9 - len(bike.GetAgents())
}
