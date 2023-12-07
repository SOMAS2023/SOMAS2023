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

type SmartAgent struct {
	objects.BaseBiker
	targetLootBox objects.ILootBox
	reputationMap map[uuid.UUID]reputation

	lootBoxCnt                     float64
	energySpent                    float64
	lastEnergyLevel                float64
	lastEnergyCost                 float64
	satisfactionOfRecentAllocation float64
	badTeam                        bool
	lastPedal                      float64
}

func (agent *SmartAgent) DecideGovernance() utils.Governance {
	currentBike, e := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	if e == false { // not on bike at initial state
		return utils.Democracy
	}
	governance := agent.which_governance_method(currentBike.GetAgents())
	return governance
}

// DecideAction change bike if find badLeader or badTeam
func (agent *SmartAgent) DecideAction() objects.BikerAction {
	if agent.GetEnergyLevel() < agent.lastEnergyLevel {
		agent.lastEnergyCost = agent.lastEnergyLevel - agent.GetEnergyLevel()
		agent.energySpent += agent.lastEnergyCost
	} else {
		agent.recalculateSatisfaction()
		agent.badTeam = false
		if agent.satisfactionOfRecentAllocation < 0.5 && agent.energySpent/agent.lootBoxCnt > agent.GetEnergyLevel() {
			// if not enough energy received and not fair allocation happened
			agent.badTeam = true
		}
	}
	agent.lastEnergyLevel = agent.GetEnergyLevel()
	agent.updateRepMap()
	if agent.badTeam {
		return objects.ChangeBike
	}
	return objects.Pedal
}

// DecideForce considering Hegselmann-Krause model, Ramirez-Cano-Pitt model and Satisfaction
func (agent *SmartAgent) DecideForce(direction uuid.UUID) {
	pedalForce := 0.0
	if agent.lastPedal == 0 {
		pedalForce = utils.BikerMaxForce
	} else {
		agentsOnBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetAgents()
		scores := make(map[uuid.UUID]float64)
		totalScore := 0.0
		for _, others := range agentsOnBike {
			id := others.GetID()
			rep := agent.reputationMap[id]
			// Cognitive dimension: is same belief?
			// Pareto principle: give more energy to those with more outcome
			// Forgiveness: forgive agents pedal harder recently
			score := rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution
			scores[others.GetID()] = score
			totalScore += score
		}

		for id, score := range scores {
			scores[id] = score / totalScore
		}

		energyCost := 0.0
		for id, weight := range scores {
			energyCost += weight * agent.reputationMap[id]._lastPedal
		}
		pedalForce = agent.lastPedal * (energyCost / agent.lastEnergyCost)
		pedalForce *= agent.satisfactionOfRecentAllocation
	}
	if pedalForce > utils.BikerMaxForce {
		pedalForce = utils.BikerMaxForce
	}
	steeringForce := 0.0
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	if direction != uuid.Nil {
		steeringForce = physics.ComputeOrientation(currentBike.GetPosition(), agent.GetGameState().GetLootBoxes()[direction].GetPosition()) - currentBike.GetOrientation()
	}
	forces := utils.Forces{
		Pedal: pedalForce,
		Brake: 0.0,
		Turning: utils.TurningDecision{
			SteerBike:     true,
			SteeringForce: steeringForce,
		},
	}
	agent.lastPedal = pedalForce
	agent.SetForces(forces)
}

// DecideJoining accept higher reputation score, max the number of agents on bike
func (agent *SmartAgent) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	decision := make(map[uuid.UUID]bool)
	scores := make([]float64, len(pendingAgents))
	for idx, applicant := range pendingAgents {
		rep := agent.reputationMap[applicant]
		// Cognitive dimension: is same belief?
		// Contribution and Achievement
		// Forgiveness: forgive agents pedal harder recently
		// Potential
		scores[idx] = rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution + rep.energyRemain
	}
	sort.Slice(pendingAgents, func(i, j int) bool {
		return scores[i] > scores[j]
	})
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	for idx, agentId := range pendingAgents {
		decision[agentId] = idx+len(currentBike.GetAgents()) < utils.BikersOnBike
	}
	return decision
}

