package agent

import (
	"SOMAS2023/internal/clients/team2/modules"
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"math/rand"

	"github.com/google/uuid"
)

func (a *AgentTwo) DecideWeights(action utils.Action) map[uuid.UUID]float64 {
	// TODO: All actions have equal weights. Weighting by AgentId based on social capital.
	return a.BaseBiker.DecideWeights(action)
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
	// TODO: All possibilities except dictatorship.
	return a.BaseBiker.DecideGovernance()
}

func (a *AgentTwo) DecideAllocation() voting.IdVoteMap {
	// TODO: We simply pass in Social Capital values in the map.
	// If a value does not exist in the map, we set it as the average social capital.
	// We give ourselves the highest social capital which is 1.
	return a.BaseBiker.DecideAllocation()
}

func (a *AgentTwo) VoteForKickout() map[uuid.UUID]int {
	// TODO: Vote for the agents with a Social Capital lower than a threshold.
	return a.BaseBiker.VoteForKickout()
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
	if agentEnergy < modules.EnergyToOptimalLootboxThreshold {
		return nearestLootbox
	}
	return optimalLootbox
}

func (a *AgentTwo) FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap {
	// TODO: If Social Capital of agent who proposed a lootbox is higher than a threshold, vote for it. Weight based on SC.
	// Otherwise, set a weight of 0.
	return a.BaseBiker.FinalDirectionVote(proposals)
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

func (a *AgentTwo) DecideAction() objects.BikerAction {
	avgSocialCapital := a.Modules.SocialCapital.GetAverage(a.Modules.SocialCapital.SocialCapital)

	if avgSocialCapital > ChangeBikeSocialCapitalThreshold {
		// Pedal if members of the bike have high social capital.
		return objects.Pedal
	} else {
		// Otherwise, change bikes.
		return objects.ChangeBike
	}
}

func (a *AgentTwo) DecideForce(direction uuid.UUID) {

	a.Modules.VotedDirection = direction

	if a.Modules.Environment.IsAudiNear() {
		// Move in opposite direction to Audi in full force
		bikePos, audiPos := a.Modules.Environment.GetBike().GetPosition(), a.Modules.Environment.GetAudi().GetPosition()
		force := a.Modules.Utils.GetForcesToTargetWithDirectionOffset(utils.BikerMaxForce, -180.0, bikePos, audiPos)
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
	force := a.Modules.Utils.GetForcesToTarget(agentPosition, lootboxPosition)
	a.SetForces(force)
}

func (a *AgentTwo) DictateDirection() uuid.UUID {
	// Move in opposite direction to Audi in full force
	if a.Modules.Environment.IsAudiNear() {
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
