package team4

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

	////////////////// messaging.go ///////////////////////
	// HandleKickoutMessage(msg KickoutAgentMessage) (uuid.UUID, uuid.UUID)
	// HandleReputationMessage(msg ReputationOfAgentMessage) (uuid.UUID, float64)
	// HandleJoiningMessage(msg JoiningAgentMessage) (uuid.UUID, uuid.UUID)
	// HandleLootboxMessage(msg LootboxMessage) (uuid.UUID, uuid.UUID)
	// HandleGovernanceMessage(msg GovernanceMessage) (uuid.UUID, int)
	// HandleForcesMessage(msg ForcesMessage) (uuid.UUID, utils.Forces)
	// HandleVoteGovernanceMessage(msg VoteGoveranceMessage) (uuid.UUID, voting.IdVoteMap)
	// HandleVoteLootboxDirectionMessage(msg VoteLootboxDirectionMessage) (uuid.UUID, voting.IdVoteMap)
	// HandleVoteRulerMessage(msg VoteRulerMessage) (uuid.UUID, voting.IdVoteMap)
	// HandleVoteKickoutMessage(msg VoteKickoutMessage) (uuid.UUID, voting.IdVoteMap)

	// GetAllMessages([]IBaselineAgent) []messaging.IMessage[IBaselineAgent]
}

// general weights
const audiDistanceThreshold = 75
const minEnergyThreshold = 0.4

const audiDistanceWeight = 8.0
const distanceWeight = 7.0
const reputationWeight = 2.0
const honestyWeight = 1.0
const energySpentWeight = 1.0
const energyLevelWeight = 1.4
const resourceWeight = 1.0

const minFellowBikers = 6         //if we have less than this number of fellows, we will not vote to kick anyone out
const dictatorMinFellowBikers = 6 //if we have less than this number of fellows, we will not kick anyone out

// for voting for leader and dictator
const leaderRepWeight = 2.0
const leaderHonestWeight = 1.0
const dictatorRepWeight = 2.0
const dictatorHonestWeight = 1.0

type BaselineAgent struct {
	*objects.BaseBiker
	currentBike       uuid.UUID
	capacity          int       //number of agents on my bike
	audiTarget        uuid.UUID //current bike audi is targeting
	currentGovernance utils.Governance
	currentRuler      uuid.UUID //ruler of the current bike = uuid.Nil if no ruler
	targetLoot        uuid.UUID
	lootBoxColour     utils.Colour
	lootBoxLocation   utils.Coordinates
	timeInLimbo       int
	onBike            bool
	optimalBike       uuid.UUID               //best bike on the map
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

	if agent.evaluateBike(agent.currentBike) {
		return objects.Pedal
	} else if agent.GetEnergyLevel() <= 0.65 || agent.ChangeBike() == agent.currentBike {
		return objects.Pedal
	} else {
		return objects.ChangeBike
	}

}

// called by the spawner/server to instantiate bikers in the MVP
func GetBiker4(baseBiker *objects.BaseBiker) objects.IBaseBiker {
	team4Agent := &BaselineAgent{
		BaseBiker: baseBiker,
	}
	team4Agent.BaseBiker.GroupID = 4
	//fmt.Println("Team 4 Agent Created", team4Agent.GetID())
	return team4Agent
}
