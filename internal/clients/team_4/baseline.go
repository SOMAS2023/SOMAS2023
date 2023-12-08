package team_4

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"

	"github.com/google/uuid"
)

type IBaselineAgent interface {
	objects.IBaseBiker

	DecideAction() objects.BikerAction
	ChangeBike() uuid.UUID

	////////////////// opinion.go ///////////////////////
	IncreaseHonesty(agentID uuid.UUID, increaseAmount float64)
	DecreaseHonesty(agentID uuid.UUID, decreaseAmount float64)
	CalculateReputation() map[uuid.UUID]float64
	CalculateHonestyMatrix() map[uuid.UUID]float64
	GetReputation() map[uuid.UUID]float64
	QueryReputation(uuid.UUID) float64

	////////////////// goverance.go ///////////////////////
	DecideGovernance() utils.Governance
	DecideJoining(pendinAgents []uuid.UUID) map[uuid.UUID]bool
	VoteForKickout() map[uuid.UUID]int
	VoteLeader() voting.IdVoteMap
	DecideWeights(action utils.Action) map[uuid.UUID]float64
	VoteDictator() voting.IdVoteMap
	DecideKickOut() []uuid.UUID

	////////////////// allocation.go ///////////////////////
	DecideDictatorAllocation() voting.IdVoteMap
	DecideAllocation() voting.IdVoteMap

	////////////////// direction.go ///////////////////////
	nearestLoot() uuid.UUID
	rankTargetProposals(proposedLootBox []objects.ILootBox) (map[uuid.UUID]float64, error)
	ProposeDirection() uuid.UUID
	FinalDirectionVote(proposals map[uuid.UUID]uuid.UUID) voting.LootboxVoteMap
	DecideForce(direction uuid.UUID)
	DictateDirection() uuid.UUID

	////////////////// data.go ///////////////////////
	UpdateDecisionData()
	getHonestyAverage() float64
	getReputationAverage() float64
	rankFellowsReputation(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error) //returns normal rank of fellow bikers reputation
	rankFellowsHonesty(agentsOnBike []objects.IBaseBiker) (map[uuid.UUID]float64, error)    //returns normal rank of fellow bikers honesty
	DisplayFellowsEnergyHistory()
	DisplayFellowsHonesty()
	DisplayFellowsReputation()
}

// general weights
const audiDistanceThreshold = 75
const minEnergyThreshold = 0.2

const audiDistanceWeight = 10.0
const distanceWeight = 5.0
const reputationWeight = 1.0
const honestyWeight = 1.0
const energySpentWeight = 0.5
const energyLevelWeight = 1.0

const minFellowBikers = 6         //if we have less than this number of fellows, we will not kick anyone out
const dictatorMinFellowBikers = 6 //if we have less than this number of fellows, we will not kick anyone out

// for voting for leader and dictator
const leaderRepWeight = 2.0
const leaderHonestWeight = 1.0
const dictatorRepWeight = 2.0
const dictatorHonestWeight = 1.0

type BaselineAgent struct {
	objects.BaseBiker
	lootBoxColour     utils.Colour
	mylocationHistory []utils.Coordinates     //log location history for this agent
	energyHistory     map[uuid.UUID][]float64 //log energy level for all agents
	reputation        map[uuid.UUID]float64   //record reputation for other agents, 0-1
	honestyMatrix     map[uuid.UUID]float64   //record honesty for other agents, 0-1
}

type agentScore struct {
	ID    uuid.UUID
	Score float64
}

// DecideAction only pedal
func (agent *BaselineAgent) DecideAction() objects.BikerAction {
	// fmt.Println("Team 4")
	return objects.Pedal
}