// ChangeBike rank the average reputation score of agents on bike with empty place, go for highest rank one
func (agent *SmartAgent) ChangeBike() (targetId uuid.UUID) {
	highestAvgScore := 0.0
	for id, bike := range agent.GetGameState().GetMegaBikes() {
		if targetId == uuid.Nil {
			targetId = id
		}
		if len(bike.GetAgents()) == utils.BikersOnBike {
			// ignore the bike with no empty space
			continue
		}
		score := 0.0
		for _, biker := range bike.GetAgents() {
			rep := agent.reputationMap[biker.GetID()]
			// Cognitive dimension: is same belief?
			// Contribution and Achievement
			// Forgiveness: forgive agents pedal harder recently
			// Potential
			score += rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution + rep.energyRemain
		}
		score /= float64(len(bike.GetAgents()))
		if score > highestAvgScore {
			highestAvgScore = score
			targetId = id
		}
	}
	return targetId
}

func (agent *SmartAgent) ProposeDirection() uuid.UUID {
	// direction is targetLootBox
	e := agent.decideTargetLootBox(agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetAgents(), agent.GetGameState().GetLootBoxes())
	// An agent has already proposed its proposal (BordaCount)
	if e != nil {
		panic("unexpected error!")
	}
	return agent.targetLootBox.GetID()
}

func (agent *SmartAgent) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	boxesInMap := agent.GetGameState().GetLootBoxes()

	rank := agent.rankTargetProposals(boxesInMap)
	// need to be map[uuid.UUID]voting.LootboxVoteMap
	return rank
}

func (agent *SmartAgent) DecideAllocation() voting.IdVoteMap {
	agent.lootBoxCnt += 1
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	vote, _ := agent.scoreAgentsForAllocation(currentBike.GetAgents())
	return vote
}

// VoteForKickout try to kick out the reputation score below half of the average on bike
func (agent *SmartAgent) VoteForKickout() map[uuid.UUID]int {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	scores := make([]float64, len(currentBike.GetAgents()))
	threshold := 0.0
	kickOutVote := make(map[uuid.UUID]int)
	for idx, onBikeAgent := range currentBike.GetAgents() {
		rep := agent.reputationMap[onBikeAgent.GetID()]
		// Cognitive dimension: is same belief?
		// Contribution and Achievement
		// Forgiveness: forgive agents pedal harder recently
		// Potential
		scores[idx] = rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution + rep.energyRemain
		threshold += scores[idx]
	}
	threshold = threshold / float64(len(currentBike.GetAgents())) / 2.0
	for idx, onBikeAgent := range currentBike.GetAgents() {
		if scores[idx] < threshold {
			kickOutVote[onBikeAgent.GetID()] = 1
		} else {
			kickOutVote[onBikeAgent.GetID()] = 0
		}
	}
	return kickOutVote
}

// VoteDictator not prefer a dictatorship, if have to then follow the same logic with choosing leader
func (agent *SmartAgent) VoteDictator() voting.IdVoteMap {
	return agent.VoteLeader()
}

func (agent *SmartAgent) VoteLeader() voting.IdVoteMap {
	// defaults to voting for first agent in the list
	agentsOnBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()].GetAgents()
	lootboxes := agent.GetGameState().GetLootBoxes()
	votes := agent.vote_leader(agentsOnBike, lootboxes)
	return votes
}

// DictateDirection assume that if this agent is chosen to be dictator,
// then it is believed by others that the choice of this agent is correct,
// then no extra modification but following the same target choice as normal is acceptable
func (agent *SmartAgent) DictateDirection() uuid.UUID {
	return agent.ProposeDirection()
}

