package team4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"math"
	"sort"

	"github.com/google/uuid"
)

func (agent *BaselineAgent) rankTargetProposals(proposedLootBox []objects.ILootBox) (voting.LootboxVoteMap, error) {
	rank := make(voting.LootboxVoteMap) //make(map[uuid.UUID]float64)
	ranksum := make(map[uuid.UUID]float64)
	totalsum := float64(0)
	totaloptions := len(proposedLootBox)
	audiPos := agent.GetGameState().GetAudi().GetPosition()

	fellowBikers := agent.GetFellowBikers()
	//This is the relavtive reputation and honest for bikers my bike
	reputationRank, e1 := agent.rankFellowsReputation(fellowBikers)
	honestyRank, e2 := agent.rankFellowsHonesty(fellowBikers)
	//This is the absolute reputation and honest for bikers my bike
	// reputationRank  := agent.reputation
	// honestyRank  := agent.honestyMatrix
	if e1 != nil || e2 != nil {
		panic("unexpected error!")
	}
	//sort proposed loot boxes by distance from agent
	sort.Slice(proposedLootBox, func(i, j int) bool {
		return physics.ComputeDistance(agent.GetLocation(), proposedLootBox[i].GetPosition()) < physics.ComputeDistance(agent.GetLocation(), proposedLootBox[j].GetPosition())
	})

	for i, lootBox := range proposedLootBox {
		lootboxID := lootBox.GetID()
		lootboxResources := lootBox.GetTotalResources()
		distanceFromAudi := physics.ComputeDistance(audiPos, lootBox.GetPosition())
		//if energy level is below threshold, increase weighting towards distance
		distanceRank := float64(totaloptions - i)
		if agent.GetEnergyLevel() < minEnergyThreshold {
			distanceRank *= 2.0
		}

		//loop through all fellow bikers and check if they have the same colour as the lootbox
		for _, fellow := range fellowBikers {
			fellowID := fellow.GetID()
			if fellow.GetColour() == lootBox.GetColour() {
				weight := (distanceWeight * distanceRank) + (reputationWeight * reputationRank[fellowID]) + (honestyWeight * honestyRank[fellowID]) + (audiDistanceWeight * distanceFromAudi) + (resourceWeight * lootboxResources)
				ranksum[lootboxID] += weight
				totalsum += weight
			}
		}

		if lootBox.GetColour() == agent.GetColour() {
			weight := (distanceRank * distanceWeight * 1.25) + (audiDistanceWeight * distanceFromAudi) + (resourceWeight * lootboxResources)
			ranksum[lootboxID] += weight
			totalsum += weight
		}
		if ranksum[lootboxID] == 0 {
			weight := (distanceRank * distanceWeight * 2.6) + (audiDistanceWeight * distanceFromAudi) + (resourceWeight * lootboxResources)
			ranksum[lootboxID] = weight
			totalsum += weight
		}
	}
	for _, lootBox := range proposedLootBox {
		rank[lootBox.GetID()] = ranksum[lootBox.GetID()] / totalsum
	}

	return rank, nil
}

func (agent *BaselineAgent) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	//fmt.Println("Final Direction Vote")
	agent.UpdateDecisionData()
	//We need to fix this ASAP
	boxesInMap := agent.GetGameState().GetLootBoxes()
	boxProposed := make([]objects.ILootBox, len(proposals))
	count := 0
	for _, i := range proposals {
		boxProposed[count] = boxesInMap[i]
		count++
	}

	rank, e := agent.rankTargetProposals(boxProposed)
	if e != nil {
		panic("unexpected error!")
	}
	return rank
}

