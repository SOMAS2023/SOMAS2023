package team_3

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/physics"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"math"
	"sort"

	"github.com/google/uuid"
)

type ISmartAgent interface {
	objects.IBaseBiker
}

type KeyValuePair struct {
	Key   uuid.UUID
	Value float64
}

/*
func SomeFunction(s *server.Server) {
	s.GetLootBoxes()
}
*/

type SmartAgent struct {
	objects.BaseBiker
	currentBike   *objects.MegaBike
	targetLootBox objects.ILootBox
	reputationMap map[uuid.UUID]reputation
	creditMap     map[uuid.UUID]credit

	lootBoxCnt                     float64
	energySpent                    float64
	lastEnergyLevel                float64
	satisfactionOfRecentAllocation float64
}

// DecideAction only pedal
func (agent *SmartAgent) DecideAction() objects.BikerAction {
	if agent.GetEnergyLevel() < agent.lastEnergyLevel {
		agent.energySpent += agent.lastEnergyLevel - agent.GetEnergyLevel()
	} else {
		agent.recalculateSatisfaction()
	}
	agent.lastEnergyLevel = agent.GetEnergyLevel()

	agent.updateRepMap()

	return objects.Pedal
}

// DecideForces considering Hegselmann-Krause model, Ramirez-Cano-Pitt model and Satisfaction
func (agent *SmartAgent) DecideForces(direction uuid.UUID) {
	agentsOnBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()].GetAgents()
	scores := make(map[uuid.UUID]float64)
	totalScore := 0.0
	for _, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		score := rep.isSameColor/ // Cognitive dimension: is same belief?
			+rep.historyContribution + rep.lootBoxGet/ // Pareto principle: give more energy to those with more outcome
			+rep.recentContribution // Forgiveness: forgive agents pedal harder recently
		scores[others.GetID()] = score
		totalScore += score
	}

	for id, score := range scores {
		scores[id] = score / totalScore
	}

	pedalForce := 0.0
	for id, weight := range scores {
		pedalForce += weight * agent.reputationMap[id]._lastPedal
	}
	pedalForce /= agent.satisfactionOfRecentAllocation

	// 因为force是一个struct,包括pedal, brake,和turning，因此需要一起定义，不能够只有pedal
	forces := utils.Forces{
		Pedal: pedalForce,
		Brake: 0.0, // 这里默认刹车为 0
		Turning: utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: physics.ComputeOrientation(agent.GetLocation(), agent.GetGameState().GetMegaBikes()[direction].GetPosition()) - agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()].GetOrientation(),
		}, // 这里默认转向为 0
	}

	agent.SetForces(forces)
}

// DecideJoining accept all
func (agent *SmartAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		decision[agent] = true
	}
	return decision
}

func (agent *SmartAgent) ProposeDirection() utils.Coordinates {
	e := agent.decideTargetLootBox(agent.GetGameState().GetLootBoxes())
	if e != nil {
		panic("unexpected error!")
	}
	return agent.targetLootBox.GetPosition()
}

func (agent *SmartAgent) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	boxesInMap := agent.GetGameState().GetLootBoxes()
	//boxProposed := make([]objects.ILootBox, len(proposals))
	boxProposed := make(map[uuid.UUID]objects.ILootBox, len(proposals))
	for i, pp := range proposals {
		boxProposed[i] = boxesInMap[pp]
	}
	rank := agent.rankTargetProposals(boxProposed)
	return rank
}

func (agent *SmartAgent) DecideAllocation() voting.IdVoteMap {
	agent.lootBoxCnt += 1
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()]
	vote, _ := agent.scoreAgentsForAllocation(currentBike.GetAgents())
	return vote
}

func (agent *SmartAgent) vote_off_leader() float64 {
	vote_off := 0.0
	// vote_off: 1.0 means to vote_off the leader
	id := agent.GetID()
	rep := agent.reputationMap[id]
	if (rep.recentGetEnergy == false) && (rep.isSameColor == 0.0) {
		vote_off = 1.0
	}
	return vote_off
}