func (agent *SmartAgent) DecideKickOut() []uuid.UUID {
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	scores := make([]float64, len(currentBike.GetAgents()))
	threshold := 0.0
	kickOutVote := make(map[uuid.UUID]int)
	for idx, onBikeAgent := range currentBike.GetAgents() {
		rep := agent.reputationMap[onBikeAgent.GetID()]
		// Cognitive dimension: is same belief?
		// Contribution and Achievement
		// Forgiveness: forgive agents pedal harder recently
		// Potential
		scores[idx] = rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution + rep.energyRemain
		threshold += scores[idx]
	}
	threshold = threshold / float64(len(currentBike.GetAgents())) / 2.0
	count := 0
	for idx, onBikeAgent := range currentBike.GetAgents() {
		if scores[idx] < threshold {
			kickOutVote[onBikeAgent.GetID()] = 1
			count += 1
		} else {
			kickOutVote[onBikeAgent.GetID()] = 0
		}
	}
	decideKickOut := make([]uuid.UUID, count)
	for idx, decision := range kickOutVote {
		if decision == 1 {
			count -= 1
			decideKickOut[count] = idx
		}
	}
	return decideKickOut
}

func (agent *SmartAgent) DecideDictatorAllocation() voting.IdVoteMap {
	return agent.DecideAllocation()
}

func (agent *SmartAgent) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	weights := make(map[uuid.UUID]float64)
	totalW := 0.0
	currentBike := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	for _, onBikeAgent := range currentBike.GetAgents() {
		rep := agent.reputationMap[onBikeAgent.GetID()]
		// Cognitive dimension: is same belief?
		// Contribution and Achievement
		// Forgiveness: forgive agents pedal harder recently
		// Potential
		weights[onBikeAgent.GetID()] = rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution + rep.energyRemain
		totalW += weights[onBikeAgent.GetID()]
	}
	for id, w := range weights {
		weights[id] = w / totalW
	}
	return weights
}

func (agent *SmartAgent) vote_off_leader() bool {
	decision_to_vote_off := false
	vote_off := 0.0
	// vote_off: 1.0 means to vote_off the leader
	id := agent.GetID()
	rep := agent.reputationMap[id]
	if (rep.recentGetEnergy == false) && (rep.isSameColor == 0.0) {
		vote_off = 1.0
	}
	if vote_off == 1.0 {
		decision_to_vote_off = true
	}
	return decision_to_vote_off
}

func (agent *SmartAgent) which_governance_method(agentsOnBike []objects.IBaseBiker) utils.Governance {
	//assume agent only accepts democracy or leadership
	// By default, it accpets leadership
	need_deomocracy := 0.0
	need_leadership := 1.0
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
		need_deomocracy = 1.0
		need_leadership = 0.0
		// selfish personality
	}
	if agent_rep.energyRemain > 2*average_energyRemain {
		need_deomocracy = 1.0
		need_leadership = 0.0
		// fear of being taken advantage of
	}

	if need_deomocracy > need_leadership {
		return utils.Democracy
	} else {
		return utils.Leadership
	}
}

