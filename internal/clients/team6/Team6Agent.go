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
	distAudiThreshold     = 75
	energyChangeThreshold = 0.8 // energy level which allows us to change bike
)

type Team6Biker struct {
	*objects.BaseBiker
	Changeflag bool
	Trust      map[uuid.UUID]Trust
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

	currColour := bb.GetColour()
	nearestLootColour := bb.GetGameState().GetLootBoxes()[bb.nearestLoot()].GetColour()

	var t bool = true
	if bb.nearestSameColourLoot() == uuid.Nil {
		t = false
	}

	if t == true {
		nearestLootColour = bb.GetGameState().GetLootBoxes()[bb.nearestSameColourLoot()].GetColour()
	}

	if bb.GetEnergyLevel() > energyChangeThreshold {
		bb.Changeflag = true
	}

	if (bb.GetEnergyLevel() < energyThreshold) && (bb.Changeflag == true) {
		bb.Changeflag = false
		return objects.ChangeBike
	}

	if nearestLootColour == currColour {
		// keep pedaling if current colour = goal

		return objects.Pedal
	} else {

		return objects.ChangeBike

	}

	//return objects.ChangeBike
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

	var nearestBike1 uuid.UUID
	var currDist float64
	shortest := math.MaxFloat64

	for _, megabike := range bb.GetGameState().GetMegaBikes() {

		currDist = bb.bikeToNearestLoot(megabike)

		//currLocation := megabike.GetPosition()

		if currDist < shortest {

			nearestBike1 = megabike.GetID()
			shortest = currDist
		}
	}
	return nearestBike1
}

// Return the nearest BikeID to the goal color
func (bb *Team6Biker) GetNearSameColourBikeID() uuid.UUID {

	var nearestBike uuid.UUID
	var currDist float64
	shortest := math.MaxFloat64

	sameColourLootList := []objects.ILootBox{}
	for _, loot := range bb.GetGameState().GetLootBoxes() {
		if loot.GetColour() == bb.GetColour() {
			sameColourLootList = append(sameColourLootList, loot)
		}
	}

	if len(sameColourLootList) == 0 {
		return bb.GetBike()
	}

	for _, megabike := range bb.GetGameState().GetMegaBikes() {
		currDist = bb.bikeToSameColorNearestLoot(megabike, sameColourLootList)

		//currLocation := megabike.GetPosition()

		if currDist < shortest {

			nearestBike = megabike.GetID()
			shortest = currDist
		}
	}
	return nearestBike
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
	bikeID := bb.GetNearSameColourBikeID()
	//bikeID := bb.GetNearLootBoxBikeID()
	if bb.GetEnergyLevel() < energyThreshold {
		bikeID = bb.GetNearLootBoxBikeID()
	}
	//fmt.Print(bikeID)
	return bikeID

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
		BaseBiker:  bb,
		Changeflag: true,
		Trust:      make(map[uuid.UUID]Trust),
	}
}