func (agent *SmartAgent) whether_need_leader(agentsOnBike []objects.IBaseBiker) *float64 {
	need := 1.0
	// need: 1.0 means it needs the leader
	agent_id := agent.GetID()
	agent_rep := agent.reputationMap[agent_id]

	average_recent_contribution := 0.0
	average_contribution := 0.0
	average_energyRemain := 0.0
	for _, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		average_recent_contribution += rep.recentContribution
		average_contribution += rep.historyContribution
		average_energyRemain += rep.recentContribution
		average_recent_contribution = average_recent_contribution / float64(len(agentsOnBike))
		average_contribution = average_contribution / float64(len(agentsOnBike))
		average_energyRemain = average_energyRemain / float64(len(agentsOnBike))
	}

	if (agent_rep.recentContribution < average_recent_contribution) && (agent_rep.historyContribution < average_contribution) {
		need = 0.0 // selfish personality
	}
	if agent_rep.energyRemain > 2*average_energyRemain {
		need = 0.0 // fear of being taken advantage of
	}
	return &need
}

func (agent *SmartAgent) find_collusion(agentsOnBike []SmartAgent, agentsOnBike2 []objects.IBaseBiker) {
	for i := 0; i < len(agentsOnBike)-1; i++ {
		for j := i + 1; j < len(agentsOnBike); j++ {
			firstAgent := &agentsOnBike[i]
			secondAgent := &agentsOnBike[j]

			firstID := firstAgent.GetID()
			secondID := secondAgent.GetID()

			firstCredit := agent.creditMap[firstID]
			secondCredit := agent.creditMap[secondID]
			firstRep := agent.reputationMap[firstID]
			secondRep := agent.reputationMap[secondID]

			firstNeedPtr := firstAgent.whether_need_leader(agentsOnBike2)
			secondNeedPtr := secondAgent.whether_need_leader(agentsOnBike2)

			first_recent_contribution := firstRep.recentContribution
			second_recent_contribution := secondRep.recentContribution

			average_recent_contribution := 0.0
			for _, others := range agentsOnBike2 {
				id := others.GetID()
				rep := agent.reputationMap[id]
				average_recent_contribution += rep.recentContribution
				average_recent_contribution = average_recent_contribution / float64(len(agentsOnBike2))
			}

			if firstCredit.consecutiveNegativeCount == 3 && secondCredit.consecutiveNegativeCount == 3 && (first_recent_contribution < average_recent_contribution) && (second_recent_contribution < average_recent_contribution) {
				if *firstNeedPtr == 0.0 {
					*firstNeedPtr = 1.0
				}
				if *secondNeedPtr == 0.0 {
					*secondNeedPtr = 1.0
				}
			}
		}
	}
}

