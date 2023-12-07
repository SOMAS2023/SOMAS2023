package agent

import (
	// "SOMAS2023/internal/clients/team2/agent"
	"SOMAS2023/internal/clients/team2/modules"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"fmt"
	"maps"
	"math/rand"

	"github.com/google/uuid"
)

// We vote for ourselves and the agent with the highest social capital.
func (a *AgentTwo) VoteDictator() voting.IdVoteMap {
	votes := make(voting.IdVoteMap)
	agentId, _ := a.Modules.Environment.GetBikerWithMaxSocialCapital(a.Modules.SocialCapital)
	if len(a.GetFellowBikers()) > 1 && agentId != a.GetID() {
		fellowBikers := a.GetFellowBikers()
		for _, fellowBiker := range fellowBikers {
			if fellowBiker.GetID() == agentId || fellowBiker.GetID() == a.GetID() {
				votes[fellowBiker.GetID()] = 0.5
			} else {
				votes[fellowBiker.GetID()] = 0.0
			}
		}
	} else {
		fellowBikers := a.GetFellowBikers()
		for _, fellowBiker := range fellowBikers {
			if fellowBiker.GetID() == a.GetID() {
				votes[fellowBiker.GetID()] = 1.0
			} else {
				votes[fellowBiker.GetID()] = 0.0
			}
		}
	}
	return votes
}

func (a *AgentTwo) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	// All actions have equal weights. Weighting by AgentId based on social capital.
	// We set the weight for an Agent to be equal to its Social Capital.
	weights := make(map[uuid.UUID]float64)
	agents := a.GetFellowBikers()
	for _, agent := range agents {
		// if agent Id is not in the a.Modules.SocialCapital.SocialCapital map, set the weight to 0.5 (neither trust or distrust)
		if _, ok := a.Modules.SocialCapital.SocialCapital[agent.GetID()]; !ok {
			// add agent to the map
			a.Modules.SocialCapital.SocialCapital[agent.GetID()] = 0.5
		}
		weights[agent.GetID()] = a.Modules.SocialCapital.SocialCapital[agent.GetID()]
	}
	return weights
}

func (a *AgentTwo) DecideKickOut() []uuid.UUID {
	// Only called when the agent is the dictator.
	// We kick out the agent with the lowest social capital on the bike.
	// GetBikerWithMinSocialCapital returns only one agent, if more agents with min SC, it randomly chooses one.
	kickOut_agents := make([]uuid.UUID, 0)
	agentId, _ := a.Modules.Environment.GetBikerWithMinSocialCapital(a.Modules.SocialCapital)
	kickOut_agents = append(kickOut_agents, agentId)
	return kickOut_agents
}

func (a *AgentTwo) VoteLeader() voting.IdVoteMap {
	// We vote 0.5 for ourselves if the agent with the highest SC Agent(that we've met so far) on our bike. If we're alone on a bike, we vote 1 for ourselves.
	votes := make(voting.IdVoteMap)
	fellowBikers := a.GetFellowBikers()
	if len(a.GetFellowBikers()) > 0 {
		agentId, _ := a.Modules.Environment.GetBikerWithMaxSocialCapital(a.Modules.SocialCapital)
		for _, fellowBiker := range fellowBikers {
			if fellowBiker.GetID() == agentId {
				votes[fellowBiker.GetID()] = 0.5
			} else if fellowBiker.GetID() == a.GetID() {
				votes[a.GetID()] = 0.5
			} else {
				votes[fellowBiker.GetID()] = 0.0
			}
		}
	} else {
		votes[a.GetID()] = 1.0
	}

	return votes
}

func (a *AgentTwo) DecideGovernance() utils.Governance {
	fmt.Printf("[DecideGovernance] Agent %s has Social Capitals %v\n", a.GetID(), a.Modules.SocialCapital.SocialCapital)
	a.Modules.SocialCapital.UpdateSocialCapital()
	// All possibilities except dictatorship.
	// Need to decide weights for each type of Governance
	// Can add an invalid weighting so that it is not 50/50

	// randomNumber := rand.Float64()
	// if randomNumber < democracyWeight {
	// 	return utils.Democracy
	// } else {
	// 	return utils.Leadership
	// }
	return utils.Democracy
}