func (agent *SmartAgent) vote_leader(agentsOnBike []objects.IBaseBiker, proposedLootBox map[uuid.UUID]objects.ILootBox) voting.IdVoteMap {
	// two-round run-off

	// the first round
	scores1 := make(map[uuid.UUID]float64)
	total_score_1 := 0.0

	for _, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		// Pareto principle: give more energy to those with more outcome
		// Cognitive dimension: is same belief?
		// necessity: must stay alive
		score_1 := rep.historyContribution + rep.lootBoxGet + rep.isSameColor + rep.energyRemain

		scores1[id] = score_1
		total_score_1 += score_1
	}

	for _, others := range agentsOnBike {
		id := others.GetID()
		scores1[id] = scores1[id] / total_score_1 //normalize
	}

	//the second round
	scores2 := make(map[uuid.UUID]float64)
	total_score_2 := 0.0
	for _, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		score_2 := rep.recentContribution // recent progress, Forgiveness if performed bad before
		scores2[id] = score_2
		total_score_2 += score_2
	}

	for _, others := range agentsOnBike {
		id := others.GetID()
		scores2[id] = scores2[id] / total_score_2 //normalize
	}

	// total
	scores := make(map[uuid.UUID]float64)
	for _, others := range agentsOnBike {
		id := others.GetID()
		scores[id] = 0.7*scores1[id] + 0.3*scores2[id]
	}

	var votes voting.IdVoteMap = scores

	return votes
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
			//opinion, trustness from direct experience
			//forgiveness
			//improve trustness by decreasing the risk of no efforts
			score += 0.5*0.4*rep.historyContribution + 0.5*0.2*rep.recentContribution + 0.5*0.4*rep.energyRemain
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

	sum_score := 0.0
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
		sum_score += score
	}
	// We choose to use the Borda count method to pick a proposal because it can mitigate the Condorcet paradox.
	// Borda count needs to get the rank of all candidates to score Borda points.
	// In this case, according to the Gibbard-Satterthwaite Theorem, Borda count is susceptible to tactical voting.
	// The following steps tend to achieve the rank of lootbox proposals according to their scores calculated. We will return the highest rank to pick the agent with it. (Another Borda score would consider reputation function)这个后面如果可以再考虑如果能得到的话

	// normalize
	for _, lootbox := range proposedLootBox {
		scores[lootbox.GetID()] = scores[lootbox.GetID()] / sum_score
	}

	var lootboxVotes voting.LootboxVoteMap = scores

	return lootboxVotes
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
			// Cognitive dimension: is same belief?
			// Pareto principle: give more energy to those with more outcome
			// Forgiveness: forgive agents pedal harder recently
			// Equality: Agents received more energy before should get less this time
			// Need: Agents with lower energyLevel require more, try to meet their need
			score := rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution - rep.energyGain + rep.energyRemain
			scores[id] = score
			totalScore += score
		}
	}

	// normalize scores
	for id, score := range scores {
		scores[id] = score / totalScore
	}

	return scores, nil
}

func (agent *SmartAgent) UpdateGameState(gameState objects.IGameState) {
	agent.BaseBiker.UpdateGameState(gameState)
}

func (agent *SmartAgent) updateRepMap() {
	if agent.reputationMap == nil {
		agent.reputationMap = make(map[uuid.UUID]reputation)
	}
	for _, otherAgent := range agent.GetGameState().GetAgents() {
		rep, exist := agent.reputationMap[otherAgent.GetID()]
		if !exist {
			rep = reputation{}
		}
		rep.updateScore(otherAgent, agent.GetColour())
		agent.reputationMap[otherAgent.GetID()] = rep
	}
}

func (agent *SmartAgent) recalculateSatisfaction() {
	agent.satisfactionOfRecentAllocation = 1.0
	currentBike, e := agent.GetGameState().GetMegaBikes()[agent.GetBike()]
	if !e {
		return
	}
	agentsOnBike := currentBike.GetAgents()
	scores := make([]float64, len(agentsOnBike))
	gains := make([]float64, len(agentsOnBike))
	for idx, others := range agentsOnBike {
		id := others.GetID()
		rep := agent.reputationMap[id]
		// Cognitive dimension: is same belief?
		// Pareto principle: give more energy to those with more outcome
		// Forgiveness: forgive agents pedal harder recently
		// Equality: Agents received more energy before should get less this time
		// Need: Agents with lower energyLevel require more, try to meet their need
		score := rep.isSameColor + rep.historyContribution + rep.lootBoxGet + rep.recentContribution - rep.energyGain + rep.energyRemain
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
			if n > input[j] {
				inversionCnt += 1
			}
			j += 1
		}
	}
	return 1.0 - 2.0*inversionCnt/float64(size*(size-1))
}

// Creates an instance of Team 3 Biker
func NewTeam3Agent(totColours utils.Colour, bikeId uuid.UUID) *SmartAgent {
	baseBiker := objects.GetBaseBiker(totColours, bikeId) // Use the constructor function
	baseBiker.GroupID = 3
	// print
	// fmt.Println("team5Agent: newTeam5Agent: baseBiker: ", baseBiker)
	return &SmartAgent{
		BaseBiker: *baseBiker,
	}
}