func (agent *SmartAgent) vote_leader(agentsOnBike []objects.IBaseBiker, proposedLootBox map[uuid.UUID]objects.ILootBox) (map[uuid.UUID]float64, error) {
	// two-round run-off

	// the first round: top three
	scores := make(map[uuid.UUID]float64)

	for _, others := range agentsOnBike {
		id := others.GetID()
		if id != agent.GetID() {
			rep := agent.reputationMap[id]
			score := rep.historyContribution + rep.lootBoxGet/ // Pareto principle: give more energy to those with more outcome
				+rep.isSameColor/ // Cognitive dimension: is same belief?
				+rep.energyRemain // necessity: must stay alive

			scores[id] = score
		}
	}

	// scores

	var keyValuePairs []KeyValuePair
	for key, value := range scores {
		keyValuePairs = append(keyValuePairs, KeyValuePair{Key: key, Value: value})
	}

	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i].Value > keyValuePairs[j].Value
	})

	indexMap := make(map[uuid.UUID]int)
	for i, pair := range keyValuePairs {
		indexMap[pair.Key] = i
	}

	floatIndexMap := make(map[uuid.UUID]float64)
	for key, index := range indexMap {
		floatIndexMap[key] = float64(index)
	}

	return floatIndexMap, nil

	/*
		sortedIDs := make([]uuid.UUID, 0, len(scores))
		for score := range scores {
			sortedIDs = append(sortedIDs, score)
		}

		sort.Slice(sortedIDs, func(i, j int) bool {
			return scores[sortedIDs[i]] > scores[sortedIDs[j]]
		})

		var topThree []uuid.UUID
		if len(sortedIDs) >= 3 {
			topThree = sortedIDs[:3]
		} else {
			topThree = sortedIDs
		}

		// the second round: borda count

		scores2 := make([]float64, 0)
		for _, uuid := range topThree {
			rep := agent.reputationMap[uuid]
			score := rep.recentContribution // Forgiveness: forgive agents pedal harder recently
			scores2 = append(scores2, score)
		}

		elementCount := make(map[float64]int)
		for _, num := range scores2 {
			elementCount[num]++
		}
		uniqueElements := make([]float64, 0, len(elementCount))
		for num := range elementCount {
			uniqueElements = append(uniqueElements, num)
		}
		sort.Float64s(uniqueElements)
		elementOrder := make(map[float64]int)
		for i, num := range uniqueElements {
			elementOrder[num] = i + 1
		}
		elementOrderList := make([]int, len(scores))
		for i, num := range scores2 {
			elementOrderList[i] = elementOrder[num]
		}

		rank := make(map[uuid.UUID]float64)
		for i, lootBox := range proposedLootBox {
			rank[lootBox.GetID()] = float64(elementOrderList[i])
		}
	*/
}

func (agent *SmartAgent) find_same_colour_highest_loot_lootbox(proposedLootBox map[uuid.UUID]objects.ILootBox) error {
	max_loot := 0.0
	for _, lootbox := range proposedLootBox {
		loot := lootbox.GetTotalResources()
		if loot > max_loot {
			max_loot = loot
			agent.targetLootBox = lootbox
		}
	}
	return nil
}

func (agent *SmartAgent) other_agents_strong(agentsOnBike []objects.IBaseBiker, proposedLootBox map[uuid.UUID]objects.ILootBox) bool {
	// other_agents' energy is higher than the farthest lootbox

	// other_agents' energy
	other_agents_energy := 0.0
	for _, others := range agentsOnBike {
		other_agents_energy += others.GetEnergyLevel()
	}

	//farthest lootbox
	max_distance := 0.0
	for _, lootbox := range proposedLootBox {
		distance := physics.ComputeDistance(lootbox.GetPosition(), agent.GetLocation())
		if distance > float64(max_distance) {
			max_distance = distance
		}
	}
	max_energy := max_distance * 1

	return other_agents_energy > max_energy
}

func (agent *SmartAgent) all_weak(agentsOnBike []objects.IBaseBiker, proposedLootBox map[uuid.UUID]objects.ILootBox) bool {
	total_energy := 0.0

	// total_energy
	for _, others := range agentsOnBike {
		total_energy += others.GetEnergyLevel()
	}

	// nearest same_colour lootbox
	nearest_same_colour_lootbox_distance := math.MaxFloat64
	for _, lootbox := range proposedLootBox {
		if lootbox.GetColour() == agent.GetColour() {
			distance := physics.ComputeDistance(lootbox.GetPosition(), agent.GetLocation())
			if distance < nearest_same_colour_lootbox_distance {
				nearest_same_colour_lootbox_distance = distance
			}
		}
	}
	nearest_same_colour_lootbox_energy := nearest_same_colour_lootbox_distance * 1

	return total_energy < nearest_same_colour_lootbox_energy
}

func (agent *SmartAgent) find_closest_lootbox(proposedLootBox map[uuid.UUID]objects.ILootBox) error {
	min_distance := math.MaxFloat64

	for _, lootbox := range proposedLootBox {
		distance := physics.ComputeDistance(lootbox.GetPosition(), agent.GetLocation())
		// no need to normalize
		if distance < min_distance {
			min_distance = distance
			agent.targetLootBox = lootbox
		}
	}
	return nil
}