func (a *AgentTwo) DecideAllocation() voting.IdVoteMap {
	socialCapital := maps.Clone(a.Modules.SocialCapital.SocialCapital)
	// Iterate through agents in social capital
	for id := range socialCapital {
		// Iterate through fellow bikers
		for _, biker := range a.GetFellowBikers() {
			// If this agent is a fellow biker, move on
			if biker.GetID() == id {
				continue
			}
		}
		// This agent is not a fellow biker - remove it from SC
		delete(socialCapital, id)
	}
	// We give ourselves 1.0
	socialCapital[a.GetID()] = 1.0
	return socialCapital
}

func (a *AgentTwo) DecideDictatorAllocation() voting.IdVoteMap {
	socialCapital := a.DecideAllocation()

	// Calculate the total social capital
	totalSocialCapital := 0.0
	for _, sc := range socialCapital {
		totalSocialCapital += sc
	}

	// Distribute the allocation based on each agent's share of the total social capital
	result := make(voting.IdVoteMap)
	for agentID, sc := range socialCapital {
		result[agentID] = sc / totalSocialCapital
	}
	return result
}

func (a *AgentTwo) VoteForKickout() map[uuid.UUID]int {
	VoteMap := make(map[uuid.UUID]int)
	kickoutThreshold := ChangeBikeSocialCapitalThreshold
	agentTwoID := a.GetID()

	// check all bikers on the bike but ignore ourselves
	for _, agent := range a.GetFellowBikers() {
		if agent.GetID() != agentTwoID {
			_, exists := a.Modules.SocialCapital.SocialCapital[agent.GetID()]

			if a.Modules.SocialCapital.SocialCapital[agent.GetID()] < kickoutThreshold && exists {
				VoteMap[agent.GetID()] = 1
			} else {
				VoteMap[agent.GetID()] = 0
			}

		}
	}

	return VoteMap
}

func (a *AgentTwo) DecideJoining(pendingAgents []uuid.UUID) map[uuid.UUID]bool {
	// Accept all agents we don't know about or are higher in social capital.
	// If we know about them and they have a lower social capital, reject them.

	decision := make(map[uuid.UUID]bool)
	for _, agent := range pendingAgents {
		// If we know about them and they have a higher social capital than threshold, accept them.
		if _, ok := a.Modules.SocialCapital.SocialCapital[agent]; ok {
			if a.Modules.SocialCapital.SocialCapital[agent] > modules.AcceptThreshold {
				decision[agent] = true
			} else {
				decision[agent] = false
			}
		} else {
			decision[agent] = true
		}
	}
	return decision
}

func (a *AgentTwo) ProposeDirection() uuid.UUID {
	agentID, agentColour, agentEnergy := a.GetID(), a.GetColour(), a.GetEnergyLevel()
	optimalLootbox := a.Modules.Environment.GetNearestLootboxByColor(agentID, agentColour)
	nearestLootbox := a.Modules.Environment.GetNearestLootbox(agentID)
	if agentEnergy < modules.EnergyToOptimalLootboxThreshold || optimalLootbox == uuid.Nil {
		fmt.Printf("[PProposeDirection] Agent %s proposed nearest lootbox %s\n", agentID, nearestLootbox)
		return nearestLootbox
	}
	fmt.Printf("[PProposeDirection] Agent %s proposed optimal lootbox %s\n", agentID, optimalLootbox)
	return optimalLootbox

}

func (a *AgentTwo) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	fmt.Printf("[FFinalDirectionVote] Agent %s got proposals %v\n", a.GetID(), proposals)
	fmt.Printf("[FFinalDirectionVote] Agent %s has Social Capitals %v\n", a.GetID(), a.Modules.SocialCapital.SocialCapital)

	votes := make(voting.LootboxVoteMap)

	// Assume we set our own social capital to 1.0, thus need to account for it
	weight := 1.0 / (a.Modules.SocialCapital.GetSum(a.Modules.SocialCapital.SocialCapital) + 1)

	for proposerID, proposal := range proposals {
		scWeight := 0.0
		if proposerID == a.GetID() {
			// If the proposal is our own, we vote for it with full weight
			scWeight = weight
		} else {
			scWeight = weight * a.Modules.SocialCapital.SocialCapital[proposerID]
		}

		// Check if the proposal already exists in votes, if not add it with the calculated weight
		if _, ok := votes[proposal]; !ok {
			votes[proposal] = scWeight
		} else {
			// If the proposal is already there, update it
			votes[proposal] += scWeight
		}
	}
	fmt.Printf("[FFinalDirectionVote] Agent %s voted %v\n", a.GetID(), votes)
	return votes
}

