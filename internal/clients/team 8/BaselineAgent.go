package team_8

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/voting"
	"math"

	"sort"

	"github.com/google/uuid"
)

type Agent8 struct {
	*objects.BaseBiker
	energyLevel float64
	color       string
}

// in the MVP the biker's action defaults to pedaling (as it won't be able to change bikes)
// in future implementations this function will be overridden by the agent's specific strategy
// which will be used to determine whether to pedalor try to change bike
func (bb *Agent8) DecideAction() objects.BikerAction {
	//TODO： need to be implemented
	return 0
}

// determine the forces (pedalling, breaking and turning)
// in the MVP the pedalling force will be 1, the breaking 0 and the tunring is determined by the
// location of the nearest lootbox

// the function is passed in the id of the voted lootbox, for now ignored
func (bb *Agent8) DecideForce(direction uuid.UUID) {
	//TODO： need to be implemented
}

// AgentRanking holds an agent's UUID and their energy level
type AgentRanking struct {
	ID          uuid.UUID
	energyLevel float64
}

func GetAgentByUUID(id uuid.UUID) *Agent8 {
	// Implement the logic to retrieve Agent8 instance
	// This is just a placeholder
	return &Agent8{}
}

// GetEnergyLevel returns the energy level of the agent.
func (bb *Agent8) GetEnergyLevel() float64 {
	return bb.energyLevel
}

func (bb *Agent8) GetColor() string {
	return bb.color
}

type AgentBordaRank struct {
	ID         uuid.UUID
	BordaPoint float64
}

func RankAgentsWithEnergy(pendingAgents []uuid.UUID, weight float64) []AgentBordaRank {
	type agentWithEnergy struct {
		ID          uuid.UUID
		EnergyLevel float64
	}

	var agentsWithEnergy []agentWithEnergy
	for _, agentID := range pendingAgents {
		agent := GetAgentByUUID(agentID)
		agentsWithEnergy = append(agentsWithEnergy, agentWithEnergy{ID: agentID, EnergyLevel: agent.energyLevel})
	}

	// Sorting agents by energy level in descending order
	sort.Slice(agentsWithEnergy, func(i, j int) bool {
		return agentsWithEnergy[i].EnergyLevel > agentsWithEnergy[j].EnergyLevel
	})

	var AgentsScore []AgentBordaRank
	totalAgents := len(agentsWithEnergy)
	for i, agent := range agentsWithEnergy {
		bordaPoint := float64(totalAgents-i) * weight // Calculate Borda point
		AgentsScore = append(AgentsScore, AgentBordaRank{ID: agent.ID, BordaPoint: bordaPoint})
	}

	return AgentsScore
}

func SumBordaScores(canons ...[]AgentBordaRank) map[uuid.UUID]float64 {
	summedScores := make(map[uuid.UUID]float64)

	// Iterate over each canon
	for _, canon := range canons {
		// Iterate over each agent's Borda score in the canon
		for _, agentRank := range canon {
			// Sum the scores for each agent
			summedScores[agentRank.ID] += agentRank.BordaPoint
		}
	}

	return summedScores
}

func AssignPointsBasedOnColor(agentsMap map[uuid.UUID]*Agent8, ourAgentColor string, weighting float64) []AgentBordaRank {
	var ranks []AgentBordaRank
	totalAgents := float64(len(agentsMap))
	for id, agent := range agentsMap {
		points := 0.0
		if agent.color == ourAgentColor {
			points = totalAgents * weighting / 2 // Assign points if the color matches
		}
		ranks = append(ranks, AgentBordaRank{ID: id, BordaPoint: points})
	}
	return ranks
}

func (bb *Agent8) CalculateSingleAgentScore(agentID uuid.UUID, threshold float64) map[uuid.UUID]bool {
	agent := GetAgentByUUID(agentID)

	colorValue := 0.0
	energyScore := 0.0
	if agent.color == bb.color {
		colorValue = 1.0
	}
	if agent.energyLevel > bb.energyLevel {
		energyScore = 1.0
	}

	totalScore := energyScore + colorValue

	decisions := make(map[uuid.UUID]bool)
	decisions[agentID] = totalScore >= threshold

	return decisions
}

// an agent will have to rank the agents that are trying to join and that they will try to
func (bb *Agent8) DecideJoining(pendingAgents []uuid.UUID,

// n int, // number of agents we want to allow to join the bike
) map[uuid.UUID]bool {

	decisions := make(map[uuid.UUID]bool)

	if len(pendingAgents) == 0 {
		agentID := uuid.UUID{}
		decisions := bb.CalculateSingleAgentScore(agentID, 1)

		return decisions

	} else {
		agentsMap := make(map[uuid.UUID]*Agent8)

		// Functions for various canons
		// Weighting for each canon should be confirmed afterwards
		AgentsScore_energy := RankAgentsWithEnergy(pendingAgents, 0.3)
		ownAgentColor := bb.color
		AgentScore_color := AssignPointsBasedOnColor(agentsMap, ownAgentColor, 0.7)
		// AgentsScore_reputation := ()

		// sum weighted scores from different canons
		summedScores := SumBordaScores(AgentsScore_energy, AgentScore_color)

		// rank agents and select the best agent
		var summedRankings []AgentBordaRank
		for id, score := range summedScores {
			summedRankings = append(summedRankings, AgentBordaRank{ID: id, BordaPoint: score})
		}

		sort.Slice(summedRankings, func(i, j int) bool {
			return summedRankings[i].BordaPoint > summedRankings[j].BordaPoint
		})

		for i, agentRank := range summedRankings {
			decisions[agentRank.ID] = i < 1 // this should be the n from input afterwards
		}
	}
	return decisions
}

// decide which bike to go to
// for now it just returns a random uuid
func (bb *Agent8) ChangeBike() uuid.UUID {
	//TODO： need to be implemented
	return uuid.New()
}

// default implementation returns the id of the nearest lootbox
func (bb *Agent8) ProposeDirection() uuid.UUID {
	//TODO： need to be implemented
	return uuid.New()
}

// //  Nemo started

// this function will contain the agent's strategy on deciding which direction to go to
// the default implementation returns an equal distribution over all options
// this will also be tried as returning a rank of options
func (bb *Agent8) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {
	votes := make(voting.LootboxVoteMap)

	// Calculate the distance to each proposed loot box
	distances := make(map[uuid.UUID]float64)
	for _, proposal := range proposals {
		for _, lootBox := range bb.GetGameState().GetLootBoxes() {
			if lootBox.GetID() == proposal {
				currLocation := bb.GetLocation()
				x, y := lootBox.GetPosition().X, lootBox.GetPosition().Y
				distance := math.Sqrt(math.Pow(currLocation.X-x, 2) + math.Pow(currLocation.Y-y, 2))
				distances[proposal] = distance
				break
			}
		}
	}

	// Find the loot box with the minimum distance (or lowest rank)
	var minDistance float64 = math.MaxFloat64
	var chosenBox uuid.UUID
	for id, dist := range distances {
		if dist < minDistance {
			minDistance = dist
			chosenBox = id
		}
	}

	// Vote for the chosen box
	for _, proposal := range proposals {
		if proposal == chosenBox {
			votes[proposal] = 1.0 // Full vote for the chosen box
		} else {
			votes[proposal] = 0.0 // No vote for the other boxes
		}
	}

	return votes
}

// through this function the agent submits their desired allocation of resources
// in the MVP each agent returns 1 whcih will cause the distribution to be equal across all of them
func (bb *Agent8) DecideAllocation() voting.IdVoteMap {
	//TODO： need to be implemented
	distribution := make(map[uuid.UUID]float64)
	return distribution
}