func (agent *SmartAgent) decideTargetLootBox(agentsOnBike []objects.IBaseBiker, proposedLootBox map[uuid.UUID]objects.ILootBox) error {
	//dynamic decison of choosing lootbox with the changes in environment
	max_score := 0.0

	// improve all agents' satisfication: while the energy was too low, all agents desire energy
	if agent.all_weak(agentsOnBike, proposedLootBox) == true { //all weak
		agent.find_closest_lootbox(proposedLootBox)
	}

	// free rider - belief that there is a rule, assume other agents would follow the rule
	if agent.other_agents_strong(agentsOnBike, proposedLootBox) == true { //is strong
		agent.find_same_colour_highest_loot_lootbox(proposedLootBox)
	}

	for _, lootbox := range proposedLootBox {
		// agent
		// consider the agent itself's satisfaction
		loot := (lootbox.GetTotalResources() / 8.0) //normalize
		is_color := 0.0
		if lootbox.GetColour() == agent.GetColour() {
			is_color = 1.0
		}
		distance := physics.ComputeDistance(lootbox.GetPosition(), agent.GetLocation())
		normalized_distance := distance / ((utils.GridHeight) * (utils.GridWidth))
		score := 0.2*loot + 0.2*is_color + (-0.3)*normalized_distance

		// environment
		// social capital framework to decide which agents we should cooperate
		// deciding the opinion of ther agents
		same_colour_bikers := make([]objects.IBaseBiker, 0)
		same_colour := 0
		for _, others := range agentsOnBike {
			if others.GetColour() == lootbox.GetColour() {
				same_colour += 1
				same_colour_bikers = append(same_colour_bikers, others)
			}
			score += 0.5 * float64(same_colour/len(agentsOnBike))
		}

		for _, others := range same_colour_bikers {
			id := others.GetID()
			rep := agent.reputationMap[id]
			score += 0.5 * 0.4 * rep.historyContribution / //opinion, trustness from direct experience
				+0.5 * 0.2 * rep.recentContribution / //forgiveness
				+0.5 * 0.4 * rep.energyRemain //improve trustness by decreasing the risk of no efforts
		}

		if score > max_score {
			max_score = score
			agent.targetLootBox = lootbox
		}
	}
	return nil
}

func (agent *SmartAgent) rankTargetProposals(proposedLootBox map[uuid.UUID]objects.ILootBox) map[uuid.UUID]float64 {
	//scores := make([]float64, 0)
	scores := make(map[uuid.UUID]float64)

	for lootbox_agent_id, lootbox := range proposedLootBox {
		other_agents_score := 0.0
		loot := (lootbox.GetTotalResources() / 8.0)
		is_color := 0.0
		if lootbox.GetColour() == agent.GetColour() {
			is_color = 1.0
		}
		rep := agent.reputationMap[lootbox_agent_id]
		other_agents_score = rep.historyContribution + rep.recentContribution + rep.energyRemain
		distance := physics.ComputeDistance(lootbox.GetPosition(), agent.GetLocation())
		normalized_distance := distance / ((utils.GridHeight) * (utils.GridWidth))
		score := 0.2*loot + 0.4*is_color + 0.2*normalized_distance + 0.2*other_agents_score

		scores[lootbox.GetID()] = score
		//scores = append(scores, score)
	}
	// We choose to use the Borda count method to pick a proposal because it can mitigate the Condorcet paradox.
	// Borda count needs to get the rank of all candidates to score Borda points.
	// In this case, according to the Gibbard-Satterthwaite Theorem, Borda count is susceptible to tactical voting.
	// The following steps tend to achieve the rank of lootbox proposals according to their scores calculated. We will return the highest rank to pick the agent with it. (Another Borda score would consider reputation function)这个后面如果可以再考虑如果能得到的话

	// borda count
	var keyValuePairs []KeyValuePair
	for key, value := range scores {
		keyValuePairs = append(keyValuePairs, KeyValuePair{Key: key, Value: value})
	}

	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i].Value > keyValuePairs[j].Value
	})

	indexMap := make(map[uuid.UUID]int)
	for i, pair := range keyValuePairs {
		indexMap[pair.Key] = i
	}

	floatIndexMap := make(map[uuid.UUID]float64)
	for key, index := range indexMap {
		floatIndexMap[key] = float64(index)
	}

	return floatIndexMap
}