func (agent *BaselineAgent) DecideForce(direction uuid.UUID) {
	agent.targetLoot = direction
	currLocation := agent.GetLocation()
	currentLootBoxes := agent.GetGameState().GetLootBoxes()
	audiPos := agent.GetGameState().GetAudi().GetPosition()

	agent.lootBoxColour = currentLootBoxes[direction].GetColour()
	agent.lootBoxLocation = currentLootBoxes[direction].GetPosition()

	distanceFromAudi := physics.ComputeDistance(currLocation, audiPos)
	pedalForce := 1.0

	if distanceFromAudi < audiDistanceThreshold {
		deltaX := audiPos.X - currLocation.X
		deltaY := audiPos.Y - currLocation.Y
		// Steer in opposite direction to audi
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi
		// Steer in opposite direction to audi
		var flipAngle float64
		if normalisedAngle < 0.0 {
			flipAngle = normalisedAngle + 1.0
		} else if normalisedAngle > 0.0 {
			flipAngle = normalisedAngle - 1.0
		}
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: flipAngle - agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetOrientation(),
		}
		escapeAudiForces := utils.Forces{
			Pedal:   utils.BikerMaxForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		agent.SetForces(escapeAudiForces)
	} else {
		targetPos := currentLootBoxes[agent.targetLoot].GetPosition()
		deltaX := targetPos.X - currLocation.X
		deltaY := targetPos.Y - currLocation.Y
		angle := math.Atan2(deltaY, deltaX)
		normalisedAngle := angle / math.Pi

		// Default BaseBiker will always
		turningDecision := utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: normalisedAngle - agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetOrientation(),
		}
		if agent.GetEnergyLevel() <= 0.5 {
			pedalForce = pedalForce * agent.GetEnergyLevel()
		}
		nearestBoxForces := utils.Forces{
			Pedal:   pedalForce,
			Brake:   0.0,
			Turning: turningDecision,
		}
		agent.SetForces(nearestBoxForces)
	}
}

func (agent *BaselineAgent) nearestLoot() uuid.UUID {
	currLocation := agent.GetLocation()
	shortestDist := math.MaxFloat64
	var nearestBox uuid.UUID
	var currDist float64
	for _, loot := range agent.GetGameState().GetLootBoxes() {
		x, y := loot.GetPosition().X, loot.GetPosition().Y
		currDist = math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
		if currDist < shortestDist {
			nearestBox = loot.GetID()
			shortestDist = currDist
		}
	}
	return nearestBox
}

func (agent *BaselineAgent) ProposeDirection() uuid.UUID {
	if agent.GetEnergyLevel() < minEnergyThreshold+0.05 {
		return agent.nearestLoot()
	}
	//fmt.Println("Propose Direction")
	agent.UpdateDecisionData()

	var lootBoxesWithinThreshold []objects.ILootBox
	audiPos := agent.GetGameState().GetAudi().GetPosition()
	agentLocation := agent.GetLocation() // agent's location

	for _, lootbox := range agent.GetGameState().GetLootBoxes() {
		if physics.ComputeDistance(lootbox.GetPosition(), audiPos) > audiDistanceThreshold {
			lootBoxesWithinThreshold = append(lootBoxesWithinThreshold, lootbox)
		}
	}

	// Sort the lootboxes within threshold by distance from the agent
	sort.Slice(lootBoxesWithinThreshold, func(i, j int) bool {
		return physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[i].GetPosition()) <
			physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[j].GetPosition())
	})

	// Select the closest lootbox if any are within the threshold
	if len(lootBoxesWithinThreshold) > 0 {
		closestLootBox := lootBoxesWithinThreshold[0]
		return closestLootBox.GetID()
	} else {
		return agent.nearestLoot()
	}
}

/////////////////////////////////// DICATOR FUNCTIONS /////////////////////////////////////

func (agent *BaselineAgent) DictateDirection() uuid.UUID {
	//fmt.Println("Dictate Direction")
	agent.UpdateDecisionData()

	var lootBoxesWithinThreshold []objects.ILootBox
	audiPos := agent.GetGameState().GetAudi().GetPosition()
	agentLocation := agent.GetLocation() // agent's location

	for _, lootbox := range agent.GetGameState().GetLootBoxes() {
		if physics.ComputeDistance(lootbox.GetPosition(), audiPos) > audiDistanceThreshold {
			lootBoxesWithinThreshold = append(lootBoxesWithinThreshold, lootbox)
		}
	}

	// Sort the lootboxes within threshold by distance from the agent
	sort.Slice(lootBoxesWithinThreshold, func(i, j int) bool {
		return physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[i].GetPosition()) <
			physics.ComputeDistance(agentLocation, lootBoxesWithinThreshold[j].GetPosition())
	})

	// Select the closest lootbox if any are within the threshold
	if len(lootBoxesWithinThreshold) > 0 {
		closestLootBox := lootBoxesWithinThreshold[0]
		return closestLootBox.GetID()
	} else {
		return agent.nearestLoot()
	}
}
