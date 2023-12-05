package team1

import (
	obj "SOMAS2023/internal/common/objects"
	physics "SOMAS2023/internal/common/physics"
	utils "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

// -------------------SETTERS AND GETTERS-----------------------------
// Returns a list of bikers on the same bike as the agent
func (bb *Biker1) GetFellowBikers() []obj.IBaseBiker {
	gs := bb.GetGameState()
	bikeId := bb.GetBike()
	return gs.GetMegaBikes()[bikeId].GetAgents()
}

func (bb *Biker1) GetBikeInstance() obj.IMegaBike {
	gs := bb.GetGameState()
	bikeId := bb.GetBike()
	return gs.GetMegaBikes()[bikeId]
}

func (bb *Biker1) GetLootLocation(id uuid.UUID) utils.Coordinates {
	gs := bb.GetGameState()
	lootboxes := gs.GetLootBoxes()
	lootbox := lootboxes[id]
	return lootbox.GetPosition()
}

func (bb *Biker1) GetAverageOpinionOfAgent(biker uuid.UUID) float64 {
	fellowBikers := bb.GetFellowBikers()
	opinionSum := 0.0
	for _, agent := range fellowBikers {
		opinionSum += agent.QueryReputation(biker)
	}
	return opinionSum / float64(len(fellowBikers))
}

// -------------------END OF SETTERS AND GETTERS----------------------

// -------------------END OF CHANGE BIKE FUNCTIONS----------------------

// Find an agent from their id
func (bb *Biker1) GetAgentFromId(agentId uuid.UUID) obj.IBaseBiker {
	agents := bb.GetAllAgents()
	for _, agent := range agents {
		if agent.GetID() == agentId {
			return agent
		}
	}
	return nil
}

// Get all agents in the game
func (bb *Biker1) GetAllAgents() []obj.IBaseBiker {
	gs := bb.GetGameState()
	// get all agents
	agents := make([]obj.IBaseBiker, 0)
	for _, agent := range gs.GetAgents() {
		agents = append(agents, agent)
	}
	return agents
}

// -------------------END OF CHANGE BIKE FUNCTIONS----------------------

// -------------------SELFISHNESS FUNCTIONS----------------------

// // Success-Relationship algo for calculating selfishness score
func calculateSelfishnessScore(success float64, relationship float64) float64 {
	difference := math.Abs(success - relationship)
	var overallScore float64
	if success >= relationship {
		overallScore = 0.5 + ((difference) / 2)
	} else if relationship > success {
		overallScore = 0.5 - ((difference) / 2)
	}
	return overallScore
}

func (bb *Biker1) GetSelfishness(agent obj.IBaseBiker) float64 {
	pointSum := bb.GetPoints() + agent.GetPoints()
	var relativeSuccess float64
	if pointSum == 0 {
		relativeSuccess = 0.5
	} else {
		relativeSuccess = float64((agent.GetPoints() - bb.GetPoints()) / (pointSum)) //-1 to 1
		relativeSuccess = (relativeSuccess + 1.0) / 2.0                              //shift to 0 to 1
	}
	id := agent.GetID()
	ourRelationship := bb.opinions[id].opinion
	return calculateSelfishnessScore(relativeSuccess, ourRelationship)
}

// -------------------END OF SELFISHNESS FUNCTIONS----------------------
// -------------------BIKE CHANGE HELPER FUNCTIONS----------------------
func (bb *Biker1) GetNearBikeObjects(bike obj.IMegaBike) (int64, int64, int64) {
	_, reachableDistance := bb.energyToReachableDistance(bb.GetEnergyLevel(), bike)
	lootBoxCount := 0
	lootBoxOurColor := 0
	bikeCount := 0
	for _, lootbox := range bb.GetGameState().GetLootBoxes() {
		distance := physics.ComputeDistance(lootbox.GetPosition(), bike.GetPosition())
		if distance <= reachableDistance {
			lootBoxCount += 1
			if lootbox.GetColour() == bb.GetColour() {
				lootBoxOurColor += 1
			}
		}
	}
	for _, nearbyBike := range bb.GetGameState().GetMegaBikes() {
		distance := physics.ComputeDistance(nearbyBike.GetPosition(), bike.GetPosition())
		if distance <= reachableDistance {
			bikeCount += 1
		}
	}

	return int64(lootBoxCount), int64(lootBoxOurColor), int64(bikeCount)
}

// -------------------END OF BIKE CHANGE HELPER FUNCTIONS---------------