// scoreAgentsForAllocation if self energy level is low (below average cost for a lootBox), we follow 'Smallest First', else 'Ration'
func (agent *SmartAgent) scoreAgentsForAllocation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) {
	scores := make(map[uuid.UUID]float64)
	totalScore := 0.0
	if agent.energySpent/agent.lootBoxCnt > agent.GetEnergyLevel() {
		// go 'Smallest First' strategy, only take energyRemain into consideration
		for _, others := range agentsOnBike {
			id := others.GetID()
			score := agent.reputationMap[id].energyRemain
			scores[others.GetID()] = score
			totalScore += score
		}
	} else {
		// go 'Ration' strategy, considering all facts
		for _, others := range agentsOnBike {
			id := others.GetID()
			rep := agent.reputationMap[id]
			score := rep.isSameColor/ // Cognitive dimension: is same belief?
				+rep.historyContribution + rep.lootBoxGet/ // Pareto principle: give more energy to those with more outcome
				+rep.recentContribution/ // Forgiveness: forgive agents pedal harder recently
				-rep.energyGain/ // Equality: Agents received more energy before should get less this time
				+rep.energyRemain // Need: Agents with lower energyLevel require more, try to meet their need
			scores[others.GetID()] = score
			totalScore += score
		}
	}

	// normalize scores
	for id, score := range scores {
		scores[id] = score / totalScore
	}

	return scores, nil
}

func (agent *SmartAgent) updateRepMap() {
	if agent.reputationMap == nil {
		agent.reputationMap = make(map[uuid.UUID]reputation)
	}
	for _, bikes := range agent.GetGameState().GetMegaBikes() {
		for _, otherAgent := range bikes.GetAgents() {
			rep, exist := agent.reputationMap[otherAgent.GetID()]
			if !exist {
				rep = reputation{}
			}
			rep.updateScore(otherAgent, agent.GetColour())
			agent.reputationMap[otherAgent.GetID()] = rep
		}
	}
}

func (agent *SmartAgent) recalculateSatisfaction() {
	agentsOnBike := agent.GetGameState().GetMegaBikes()[agent.GetMegaBikeId()].GetAgents()
	scores := make([]float64, len(agentsOnBike))
	gains := make([]float64, len(agentsOnBike))
	for idx, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		score := rep.isSameColor/ // Cognitive dimension: is same belief?
			+rep.historyContribution + rep.lootBoxGet/ // Pareto principle: give more energy to those with more outcome
			+rep.recentContribution/ // Forgiveness: forgive agents pedal harder recently
			-rep.energyGain/ // Equality: Agents received more energy before should get less this time
			+rep.energyRemain // Need: Agents with lower energyLevel require more, try to meet their need
		scores[idx] = score
		gains[idx] = agent.reputationMap[id]._recentEnergyGain
	}
	sort.Slice(gains, func(i, j int) bool {
		return scores[i] < scores[j]
	})
	agent.satisfactionOfRecentAllocation = measureOrder(gains)
}

// To measure how an array is well sorted, result normalized to 0~1
func measureOrder(input []float64) float64 {
	inversionCnt := 0.0
	size := len(input)
	for i, n := range input {
		j := i + 1
		for j < size {
			if n > input[j] { // 升序为正序
				inversionCnt += 1
			}
			j += 1
		}
	}
	return 1.0 - 2.0*inversionCnt/float64(size*(size-1))
}
