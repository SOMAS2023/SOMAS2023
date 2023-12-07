package team1

import (
	obj "SOMAS2023/internal/common/objects"
	utils "SOMAS2023/internal/common/utils"
	"math"

	"github.com/google/uuid"
)

// -------------------MATHS HELPER FUNCTIONS ------------------------
func (bb *Biker1) ComputeDistance(a utils.Coordinates, b utils.Coordinates) float64 {
	return math.Sqrt((math.Pow((a.X-b.X), 2) + math.Pow((a.Y-b.Y), 2)))
}

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

func (bb *Biker1) GetAverageOpinionOfBike(megabike obj.IMegaBike) float64 {
	bikers := megabike.GetAgents()
	totalBikers := len(bikers)
	if totalBikers == 0 {
		return 0.5
	}
	sumOpinions := 0.0
	for _, biker := range bikers {
		if biker.GetID() == bb.GetID() {
			continue
		}
		op, ok := bb.opinions[biker.GetID()]
		if ok {
			sumOpinions += op.opinion
		} else {
			newOpinion := Opinion{
				effort:   0.5,
				trust:    0.5,
				fairness: 0.5,
				opinion:  0.5,
			}
			bb.opinions[biker.GetID()] = newOpinion
			sumOpinions += 0.5
		}
	}

	return sumOpinions / float64(totalBikers)
}

// -------------------END OF SETTERS AND GETTERS----------------------

func (bb *Biker1) DistanceFromAudi(obj.IMegaBike) float64 {
	return bb.ComputeDistance(bb.GetLocation(), bb.GetGameState().GetAudi().GetPosition())
}

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

	agentMap := gs.GetAgents()
	agents := make([]obj.IBaseBiker, 0, len(agentMap))
	for _, agent := range agentMap {
		agents = append(agents, agent)
	}
	return agents
}

// -------------------SELFISHNESS/HELPFUL FUNCTIONS----------------------

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
	maxPoints := 0
	for _, agents := range bb.GetFellowBikers() {
		if agents.GetPoints() > maxPoints {
			maxPoints = agents.GetPoints()
		}
	}
	var relativeSuccess float64
	if maxPoints == 0 {
		relativeSuccess = 0.5
	} else {
		relativeSuccess = float64((agent.GetPoints() - bb.GetPoints()) / (maxPoints)) //-1 to 1
		relativeSuccess = (relativeSuccess + 1.0) / 2.0                               //shift to 0 to 1
	}
	id := agent.GetID()
	ourRelationship := bb.opinions[id].opinion
	return calculateSelfishnessScore(relativeSuccess, ourRelationship)
}

func (bb *Biker1) getHelpfulAllocation() map[uuid.UUID]float64 {
	fellowBikers := bb.GetFellowBikers()

	sumEnergyNeeds := 0.0
	helpfulAllocation := make(map[uuid.UUID]float64)
	for _, agent := range fellowBikers {
		// energyNeed := 1.0 - agent.GetEnergyLevel()
		energyNeed := 1.0 - bb.prevEnergy[agent.GetID()]

		helpfulAllocation[agent.GetID()] = energyNeed
		sumEnergyNeeds = sumEnergyNeeds + energyNeed
		if sumEnergyNeeds != 0 {
			for agentId := range helpfulAllocation {
				helpfulAllocation[agentId] /= sumEnergyNeeds
			}
		}

	}
	return helpfulAllocation
}

// -------------------END OF SELFISHNESS/HELPFUL FUNCTIONS----------------------
// -------------------BIKE CHANGE HELPER FUNCTIONS----------------------
func (bb *Biker1) GetNearBikeObjects(bike obj.IMegaBike) (int64, int64, int64) {
	_, reachableDistance := bb.energyToReachableDistance(bb.GetEnergyLevel(), bike)
	lootBoxCount := 0
	lootBoxOurColor := 0
	bikeCount := 0
	for _, lootbox := range bb.GetGameState().GetLootBoxes() {
		distance := bb.ComputeDistance(lootbox.GetPosition(), bike.GetPosition())
		//fmt.Printf("distance from bike %v to lootox %v is %v\n", bike.GetID(), lootbox.GetID(), distance)
		if distance <= reachableDistance {
			lootBoxCount += 1
			if lootbox.GetColour() == bb.GetColour() {
				lootBoxOurColor += 1
			}
		}
	}
	for _, nearbyBike := range bb.GetGameState().GetMegaBikes() {
		if nearbyBike.GetID() == bike.GetID() {
			continue
		}
		distance := bb.ComputeDistance(nearbyBike.GetPosition(), bike.GetPosition())
		if distance <= 20.0 {
			bikeCount += 1
		}
	}

	return int64(lootBoxCount), int64(lootBoxOurColor), int64(bikeCount)
}

// -------------------END OF BIKE CHANGE HELPER FUNCTIONS---------------