func (a *AgentTwo) ChangeBike() uuid.UUID {
	decisionInputs := modules.DecisionInputs{SocialCapital: a.Modules.SocialCapital, Enviornment: a.Modules.Environment, AgentID: a.GetID()}
	isChangeBike, bikeId := a.Modules.Decision.MakeBikeChangeDecision(decisionInputs)
	if isChangeBike {
		return bikeId
	} else {
		return a.Modules.Environment.BikeId
	}
}

func (bb *AgentTwo) GetGroupID() int {
	return 5
}

func (a *AgentTwo) DecideAction() objects.BikerAction {
	avgSocialCapital := a.Modules.SocialCapital.GetAverage(a.Modules.SocialCapital.SocialCapital)

	if avgSocialCapital >= ChangeBikeSocialCapitalThreshold {
		// Pedal if members of the bike have high social capital.
		return objects.Pedal
	} else {
		// Otherwise, change bikes.
		return objects.ChangeBike
	}
}

func (a *AgentTwo) DecideForce(direction uuid.UUID) {
	if direction == uuid.Nil {
		lootboxId := a.Modules.Environment.GetHighestGainLootbox()
		lootboxPos := a.Modules.Environment.GetLootboxPos(lootboxId)
		a.SetForces(a.Modules.Utils.GetForcesToTarget(a.GetLocation(), lootboxPos))
		return
	}

	a.Modules.VotedDirection = direction

	if a.Modules.Environment.IsAudiNear() {
		fmt.Printf("[DecideForce] Agent %s is near Audi\n", a.GetID())
		// Move in opposite direction to Audi in full force
		bikePos, audiPos := a.Modules.Environment.GetBike().GetPosition(), a.Modules.Environment.GetAudi().GetPosition()
		force := a.Modules.Utils.GetForcesToTargetWithDirectionOffset(utils.BikerMaxForce, 1.0-a.Modules.Environment.GetBikeOrientation(), bikePos, audiPos)
		a.SetForces(force)
		return
	}
	// Use the average social capital to decide whether to pedal in the voted direciton or not
	probabilityOfConformity := a.Modules.SocialCapital.GetAverage(a.Modules.SocialCapital.SocialCapital)
	randomNumber := rand.Float64()
	agentPosition := a.GetLocation()
	lootboxID := direction
	if randomNumber > probabilityOfConformity {
		lootboxID = a.Modules.Environment.GetHighestGainLootbox()
	}
	lootboxPosition := a.Modules.Environment.GetLootboxPos(lootboxID)
	force := a.Modules.Utils.GetForcesToTargetWithDirectionOffset(utils.BikerMaxForce, -a.Modules.Environment.GetBikeOrientation(), agentPosition, lootboxPosition)
	a.SetForces(force)
}

func (a *AgentTwo) DictateDirection() uuid.UUID {
	// Move in opposite direction to Audi in full force
	if a.Modules.Environment.IsAudiNear() {
		// fmt.Printf("[DictateDirection] Agent %s is near Audi\n", a.GetID())
		return a.Modules.Environment.GetNearestLootboxAwayFromAudi()
	}
	// Otherwise, move towards the lootbox with the highest gain
	return a.Modules.Environment.GetHighestGainLootbox()
}

func (a *AgentTwo) UpdateGameState(gameState objects.IGameState) {
	a.BaseBiker.UpdateGameState(gameState)
	a.Modules.Environment.SetGameState(gameState)
}

func (a *AgentTwo) SetBike(bikeId uuid.UUID) {
	a.Modules.Environment.BikeId = bikeId
	a.BaseBiker.SetBike(bikeId)
}
